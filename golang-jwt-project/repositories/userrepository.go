package repositories

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/Preetham-P/go-projects/golang-jwt-project/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func GetUsers(client *mongo.Client, ctx context.Context, cancel context.CancelFunc, database string, collection string, filter interface{}) (foundUsers []models.User, returnError error) {
	coll := client.Database(database).Collection(collection)
	ctx, cancel = context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var users []models.User

	cursor, err := coll.Find(ctx, filter)
	if err != nil {
		return nil, fmt.Errorf("Find error: %v", err)
	}
	defer cursor.Close(ctx)

	if err := cursor.All(ctx, &users); err != nil {
		return nil, errors.New("There was an error fetching users, ")
	}

	return users, nil
}

func GetUser(client *mongo.Client, ctx context.Context, cancel context.CancelFunc, database string, collection string, filter interface{}) (returnuser models.User, err error) {
	coll := client.Database(database).Collection(collection)
	defer cancel()
	var user models.User
	if err := coll.FindOne(ctx, filter).Decode(&user); err != nil {
		return user, errors.New("Cannot fetch user with id ")
	}
	return user, nil
}

func CountUsers(client *mongo.Client, ctx context.Context, cancel context.CancelFunc, database string, collection string, filter interface{}) (count int64, err error) {
	coll := client.Database(database).Collection(collection)
	//defer cancel()

	count, err = coll.CountDocuments(ctx, filter)

	if err != nil {
		fmt.Printf(err.Error())
		return 0, err
	}
	return count, nil
}

func CreateUser(client *mongo.Client, ctx context.Context, cancel context.CancelFunc, database string, collection string, user models.User) (mongoresult *mongo.InsertOneResult, err error) {
	coll := client.Database(database).Collection(collection)
	defer cancel()

	insertResult, err := coll.InsertOne(ctx, user)
	if err != nil {
		return nil, fmt.Errorf("Find error: %v", err)
	}
	return insertResult, nil

}

func UpdateUser(client *mongo.Client, ctx context.Context, cancel context.CancelFunc, database string, collection string, filter interface{}, updateObj primitive.D, options options.UpdateOptions) (updateresult *mongo.UpdateResult, err error) {
	coll := client.Database(database).Collection(collection)
	defer cancel()

	updateresult, _ = coll.UpdateOne(ctx,
		filter,
		bson.D{
			{Key: "$set", Value: updateObj},
		}, &options)

	return updateresult, nil
}
