package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/juniorAkp/delivery-go/database"
	"github.com/juniorAkp/delivery-go/middleware"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

func main() {
	port := os.Getenv("PORT")

	if port == "" {
		port = "8080"
	}

	r := NewRouter()

	srv := &http.Server{
		Addr:    ":" + port,
		Handler: r,
		// set timeout due CWE-400 - Potential Slowloris Attack
		ReadHeaderTimeout: 5 * time.Second,
	}

	if err := srv.ListenAndServe(); err != nil {
		log.Printf("Failed to start server: %v", err)
	}
}

func NewRouter() *gin.Engine {
	r := gin.Default()
	var client *mongo.Client = database.Connect()

	if err := client.Ping(context.Background(), nil); err != nil {
		log.Fatal("Failed to connect to mongo db")
	}

	defer func() {
		if err := client.Disconnect(context.Background()); err != nil {
			panic(err)
		}
	}()

	r.Use(middleware.Cors())

	api := r.Group("/api/v1")
	{
		api.GET("/ping", func(c *gin.Context) {
			c.JSON(200, gin.H{
				"message": "pong",
			})
		})
	}

	return r
}
