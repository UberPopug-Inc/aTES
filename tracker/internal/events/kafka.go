package events

import (
	"context"
	"log"
	"os"

	"github.com/UberPopug-Inc/aTES/auth/internal/service"
	"github.com/segmentio/kafka-go"
)

const (
	topic         = "tracker"
	brokerAddress = "localhost:29092"
)

type Kafka struct {
	writer *kafka.Writer
}

func NewKafka() *Kafka {
	l := log.New(os.Stdout, "kafka writer: ", 0)

	k := &Kafka{writer: kafka.NewWriter(kafka.WriterConfig{
		Brokers: []string{brokerAddress},
		Topic:   topic,
		Logger:  l,
	})}
	k.writer.AllowAutoTopicCreation = true

	return k
}

func (k *Kafka) Done(ctx context.Context, taskID string) error {
	return k.writer.WriteMessages(ctx, kafka.Message{
		Topic: topic,
		Key:   []byte("task_done"),
		Value: []byte(taskID),
	})
}

func (k *Kafka) Created(ctx context.Context, taskID string) error {
	return k.writer.WriteMessages(ctx, kafka.Message{
		Topic: topic,
		Key:   []byte("task_created"),
		Value: []byte(taskID),
	})
}

func (k *Kafka) Assigned(ctx context.Context, task service.Task) error {
	return k.writer.WriteMessages(ctx, kafka.Message{
		Topic: topic,
		Key:   []byte("task_assigned"),
		Value: []byte(task.String()),
	})
}
