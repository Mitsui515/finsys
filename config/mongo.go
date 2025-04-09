package config

import (
	"context"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var MongoDB *mongo.Database

type MongoDBConfig struct {
	URI      string
	Database string
	Timeout  time.Duration
}

func DefaultMongoDBConfig() *MongoDBConfig {
	return &MongoDBConfig{
		URI:      "mongodb://localhost:27017",
		Database: "finsys",
		Timeout:  10 * time.Second,
	}
}

func InitMongoDB() *mongo.Database {
	mongoConfig := DefaultMongoDBConfig()
	ctx, cancel := context.WithTimeout(context.Background(), mongoConfig.Timeout)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(mongoConfig.URI))
	if err != nil {
		log.Fatalf("Failed to connect to MongoDB: %v", err)
	}

	err = client.Ping(ctx, nil)
	if err != nil {
		log.Fatalf("Failed to ping MongoDB: %v", err)
	}

	MongoDB = client.Database(mongoConfig.Database)
	log.Println("Successfully connected to MongoDB")
	return MongoDB
}
