package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/juniorAkp/delivery-go/controller"
	"github.com/juniorAkp/delivery-go/middleware"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

func ProtectedRoute(router *gin.RouterGroup, client *mongo.Client) {
	router.Use(middleware.AuthRequired())
	router.POST("/orders", controller.CreateOrder(client))

}
