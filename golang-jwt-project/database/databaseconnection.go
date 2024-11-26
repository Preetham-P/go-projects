package database

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func DBInstance() (*mongo.Client, context.Context, context.CancelFunc) {
	if err := godotenv.Load(".env"); err != nil {
		log.Fatal("There was an error loading env file", err)
	}

	mongodburl := os.Getenv("MONGODBURL")

	clientOptions := options.Client().ApplyURI(mongodburl)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)

	client, err := mongo.Connect(ctx, clientOptions)

	if err != nil {
		log.Fatal("There was an error connecting to the mongodb ", err)
	}

	if err = client.Ping(ctx, nil); err != nil {
		log.Fatalf("Could not ping MongoDB: %v", err)
	}

	fmt.Println("Connected to MongoDB!")

	return client, ctx, cancel
}
