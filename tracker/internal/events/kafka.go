package events

import (
	"context"
	"log"
	"os"
	"time"

	"github.com/UberPopug-Inc/aTES/tracker/internal/service"
	"github.com/google/uuid"
	"github.com/segmentio/kafka-go"
)

const (
	brokerAddress = "localhost:29092"
)

type Kafka struct {
	writer *kafka.Writer
}

func NewKafka() *Kafka {
	l := log.New(os.Stdout, "kafka writer: ", 0)

	k := &Kafka{writer: &kafka.Writer{
		Addr:                   kafka.TCP(brokerAddress),
		RequiredAcks:           kafka.RequireAll,
		Logger:                 l,
		AllowAutoTopicCreation: true,
	}}

	return k
}

func (k *Kafka) Done(ctx context.Context, task service.Task) error {
	event := TaskV1{
		Title:        "Task v1",
		Description:  "Descr v1",
		EventID:      uuid.New().String(),
		EventVersion: 1,
		EventName:    "created",
		EventTime:    time.Now().String(),
		Data: TaskDataV1{
			TaskID:     task.ID,
			TaskTitle:  task.Description,
			AssignUUID: task.WorkerID,
		},
	}

	if err := event.validate(); err != nil {
		return err
	}

	return k.writer.WriteMessages(ctx, kafka.Message{
		Topic: "task_done",
		Key:   []byte("event"),
		Value: []byte(event.string()),
	})
}

func (k *Kafka) Created(ctx context.Context, task service.Task) error {
	event := TaskV1{
		Title:        "Task v1",
		Description:  "Descr v1",
		EventID:      uuid.New().String(),
		EventVersion: 1,
		EventName:    "created",
		EventTime:    time.Now().String(),
		Data: TaskDataV1{
			TaskID:     task.ID,
			TaskTitle:  task.Description,
			AssignUUID: task.WorkerID,
		},
	}

	if err := event.validate(); err != nil {
		return err
	}

	return k.writer.WriteMessages(ctx, kafka.Message{
		Topic: "task_done",
		Key:   []byte("event"),
		Value: []byte(event.string()),
	})
}

func (k *Kafka) Assigned(ctx context.Context, task service.Task) error {
	event := TaskV1{
		Title:        "Task v1",
		Description:  "Descr v1",
		EventID:      uuid.New().String(),
		EventVersion: 1,
		EventName:    "created",
		EventTime:    time.Now().String(),
		Data: TaskDataV1{
			TaskID:     task.ID,
			TaskTitle:  task.Description,
			AssignUUID: task.WorkerID,
		},
	}

	if err := event.validate(); err != nil {
		return err
	}

	return k.writer.WriteMessages(ctx, kafka.Message{
		Topic: "task_assigned",
		Key:   []byte("event"),
		Value: []byte(event.string()),
	})
}
