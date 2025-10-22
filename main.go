package main

import (
	"context"
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/juniorAkp/delivery-go/database"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

func main() {
	port := os.Getenv("PORT")

	if port == "" {
		port = ":8080"
	}

	r := gin.Default()

	var client *mongo.Client = database.Connect()

	if err := client.Ping(context.Background(),nil); err != nil {
		log.Fatal("Failed to connect to mongo db")
	}

	defer func() {
		if err := client.Disconnect(context.Background()); err != nil {
			panic(err)
		}
	}()

	api := r.Group("/api/v1") 
	{
		api.GET("/ping",func(c *gin.Context) {
		c.JSON(200,gin.H{
			"message": "pong",
		})
	})
	}

	r.Run(port)
}