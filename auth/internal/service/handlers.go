package service

import (
	"context"
	"fmt"
	"net/http"

	"github.com/Nerzal/gocloak/v13"
	"github.com/gin-gonic/gin"
)

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
