package service

import (
	"context"
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"
)

type DoneParam struct {
	JWT    string `json:"jwt"`
	TaskID string `json:"taskID"`
}

func (s *Service) DoneHandle(c *gin.Context) {
	p := DoneParam{}
	if err := c.BindJSON(&p); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})
	}

	err := s.TaskDone(context.TODO(), p.JWT, p.TaskID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "ok",
	})
}

type NewParam struct {
	JWT         string     `json:"jwt"`
	Description string     `json:"description"`
	Status      TaskStatus `json:"status"`
}

func (s *Service) NewHandle(c *gin.Context) {
	p := NewParam{}
	if err := c.BindJSON(&p); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})
	}

	taskID, err := s.NewTask(context.TODO(), p.Description, p.Status, s.popugs.GetRandom().ID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})
	}

	if err := s.taskEvents.Created(context.TODO(), taskID); err != nil {
		slog.Error("error create task", err)

		c.JSON(http.StatusOK, gin.H{
			"message": taskID,
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"message": taskID,
	})
}

type ShuffleParam struct {
	JWT string `json:"jwt"`
}

func (s *Service) ShuffleHandle(c *gin.Context) {
	jwt := ShuffleParam{}
	if err := c.BindJSON(&jwt); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})
	}

	if err := s.ShuffleTasks(context.TODO(), jwt.JWT); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "ok",
	})
}
