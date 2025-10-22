package main

import (
	"os"

	"github.com/gin-gonic/gin"
)

func main() {
	port := os.Getenv("PORT")

	if port == "" {
		port = ":8080"
	}

	r := gin.Default()

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