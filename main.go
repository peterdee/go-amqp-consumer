package main

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	amqp "github.com/rabbitmq/amqp091-go"
)

func main() {
	envError := godotenv.Load()
	if envError != nil {
		log.Fatal("Could not load .env file!")
	}

	// Connect to RabbitMQ
	rabbitMQHost := os.Getenv("RABBITMQ_HOST")
	rabbitMQPassword := os.Getenv("RABBITMQ_PASSWORD")
	rabbitMQPort := os.Getenv("RABBITMQ_PORT")
	rabbitMQUser := os.Getenv("RABBITMQ_USER")
	if rabbitMQHost == "" || rabbitMQPassword == "" ||
		rabbitMQPort == "" || rabbitMQUser == "" {
		log.Fatal("Could not load RabbitMQ configuration")
	}
	rabbitMQConnection, connectionError := amqp.Dial(
		fmt.Sprintf(
			"amqp://%s:%s@%s:%s/",
			rabbitMQUser,
			rabbitMQPassword,
			rabbitMQHost,
			rabbitMQPort,
		),
	)
	if connectionError != nil {
		log.Fatal("Could not connect to RabbitMQ:", connectionError)
	}

	channel, channelError := rabbitMQConnection.Channel()
	if channelError != nil {
		log.Fatal(channelError)
	}
	rabbitMQQueue, queueError := channel.QueueDeclare(
		"quotes",
		false,
		false,
		false,
		false,
		nil,
	)
	if queueError != nil {
		log.Fatal(queueError)
	}

	messages, consumeError := channel.Consume(
		rabbitMQQueue.Name,
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
			fmt.Printf("Received a message: %s", message.Body)
		}
	}()

	fmt.Println("Waiting for incoming messages...")
	<-forever
}
