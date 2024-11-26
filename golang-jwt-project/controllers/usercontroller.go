package controllers

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/Preetham-P/go-projects/golang-jwt-project/database"
	"github.com/Preetham-P/go-projects/golang-jwt-project/helpers"
	"github.com/Preetham-P/go-projects/golang-jwt-project/models"
	"github.com/Preetham-P/go-projects/golang-jwt-project/repositories"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
	"gopkg.in/mgo.v2/bson"
)

var ctx context.Context
var client *mongo.Client
var cancel context.CancelFunc
var validate *validator.Validate
var dbName string

func init() {
	client, ctx, cancel = database.DBInstance()
	validate = validator.New()

	if err := godotenv.Load(".env"); err != nil {
		log.Fatal("There was an error loading env file", err)
	}

	dbName = os.Getenv("DATABASE")
}

func HashPassword(password string) (hp string, err error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), 14)

	if err != nil {
		log.Panic(err)
		return "", err
	}
	return string(hashedPassword), nil

}

func VerifyPassword(userPassword string, foundUserPassword string) (valid bool, msg string) {
	err := bcrypt.CompareHashAndPassword([]byte(userPassword), []byte(foundUserPassword))
	check := true
	message := ""

	if err != nil {
		check = false
		message = fmt.Sprintf("Password doesnt match")
	}

	message = "Passwords match"
	return check, message
}

func SignUp() gin.HandlerFunc {
	return func(c *gin.Context) {
		var user models.User

		if err := c.BindJSON(&user); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
			return
		}

		if validationerr := validate.Struct(user); validationerr != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": validationerr.Error(),
			})
			return
		}

		emailfilter := bson.M{"email": user.Email}
		phoneFilter := bson.M{"phone": user.Phone}
		count, err := repositories.CountUsers(client, ctx, cancel, dbName, "users", emailfilter)
		if err != nil {
			log.Panic(err)
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
			return
		}
		count, err = repositories.CountUsers(client, ctx, cancel, dbName, "users", phoneFilter)
		if err != nil {
			log.Panic(err)
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
			return
		}

		if count > 0 {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "The user with this email or phone number already exist",
			})
		}

		hp, err := HashPassword(*user.Password)
		user.Password = &hp
		user.CreatedAt, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
		user.UpdateAt, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
		user.ID = primitive.NewObjectID()
		user.User_Id = user.ID.Hex()
		token, refreshToken := helpers.GenerateAlltokens(*user.Email, *user.FirstName, *user.LastName, *user.UserType, *&user.User_Id)
		user.Token = &token
		user.RefreshToken = &refreshToken
		result, err := repositories.CreateUser(client, ctx, cancel, dbName, "users", user)

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
		}
		c.JSON(http.StatusOK, result)
	}

}

func Login() gin.HandlerFunc {
	return func(c *gin.Context) {
		var user models.User
		var foundUser models.User

		if err := c.BindJSON(&user); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
		}

		userFilter := bson.M{"email": user.Email}

		foundUser, err := repositories.GetUser(client, ctx, cancel, dbName, "users", userFilter)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{
				"error": err.Error(),
			})
			return
		}

		passwordValid, msg := VerifyPassword(*user.Password, *foundUser.Password)

		if !passwordValid {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": msg,
			})
			return
		}

		token, refreshToken := helpers.GenerateAlltokens(*foundUser.Email, *foundUser.FirstName, *foundUser.LastName, *foundUser.UserType, foundUser.User_Id)
		updateObj, options := helpers.UpdateAllTokens(token, refreshToken, foundUser.User_Id)

		filter := bson.M{"user_id": foundUser.User_Id}

		updateResult, err := repositories.UpdateUser(client, ctx, cancel, dbName, "users", filter, *updateObj, options)

		if err != nil {
			fmt.Sprintf(err.Error())
		}
		fmt.Sprintf("Updated tokens successfuly", updateResult)

		foundUser, err = repositories.GetUser(client, ctx, cancel, dbName, "users", userFilter)

		c.JSON(http.StatusOK, foundUser)

	}

}

func GetUsers() gin.HandlerFunc {
	return func(c *gin.Context) {

		filter := bson.M{}

		users, err := repositories.GetUsers(client, ctx, cancel, dbName, "users", filter)

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
			return
		}
		c.JSON(http.StatusOK, users)
	}
}

func GetUserById() gin.HandlerFunc {
	return func(c *gin.Context) {
		user_id := c.Param("user_id")

		if err := helpers.MatchUserTypeToUid(c, user_id); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"errpr": err.Error()})
			return
		}

		if err := godotenv.Load(".env"); err != nil {
			log.Fatal("There was an error loading env file", err)
		}

		dbName := os.Getenv("database")

		userFilter := bson.M{"user_id": user_id}

		user, err := repositories.GetUser(client, ctx, cancel, dbName, "users", userFilter)

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, user)
	}

}
