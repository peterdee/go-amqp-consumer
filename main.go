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

			insert := func(body []byte) {
				ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
				defer cancel()

				var quote interface{}
				parseError := bson.UnmarshalExtJSON(body, true, &quote)
				if parseError != nil {
					log.Fatal(parseError)
				}
				_, insertionError := mongodb.Quotes.InsertOne(ctx, quote)
				if insertionError != nil {
					log.Fatal(insertionError)
				}
			}

			insert(message.Body)
		}
	}()

	fmt.Println("Waiting for incoming messages...")
	<-forever
}
