package events

import (
	"context"
	"log"
	"os"

	"github.com/Nerzal/gocloak/v13"
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

func (k *Kafka) Logged(ctx context.Context, user gocloak.UserInfo) error {
	return k.writer.WriteMessages(ctx, kafka.Message{
		Topic: "user_logged",
		Key:   []byte("event"),
		Value: []byte(user.String()),
	})
}

func (k *Kafka) Created(ctx context.Context, userID string) error {
	return k.writer.WriteMessages(ctx, kafka.Message{
		Topic: "user_created",
		Key:   []byte("event"),
		Value: []byte(userID),
	})
}

func (k *Kafka) Updated(ctx context.Context, user gocloak.User) error {
	return k.writer.WriteMessages(ctx, kafka.Message{
		Topic: "user_updated",
		Key:   []byte("event"),
		Value: []byte(user.String()),
	})
}

func (k *Kafka) Deleted(ctx context.Context, userID string) error {
	return k.writer.WriteMessages(ctx, kafka.Message{
		Topic: "user_deleted",
		Key:   []byte("event"),
		Value: []byte(userID),
	})
}
