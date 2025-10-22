package database

import (
	"log"
	"os"

	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

func Connect() *mongo.Client {
	err := godotenv.Load(".env")
	if err != nil {
		panic(err)
	}

	mongoUri := os.Getenv("MONGODB_URI")
	if mongoUri == "" {
		mongoUri = "mongodb://localhost:27017"
	}

	clientOptions := options.Client().ApplyURI(mongoUri)

	client, err := mongo.Connect(clientOptions)
	if err != nil {
		log.Fatal("failed to connect to mongo db",err)
	}
	return  client
}


func OpenCollection(collectionName string, client *mongo.Client) *mongo.Collection {
	err := godotenv.Load(".env")
	if err != nil {
		panic(err)
	}

	dbName := os.Getenv("DBNAME")

	collection := client.Database(dbName).Collection(collectionName)
	if collection == nil {
		return  nil
	}
	return  collection
}