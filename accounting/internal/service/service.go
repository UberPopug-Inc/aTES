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
	taskEvents  TaskEventer
	popugs      *PopugsStorage
	tasks       *TaskStorage
}

func New(taskEvents TaskEventer) *Service {
	s := &Service{
		taskEvents: taskEvents,
		popugs:     NewPopugsStorage(),
		tasks:      NewTaskStorage(),
	}

	s.rest = gin.Default()
	s.rest.POST("/done", s.DoneHandle)
	s.rest.POST("/new", s.NewHandle)
	s.rest.POST("/shuffle", s.ShuffleHandle)

	s.client = gocloak.NewClient("http://localhost:8080/accounting")
	s.restyClient = s.client.RestyClient()
	s.restyClient.SetDebug(true)
	s.restyClient.SetTLSClientConfig(&tls.Config{InsecureSkipVerify: true})

	return s
}

func (s *Service) Run() error {
	if err := s.rest.Run("localhost:8889"); err != nil {
		log.Panicf("rest.Run(): %v", err)
	}

	return nil
}

func (s *Service) TaskDone(ctx context.Context, jwt string, taskID string) error {
	// TODO: JWT analyze

	s.tasks.Done(taskID)

	if err := s.taskEvents.Done(ctx, taskID); err != nil {
		return err
	}

	return nil
}

func (s *Service) NewTask(ctx context.Context, desc string, status TaskStatus, workerID string) (string, error) {
	taskID := s.tasks.Add(desc, status, workerID)

	if err := s.taskEvents.Created(ctx, taskID); err != nil {
		return "", err
	}

	return taskID, nil
}

func (s *Service) ShuffleTasks(ctx context.Context, jwt string) error {
	// TODO: JWT

	for key, value := range s.tasks.tasks {
		t := value
		t.WorkerID = s.popugs.GetRandom().ID
		s.tasks.tasks[key] = t

		if err := s.taskEvents.Assigned(ctx, t); err != nil {
			return err
		}
	}

	return nil
}
