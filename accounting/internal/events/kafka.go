package events

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/segmentio/kafka-go"
)

const (
	brokerAddress = "localhost:29092"
)

type Kafka struct {
	readerAssign *kafka.Reader
	readerDone   *kafka.Reader
}

func (k *Kafka) AssignedPull(ctx context.Context) (*TaskV1, error) {
	msg, err := k.readerAssign.ReadMessage(ctx)
	if err != nil {
		panic("could not read message " + err.Error())
	}
	// after receiving the message, log its value
	fmt.Println("received: ", msg)

	task := &TaskV1{}

	if err := json.Unmarshal(msg.Value, task); err != nil {
		return nil, err
	}

	return task, nil
}

func (k *Kafka) DonePull(ctx context.Context) (*TaskV1, error) {
	msg, err := k.readerDone.ReadMessage(ctx)
	if err != nil {
		panic("could not read message " + err.Error())
	}
	// after receiving the message, log its value
	fmt.Println("received: ", msg)

	task := &TaskV1{}

	if err := json.Unmarshal(msg.Value, task); err != nil {
		return nil, err
	}

	return task, nil
}

func NewKafka() *Kafka {
	l := log.New(os.Stdout, "kafka writer: ", 0)

	k := &Kafka{
		readerAssign: kafka.NewReader(kafka.ReaderConfig{
			Brokers: []string{brokerAddress},
			GroupID: "my-group",
			Logger:  l,
			Topic:   "task_assign",
		}),
		readerDone: kafka.NewReader(kafka.ReaderConfig{
			Brokers: []string{brokerAddress},
			GroupID: "my-group",
			Logger:  l,
			Topic:   "task_done",
		}),
	}

	return k
}
