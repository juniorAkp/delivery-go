package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/juniorAkp/delivery-go/controller"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

func UnprotectedRoute(router *gin.Engine, client *mongo.Client) {
	router.GET("/products", controller.GetProducts(client))
	router.GET("/products/:productID", controller.GetProduct(client))
	router.POST("products", controller.CreateProduct(client))

	router.POST("/register", controller.RegisterCustomer(client))
	router.POST("/login", controller.LoginCustomer(client))
}
