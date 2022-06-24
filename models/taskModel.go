package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

//User is the model that governs all notes objects retrived or inserted into the DB
type Task struct {
	ID          primitive.ObjectID `bson:"_id"`
	First_name  *string            `json:"first_name" validate:"required,min=1,max=100"`
	Last_name   *string            `json:"last_name" validate:"required,min=1,max=100"`
	Task_name   *string            `json:"taskName"`
	Module_code *string            `json:"moduleCode"`
	Duration    *string            `json:"duration"`
	Hidden      bool               `json:"hidden"`
	Created_at  time.Time          `json:"created_at"`
	Updated_at  time.Time          `json:"updated_at"`
	User_id     string             `json:"user_id"`
}
