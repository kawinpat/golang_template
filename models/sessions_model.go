package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Session defines the structure of a session to be stored in the `sessions` collection
type Session struct {
	ID        primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	UserID    primitive.ObjectID `bson:"user_id" json:"user_id"`
	Token     string             `bson:"token" json:"token"`
	CreatedAt time.Time          `bson:"created_at" json:"created_at"`
}
