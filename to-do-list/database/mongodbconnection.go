package database

import (
	"context"
	"fmt"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func Connect() (*mongo.Client, context.Context, context.CancelFunc) {
	clientOptions := options.Client().ApplyURI("mongodb://localhost:27017")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)

	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		log.Fatalf("MongoDB connection error: %v", err)
	}

	// Test the connection
	if err = client.Ping(ctx, nil); err != nil {
		log.Fatalf("Could not ping MongoDB: %v", err)
	}

	fmt.Println("Connected to MongoDB!")
	return client, ctx, cancel
}
