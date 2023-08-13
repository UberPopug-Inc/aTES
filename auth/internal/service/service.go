package service

import (
	"context"
	"crypto/tls"
	"log"

	"github.com/Nerzal/gocloak/v13"
	"github.com/gin-gonic/gin"
	"github.com/go-resty/resty/v2"
)

const realm = "master"

type Service struct {
	rest        *gin.Engine
	client      *gocloak.GoCloak
	restyClient *resty.Client
	userEvents  UserEvents
}

func New(userEvents UserEvents) *Service {
	s := &Service{
		userEvents: userEvents,
	}

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
	jwt, err := s.client.Login(ctx, "account", "", username, password, realm)
	if err != nil {
		return nil, err
	}
	info, err := s.client.GetUserInfo(ctx, jwt.AccessToken, realm)
	if err != nil {
		return nil, err
	}

	if err := s.userEvents.Logged(ctx, *info); err != nil {
		return nil, err
	}

	return jwt, nil
}

func (s *Service) NewUser(ctx context.Context, jwt string, user gocloak.User) (string, error) {
	createUser, err := s.client.CreateUser(ctx, jwt, realm, user)
	if err != nil {
		return "", nil
	}

	if err := s.userEvents.Created(ctx, createUser); err != nil {
		return "", err
	}

	return createUser, nil
}

func (s *Service) UpdateUser(ctx context.Context, jwt string, user gocloak.User) error {
	if err := s.client.UpdateUser(ctx, jwt, realm, user); err != nil {
		return err
	}

	if err := s.userEvents.Updated(ctx, user); err != nil {
		return err
	}

	return nil
}

func (s *Service) DeleteUser(ctx context.Context, jwt string, userID string) error {
	if err := s.client.DeleteUser(ctx, jwt, realm, userID); err != nil {
		return err
	}

	if err := s.userEvents.Deleted(ctx, userID); err != nil {
		return err
	}

	return nil
}
