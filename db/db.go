package db

import (
	"context"
	"fmt"
	"log"
	"os"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	client       *mongo.Client
	collUsers    *mongo.Collection
	collTheaters *mongo.Collection
	collSales    *mongo.Collection
	collSessions *mongo.Collection
)

// ConnectMongoDb connects to MongoDB
func ConnectMongoDb() error {
	mongoURI := os.Getenv("MONGODB_URI")
	if mongoURI == "" {
		log.Println("MONGODB_URI environment variable is not set")
		return fmt.Errorf("MONGODB_URI environment variable is not set")
	}

	clientOptions := options.Client().ApplyURI(mongoURI)
	var err error
	client, err = mongo.Connect(context.Background(), clientOptions)
	if err != nil {
		log.Println("Error connecting to MongoDB:", err)
		return fmt.Errorf("could not connect to MongoDB: %v", err)
	}

	// if err = client.Ping(context.Background(), nil); err != nil {
	// 	log.Println("Error pinging MongoDB:", err)
	// 	return fmt.Errorf("could not ping MongoDB: %v", err)
	// }

	collUsers = client.Database(os.Getenv("MONGO_DB_KGO")).Collection(os.Getenv("MONGO_KGO_USERS"))
	collSessions = client.Database(os.Getenv("MONGO_DB_KGO")).Collection(os.Getenv("MONGO_KGO_SESSIONS"))
	collTheaters = client.Database(os.Getenv("MONGO_DB_MFX")).Collection(os.Getenv("MONGO_MFX_THEATERS"))
	collSales = client.Database(os.Getenv("MONGO_DB_SUP")).Collection(os.Getenv("MONGO_KGO_SALSES"))

	log.Println("Connected to MongoDB!")
	return nil
}

func CollUsers() *mongo.Collection {
	if collUsers == nil {
		log.Panic("CollUsers is not initialized.")
	}
	return collUsers
}

func CollSessions() *mongo.Collection {
	if collSessions == nil {
		log.Panic("collSessions is not initialized.")
	}
	return collSessions
}

func CollTheaters() *mongo.Collection {
	if collTheaters == nil {
		log.Panic("collTheaters is not initialized.")
	}
	return collTheaters
}

func CollSales() *mongo.Collection {
	if collSales == nil {
		log.Panic("collSales is not initialized.")
	}
	return collSales
}
