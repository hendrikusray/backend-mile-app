package config

import (
	"context"
	"fmt"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var MongoClient *mongo.Client

func Connect() *mongo.Client {
	clientOptions := options.Client().
		ApplyURI(Cold.DBMongoURI).
		SetMaxPoolSize(Cold.DBMongoMaxPoolSize).
		SetMinPoolSize(Cold.DBMongoMinPoolSize)

	client, err := mongo.Connect(context.Background(), clientOptions)
	if err != nil {
		log.Fatalf("❌ Failed to connect MongoDB: %v", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := client.Ping(ctx, nil); err != nil {
		log.Fatalf("❌ Ping MongoDB failed: %v", err)
	}

	MongoClient = client
	fmt.Println("✅ MongoDB connected successfully to:", Cold.DBMongoURI)
	return client
}
