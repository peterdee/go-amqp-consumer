package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/bson"

	"go-amqp-consumer/mongodb"
	"go-amqp-consumer/rabbitmq"
)

func main() {
	envError := godotenv.Load()
	if envError != nil {
		log.Fatal("Could not load .env file!")
	}

	mongodb.CreateConnection()
	rabbitmq.CreateConnection()

	messages, consumeError := rabbitmq.Channel.Consume(
		rabbitmq.Queue.Name,
		"",
		true,
		false,
		false,
		false,
		nil,
	)
	if consumeError != nil {
		log.Fatal(consumeError)
	}
	var forever chan struct{}

	go func() {
		for message := range messages {
			fmt.Printf(" -> Received a message:\n%s\n\n", message.Body)
			ctx, cancel = context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()
			res, err := collection.InsertOne(ctx, bson.D{{"name", "pi"}, {"value", 3.14159}})
			id := res.InsertedID
		}
	}()

	fmt.Println("Waiting for incoming messages...")
	<-forever
}
