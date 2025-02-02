package database

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var Client *mongo.Client = ConnectToDb()

func ConnectToDb() *mongo.Client {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("Error in Loading the Env file")
	}
	mongoUri := os.Getenv("MONGO_URI")
	if mongoUri == "" {
		log.Fatal("Mongo url missing from env")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	fmt.Println("Context", ctx)
	defer cancel()
	mongoClient, err := mongo.Connect(ctx, options.Client().ApplyURI(mongoUri))
	if err != nil {
		log.Fatal("Error in mongoDb Connection")
	}
	err = mongoClient.Ping(ctx, nil)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Connected to DB")
	fmt.Println(mongoClient)
	return mongoClient
}

func GetCollection(collectionName string) *mongo.Collection {
	dbName := os.Getenv("DB_NAME") // Add DB_NAME in your .env file
	if dbName == "" {
		log.Fatal("Database name is missing from .env")
	}

	return Client.Database(dbName).Collection(collectionName)
}
