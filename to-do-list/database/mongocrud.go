package database

import (
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func CreateDocument(client *mongo.Client, ctx context.Context, database, collection string, document interface{}) (*mongo.InsertOneResult, error) {
	coll := client.Database(database).Collection(collection)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	result, err := coll.InsertOne(ctx, document)
	if err != nil {
		fmt.Println("error occured", err)
		return nil, fmt.Errorf("InsertOne error: %v", err)
	}
	return result, nil
}

func FindOneDocument(client *mongo.Client, ctx context.Context, database, collection string, filter interface{}) (bson.M, error) {
	coll := client.Database(database).Collection(collection)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	var result bson.M
	err := coll.FindOne(ctx, filter).Decode(&result)
	if err != nil {
		return nil, fmt.Errorf("FindOne error: %v", err)
	}
	return result, nil
}

func FindDocuments(client *mongo.Client, ctx context.Context, database, collection string, filter interface{}) ([]bson.M, error) {
	coll := client.Database(database).Collection(collection)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	cursor, err := coll.Find(ctx, filter)
	if err != nil {
		return nil, fmt.Errorf("Find error: %v", err)
	}
	defer cursor.Close(ctx)

	var results []bson.M
	if err := cursor.All(ctx, &results); err != nil {
		return nil, fmt.Errorf("Cursor error: %v", err)
	}
	return results, nil
}

func UpdateDocument(client *mongo.Client, ctx context.Context, database, collection string, filter, update interface{}) (*mongo.UpdateResult, error) {
	coll := client.Database(database).Collection(collection)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	result, err := coll.UpdateOne(ctx, filter, update)
	if err != nil {
		return nil, fmt.Errorf("UpdateOne error: %v", err)
	}
	return result, nil
}

func DeleteDocument(client *mongo.Client, ctx context.Context, database, collection string, filter interface{}) (*mongo.DeleteResult, error) {
	coll := client.Database(database).Collection(collection)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	result, err := coll.DeleteOne(ctx, filter)
	if err != nil {
		return nil, fmt.Errorf("DeleteOne error: %v", err)
	}
	return result, nil
}
