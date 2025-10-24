package controller

import (
	"context"
	"errors"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/juniorAkp/delivery-go/database"
	"github.com/juniorAkp/delivery-go/model"
	"github.com/juniorAkp/delivery-go/utils"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type CreateOrderRequest struct {
	ProductId string `json:"productId"`
	Quantity  int64  `json:"quantity" validate:"required gt=0"`
}

func CreateOrder(client *mongo.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(c, 100*time.Second)
		defer cancel()

		userId, err := utils.GetUserId(c)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		var req CreateOrderRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		if err := validate.Struct(req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "validation failed", "details": err.Error()})
			return
		}

		orderCollection := database.OpenCollection("orders", client)
		productCollection := database.OpenCollection("products", client)
		customerCollection := database.OpenCollection("customers", client)

		var customer model.Customer
		err = customerCollection.FindOne(ctx, bson.M{"_id": userId}).Decode(&customer)
		if errors.Is(err, mongo.ErrNoDocuments) {
			c.JSON(http.StatusNotFound, gin.H{"error": "Customer not found"})
			return
		}
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to verify customer"})
			return
		}

		var product model.Product
		err = productCollection.FindOne(ctx, bson.M{"_id": req.ProductId}).Decode(&product)
		if errors.Is(err, mongo.ErrNoDocuments) {
			c.JSON(http.StatusNotFound, gin.H{"error": "Product not found"})
			return
		}
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch product"})
			return
		}

		if product.StockQuantity < req.Quantity {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Insufficient stock"})
			return
		}

		totalAmount := product.Price * float64(req.Quantity)

		var order model.Order

		order.CustomerId = userId
		order.ProductId = req.ProductId
		order.TotalAmount = totalAmount
		order.TotalQuantity = float64(req.Quantity)
		order.CreatedAt = time.Now()
		order.UpdatedAt = time.Now()

		session, err := client.StartSession()
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "transaction failed"})
			return
		}
		defer session.EndSession(ctx)

		_, err = session.WithTransaction(ctx, func(sessCtx context.Context) (any, error) {
			result, err := orderCollection.InsertOne(sessCtx, order)
			if err != nil {
				return nil, err
			}

			_, err = productCollection.UpdateOne(sessCtx, bson.M{"_id": req.ProductId}, bson.M{
				"$inc": bson.M{"stockQuantity": -req.Quantity},
				"$set": bson.M{"updatedAt": time.Now()},
			})
			if err != nil {
				return nil, err
			}
			return result, nil
		})

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusCreated, gin.H{"message": "order created", "order": order})
	}
}
