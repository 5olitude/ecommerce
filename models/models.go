package models

import (
	"time"

	"startup/go/pkg/mod/go.mongodb.org/mongo-driver@v1.7.2/bson/primitive"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
	ID            primitive.ObjectID `json:"_id" bson:"_id"`
	First_Name    *string            `json:"first_name" validate:"required,min=2,max=30"`
	Last_Name     *string            `json:"last_name"  validate:"required,min=2,max=30"`
	Password      *string            `json:"password"   validate:"required,min=6"`
	Email         *string            `json:"email"      validate:"email,required"`
	Phone         *string            `json:"phone"      validate:"phone,required"`
	Token         *string            `json:"token"`
	Refresh_Token *string            `josn:"refresh_token"`
	Created_At    time.Time          `json:"created_at"`
	Updated_At    time.Time          `json:"updtaed_at"`
	User_ID       string             `json:"user_id"`
}
