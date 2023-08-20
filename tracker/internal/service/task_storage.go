package service

import (
	"encoding/json"
	"math/rand"

	"github.com/google/uuid"
)

type TaskStatus uint8

const (
	Inprogress = TaskStatus(iota)
	Done
)

type Task struct {
	ID          string     `json:"ID,omitempty"`
	Description string     `json:"description,omitempty"`
	Status      TaskStatus `json:"status,omitempty"`
	WorkerID    string     `json:"workerID,omitempty"`
	Cost        uint32     `json:"cost,omitempty"`
}

type TaskStorage struct {
	tasks map[string]Task
}

func NewTaskStorage() *TaskStorage {
	return &TaskStorage{tasks: make(map[string]Task)}
}

func (t *TaskStorage) Done(ID string) {
	task := t.tasks[ID]
	task.Status = Done
	t.tasks[ID] = task
}

func (t *TaskStorage) Add(description string, status TaskStatus, workerID string) string {
	tt := Task{
		ID:          uuid.NewString(),
		Description: description,
		Status:      status,
		WorkerID:    workerID,
		Cost:        uint32(rand.Intn(1000)),
	}

	t.tasks[tt.ID] = tt

	return tt.ID
}

func (t *Task) String() string {
	bytes, err := json.Marshal(t)
	if err != nil {
		return ""
	}

	return string(bytes)
}
