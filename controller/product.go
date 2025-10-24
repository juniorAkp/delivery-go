package controller

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/juniorAkp/delivery-go/database"
	"github.com/juniorAkp/delivery-go/model"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

var validate = validator.New()

func GetProducts(client *mongo.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(c, 100*time.Second)
		defer cancel()

		var products []model.Product

		var productCollection *mongo.Collection = database.OpenCollection("products", client)

		cur, err := productCollection.Find(ctx, bson.M{})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err})
			return
		}
		defer cur.Close(ctx)

		if err := cur.All(ctx, &products); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err})
			return
		}
		c.JSON(http.StatusOK, products)
	}
}

func GetProduct(client *mongo.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(c, 100*time.Second)
		defer cancel()

		productID := c.Param("productID")

		if productID == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "product required"})
			return
		}

		var product model.Product

		var productCollection *mongo.Collection = database.OpenCollection("products", client)

		err := productCollection.FindOne(ctx, bson.M{"ID": productID}).Decode(&product)

		if errors.Is(err, mongo.ErrNoDocuments) {
			c.JSON(http.StatusNotFound, gin.H{"error": fmt.Sprintf("no product matching this id %s found", productID)})
			return
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"err": err})
		}
		c.JSON(http.StatusOK, product)
	}
}

func CreateProduct(client *mongo.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(c, 100*time.Second)
		defer cancel()

		var product model.Product

		var productCollection *mongo.Collection = database.OpenCollection("products", client)

		if err := c.ShouldBindJSON(&product); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err})
			return
		}

		if err := validate.Struct(product); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "validation failed", "details": err})
			return
		}

		product.CreatedAt = time.Now()
		product.UpdatedAt = time.Now()
		product.ProductId = bson.NewObjectID().Hex()

		result, err := productCollection.InsertOne(ctx, product)

		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err})
			return
		}

		c.JSON(http.StatusCreated, result)
	}
}
