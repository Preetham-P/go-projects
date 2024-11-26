package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
	ID           primitive.ObjectID `bson:"_id"`
	FirstName    *string            `json:"first_name" validate:"required,min=2,max=100"`
	LastName     *string            `json:"last_name" validate:"required,min=2,max=100"`
	Email        *string            `json:"email" validate:"email,required,min=10,max=100"`
	Phone        *string            `json:"phone" validate:"required"`
	UserType     *string            `json:"user_type" validate:"required,eq=ADMIN|eq=USER"`
	Password     *string            `json:"password" validate:"required,min=2,max=100"`
	Token        *string            `json:"token"`
	RefreshToken *string            `json:"refreshtoken"`
	CreatedAt    time.Time          `json:"CreatedAt"`
	UpdateAt     time.Time          `json:"UpdateAt" `
	CreatedBy    *string            `json:"CreatedBy"`
	ModifiedBy   *string            `json:"ModifiedBy"`
	User_Id      string             `json:"user_id"`
}
