package service

import (
	"context"
	"crypto/tls"
	"fmt"
	"log"
	"net/http"

	"github.com/Nerzal/gocloak/v13"
	"github.com/gin-gonic/gin"
	"github.com/go-resty/resty/v2"
)

const realm = "master"

type Service struct {
	rest        *gin.Engine
	client      *gocloak.GoCloak
	restyClient *resty.Client
}

type EventsBroker interface {
}

func New() *Service {
	s := &Service{}

	s.rest = gin.Default()
	s.rest.POST("/login", s.LoginHandle)
	s.rest.POST("/new", s.NewUserHandle)
	s.rest.POST("/update", s.UpdateUserHandle)
	s.rest.GET("/delete", s.DeleteUserHandle)

	s.client = gocloak.NewClient("http://localhost:8080/auth")
	s.restyClient = s.client.RestyClient()
	s.restyClient.SetDebug(true)
	s.restyClient.SetTLSClientConfig(&tls.Config{InsecureSkipVerify: true})

	return s
}

func (s *Service) Run() error {
	if err := s.rest.Run("localhost:8888"); err != nil {
		log.Panicf("rest.Run(): %v", err)
	}

	return nil
}

func (s *Service) Login(ctx context.Context, username, password string) (*gocloak.JWT, error) {
	jwt, err := s.client.LoginClient(ctx, username, password, realm)
	if err != nil {
		return nil, err
	}

	return jwt, nil
}

func (s *Service) NewUser(ctx context.Context, jwt string, user gocloak.User) (string, error) {
	createUser, err := s.client.CreateUser(ctx, jwt, realm, user)
	if err != nil {
		return "", nil
	}
	return createUser, nil
}

func (s *Service) UpdateUser(ctx context.Context, jwt string, user gocloak.User) error {
	if err := s.client.UpdateUser(ctx, jwt, realm, user); err != nil {
		return err
	}

	return nil
}

func (s *Service) DeleteUser(ctx context.Context, jwt string, userID string) error {
	if err := s.client.DeleteUser(ctx, jwt, realm, userID); err != nil {
		return err
	}

	return nil
}

type UserLoginParams struct {
	User     string `json:"user"`
	Password string `json:"password"`
}

func (s *Service) LoginHandle(c *gin.Context) {
	ul := UserLoginParams{}

	if err := c.BindJSON(&ul); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})
	}

	fmt.Println(ul)

	jwt, err := s.Login(context.TODO(), ul.User, ul.Password)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"jwt": "",
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"jwt": jwt,
	})
}

type UserParam struct {
	JWT  string `json:"jwt"`
	User gocloak.User
}

func (s *Service) NewUserHandle(c *gin.Context) {
	jwt := UserParam{}
	if err := c.BindJSON(&jwt); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})
	}

	userID, err := s.NewUser(context.TODO(), jwt.JWT, jwt.User)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"userID": userID,
	})
}

func (s *Service) UpdateUserHandle(c *gin.Context) {
	jwt := UserParam{}
	if err := c.BindJSON(&jwt); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})
	}

	if err := s.UpdateUser(context.TODO(), jwt.JWT, jwt.User); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "ok",
	})
}

func (s *Service) DeleteUserHandle(c *gin.Context) {
	jwt := UserParam{}
	if err := c.BindJSON(&jwt); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})
	}

	if err := s.DeleteUser(context.TODO(), jwt.JWT, *jwt.User.ID); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "ok",
	})
}
