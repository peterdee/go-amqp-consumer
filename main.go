package main

import (
	"fmt"
	"log"

	"github.com/joho/godotenv"

	"go-amqp-consumer/rabbitmq"
)

func main() {
	envError := godotenv.Load()
	if envError != nil {
		log.Fatal("Could not load .env file!")
	}

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
		}
	}()

	fmt.Println("Waiting for incoming messages...")
	<-forever
}
