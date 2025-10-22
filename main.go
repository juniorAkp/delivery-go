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
	"go.mongodb.org/mongo-driver/v2/mongo"
)

func main() {
	port := os.Getenv("PORT")

	if port == "" {
		port = "8080"
	}

	var client *mongo.Client = database.Connect()

	if err := client.Ping(context.Background(), nil); err != nil {
		log.Fatal("Failed to connect to mongo db")
	}
	log.Println("mongoDb connected")

	r := NewRouter()

	srv := &http.Server{
		Addr:    ":" + port,
		Handler: r,
		// set timeout due CWE-400 - Potential Slowloris Attack
		ReadHeaderTimeout: 5 * time.Second,
	}

	//start server in background
	go func() {
		log.Printf("server is listening on %s", port)
		if err := srv.ListenAndServe(); err != nil {
			log.Printf("Failed to start server: %v", err)
		}
	}()

	stopChan := make(chan os.Signal, 1)
	signal.Notify(stopChan, os.Interrupt, syscall.SIGTERM)
	<-stopChan

	log.Println("initializing graceful shutdown")

	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer shutdownCancel()

	if err := srv.Shutdown(shutdownCtx); err != nil {
		log.Printf("Http server shutdown error: %v", err)
	} else {
		log.Println("Http server shutdown")
	}

	if err := client.Disconnect(shutdownCtx); err != nil {
		log.Printf("MongoDB disconnect error: %v", err)
	} else {
		log.Println("MongoDB disconnected")
	}

	log.Println("Shutdown complete")
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
	}

	return r
}
