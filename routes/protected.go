package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/juniorAkp/delivery-go/controller"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

func ProtectedRoute(router *gin.RouterGroup, client *mongo.Client) {

	router.POST("/orders", controller.CreateOrder(client))

}
