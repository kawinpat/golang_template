package repository

import (
	"go.mongodb.org/mongo-driver/mongo"
)

// UserRepository holds the client reference for MongoDB operations
type UserRepository struct {
	client     *mongo.Client
	collection *mongo.Collection
}

var (
	databaseName   = "sample_kgolang" // Database name
	collectionName = "users"          // Collection name
)

// NewUserRepository initializes and returns a UserRepository instance
func NewUserRepository(client *mongo.Client) *UserRepository {
	collection := client.Database(databaseName).Collection(collectionName)
	return &UserRepository{
		client:     client,
		collection: collection, // Store the collection in the repository
	}
}

// Collection returns the MongoDB collection for the repository
func (repo *UserRepository) Collection() *mongo.Collection {
	return repo.collection
}
