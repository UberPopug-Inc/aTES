package service

import (
	"context"
	"log"
	"math/rand"

	"github.com/gin-gonic/gin"
)

const realm = "master"

type Service struct {
	rest           *gin.Engine
	taskEvents     TaskEventer
	billingStorage *BillingStorage
}

func New(taskEvents TaskEventer) *Service {
	s := &Service{
		taskEvents:     taskEvents,
		billingStorage: NewPopugsStorage(),
	}

	s.rest = gin.Default()
	s.rest.POST("/stat", s.StatHandle)
	s.rest.POST("/list", s.ListHandle)

	ctx := context.TODO()

	go s.RunAssign(ctx)

	go s.RunDone(ctx)

	return s
}

func (s *Service) Run() error {
	if err := s.rest.Run("localhost:8889"); err != nil {
		log.Panicf("rest.Run(): %v", err)
	}

	return nil
}

func (s *Service) RunAssign(ctx context.Context) error {
	for {
		select {
		case <-ctx.Done():
			return nil
		default:
		}

		task, err := s.taskEvents.AssignedPull(ctx)
		if err == nil {
			s.billingStorage.Credit(task.Data.AssignUUID, uint64(rand.Intn(100)))
		}
	}
}

func (s *Service) RunDone(ctx context.Context) error {
	for {
		select {
		case <-ctx.Done():
			return nil
		default:
		}

		task, err := s.taskEvents.DonePull(ctx)
		if err == nil {
			s.billingStorage.Debit(task.Data.AssignUUID, uint64(rand.Intn(100)))
		}
	}
}
