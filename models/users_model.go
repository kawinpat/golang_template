package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// User defines the structure of the user collection
type Users struct {
	ID        primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Name      string             `bson:"name" json:"name"`
	Email     string             `bson:"email" json:"email"`
	Username  string             `bson:"username" json:"username"`
	Password  string             `bson:"password" json:"password"`
	CreatedAt time.Time          `bson:"created_at" json:"created_at"` // Store creation timestamp
	UpdatedAt time.Time          `bson:"updated_at" json:"updated_at"` // Store last update timestamp
}
