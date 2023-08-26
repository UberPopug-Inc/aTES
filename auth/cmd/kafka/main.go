package main

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/segmentio/kafka-go"
)

const (
	topic         = "message-log"
	brokerAddress = "localhost:9092"
)

func main() {
	// create a new context
	ctx := context.Background()
	// produce messages in a new go routine, since
	// both the produce and consume functions are
	// blocking

	go produce(ctx)

	consume(ctx)
}

func produce(ctx context.Context) {
	// initialize a counter
	i := 0

	// Construct a synchronous writer (the default mode).
	w := &kafka.Writer{
		Addr:                   kafka.TCP(brokerAddress),
		Topic:                  topic,
		RequiredAcks:           kafka.RequireAll,
		AllowAutoTopicCreation: true,
	}

	defer w.Close()

	for {
		err := w.WriteMessages(ctx, kafka.Message{
			Key: []byte(strconv.Itoa(i)),
			// create an arbitrary message payload for the value
			Value: []byte("this is message" + strconv.Itoa(i)),
		})
		if err != nil {
			panic("could not write message " + err.Error())
		}

		// log a confirmation once the message is written
		fmt.Println("writes:", i)
		i++
		// sleep for a second
		time.Sleep(time.Second)
	}
}

func consume(ctx context.Context) {
	r := kafka.NewReader(kafka.ReaderConfig{
		Brokers: []string{brokerAddress},
		GroupID: "my-group",
		Topic:   topic,
	})

	defer r.Close()

	for {
		msg, err := r.ReadMessage(ctx)
		if err != nil {
			panic("could not read message " + err.Error())
		}
		// after receiving the message, log its value
		fmt.Println("received: ", msg)
	}
}
