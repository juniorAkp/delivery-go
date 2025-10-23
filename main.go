package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/juniorAkp/delivery-go/database"
	"github.com/juniorAkp/delivery-go/middleware"
	"github.com/juniorAkp/delivery-go/routes"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

var client *mongo.Client = database.Connect()

func main() {
	port := os.Getenv("PORT")

	if port == "" {
		port = "8080"
	}

	if err := client.Ping(context.Background(), nil); err != nil {
		log.Fatal("Failed to connect to mongo db")
	}

	r := NewRouter()

	srv := &http.Server{
		Addr:              ":" + port,
		Handler:           r,
		ReadHeaderTimeout: 5 * time.Second,
	}

	go func() {
		log.Printf("server is listening on %s", port)
		if err := srv.ListenAndServe(); err != nil {
			log.Printf("Failed to start server: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	<-quit

	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer shutdownCancel()

	if err := srv.Shutdown(shutdownCtx); err != nil {
		log.Printf("Http server shutdown error: %v", err)
	}

	if err := client.Disconnect(shutdownCtx); err != nil {
		log.Printf("MongoDB disconnect error: %v", err)
	}
}

func NewRouter() *gin.Engine {
	r := gin.Default()

	r.Use(middleware.Cors())

	api := r.Group("/api/v1")
	{
		api.GET("/ping", func(c *gin.Context) {
			c.JSON(200, gin.H{
				"message": "pong",
			})
		})

		routes.UnprotectedRoute(api, client)
	}

	return r
}
