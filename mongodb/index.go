package mongodb

import (
	"context"
	"log"
	"os"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

var Client *mongo.Client

var Database *mongo.Database

var Quotes *mongo.Collection

func CreateConnection() {
	connectionString := os.Getenv("MONGODB_CONNECTION_STRING")
	database := os.Getenv("MONGODB_DATABASE")
	if connectionString == "" || database == "" {
		log.Fatal("Could not load MongoDB configuration")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, connectionError := mongo.Connect(ctx, options.Client().ApplyURI(connectionString))
	if connectionError != nil {
		log.Fatal("Could not connect to MongoDB:", connectionError)
	}
	pingError := client.Ping(ctx, readpref.Primary())
	if connectionError != nil {
		log.Fatal("Could not ping MongoDB:", pingError)
	}

	Client = client

	Database = Client.Database(database)

	Quotes = Database.Collection("Quotes")

	log.Println("Connected to MongoDB")
}
