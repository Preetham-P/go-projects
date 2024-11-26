package helpers

import (
	"log"
	"os"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type SignedDetails struct {
	Email      string
	First_Name string
	Last_Name  string
	User_Type  string
	User_Id    string
	jwt.StandardClaims
}

var secret_key string

func init() {
	if err := godotenv.Load(".env"); err != nil {
		log.Fatal("There was an error loading env file", err)
	}

	secret_key = os.Getenv("SECRET_KEY")
}

func UpdateAllTokens(token string, refreshToken string, user_id string) (updateObject *primitive.D, updateoptions options.UpdateOptions) {

	var updateObj primitive.D

	updateObj = append(updateObj, bson.E{"token", token})
	updateObj = append(updateObj, bson.E{"refreshtoken", refreshToken})

	updatedAt, _ := time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
	updateObj = append(updateObj, bson.E{"updated_at", updatedAt})

	upsert := true

	opt := options.UpdateOptions{
		Upsert: &upsert,
	}

	return &updateObj, opt

}

func GenerateAlltokens(email string, first_name string, last_name string, user_type string, user_id string) (token string, refresh_token string) {

	claims := &SignedDetails{
		Email:      email,
		First_Name: first_name,
		Last_Name:  last_name,
		User_Type:  user_type,
		User_Id:    user_id,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Local().Add(time.Hour * time.Duration(24)).Unix(),
		},
	}

	refreshClaims := &SignedDetails{
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Local().Add(time.Hour * time.Duration(24)).Unix(),
		},
	}

	token, err := jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString([]byte(secret_key))
	refresh_token, Newerr := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims).SignedString([]byte(secret_key))

	if err != nil {
		log.Panic(err)
		return
	}
	if Newerr != nil {
		log.Panic(Newerr)
		return
	}
	return token, refresh_token
}

func ValidateToken(signedtoken string) (claims *SignedDetails, err error) {
	token, err := jwt.ParseWithClaims(signedtoken, &SignedDetails{}, func(t *jwt.Token) (interface{}, error) {
		return []byte(secret_key), nil
	})
	if err != nil {
		return
	}

	claims, ok := token.Claims.(*SignedDetails)
	if !ok {
		return
	}
	return claims, nil
}
