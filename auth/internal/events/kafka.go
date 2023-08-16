package events

import (
	"context"
	"log"
	"os"

	"github.com/Nerzal/gocloak/v13"
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

func (k *Kafka) Logged(ctx context.Context, user gocloak.UserInfo) error {
	return k.writer.WriteMessages(ctx, kafka.Message{
		Key:   []byte("logged_user"),
		Value: []byte(user.String()),
	})
}

func (k *Kafka) Created(ctx context.Context, userID string) error {
	return k.writer.WriteMessages(ctx, kafka.Message{
		Key:   []byte("created_user"),
		Value: []byte(userID),
	})
}

func (k *Kafka) Updated(ctx context.Context, user gocloak.User) error {
	return k.writer.WriteMessages(ctx, kafka.Message{
		Key:   []byte("updated_user"),
		Value: []byte(user.String()),
	})
}

func (k *Kafka) Deleted(ctx context.Context, userID string) error {
	return k.writer.WriteMessages(ctx, kafka.Message{
		Key:   []byte("deleted_user"),
		Value: []byte(userID),
	})
}
