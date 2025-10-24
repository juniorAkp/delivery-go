package controller

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/juniorAkp/delivery-go/database"
	"github.com/juniorAkp/delivery-go/model"
	"github.com/juniorAkp/delivery-go/utils"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"golang.org/x/crypto/bcrypt"
)

func RegisterCustomer(client *mongo.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		var customer model.Customer

		if err := c.ShouldBindJSON(&customer); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "bind error"})
			return
		}

		var validate = validator.New()

		if err := validate.Struct(customer); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		hashedPassword := HashPassword(customer.Password)

		var ctx, cancel = context.WithTimeout(c, 100*time.Second)
		defer cancel()

		var customerCollection *mongo.Collection = database.OpenCollection("customers", client)

		count, err := customerCollection.CountDocuments(ctx, bson.M{"email": customer.Email})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to check user"})
			return
		}
		if count > 0 {
			c.JSON(http.StatusConflict, gin.H{"error": "User already exists"})
			return
		}

		customer.UserId = bson.NewObjectID().Hex()
		customer.Password = hashedPassword
		customer.CreatedAt = time.Now()
		customer.UpdatedAt = time.Now()

		result, err := customerCollection.InsertOne(ctx, customer)

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		c.JSON(http.StatusCreated, result)
	}
}

func LoginCustomer(client *mongo.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		var customerLogin model.CustomerLogin

		if err := c.BindJSON(&customerLogin); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err})
			return
		}

		var ctx, cancel = context.WithTimeout(c, 100*time.Second)
		defer cancel()

		var foundCustomer model.Customer

		var userCollection *mongo.Collection = database.OpenCollection("customers", client)
		err := userCollection.FindOne(ctx, bson.M{"email": customerLogin.Email}).Decode(&foundCustomer)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "invalid credentials"})
			return
		}

		isPasswordValid, msg := CheckPasswordHash(customerLogin.Password, foundCustomer.Password)
		if !isPasswordValid {
			c.JSON(http.StatusUnauthorized, gin.H{"error": msg})
			return
		}
		token, refreshToken, err := utils.GenerateAllTokens(foundCustomer.UserId, foundCustomer.Email)

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		err = utils.UpdateAllTokens(foundCustomer.UserId, token, refreshToken, client)

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}

		c.JSON(http.StatusOK, gin.H{"accessToken": token, "refreshToken": refreshToken})
	}
}

func HashPassword(password string) string {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	if err != nil {
		log.Panic(err)
	}
	return string(bytes)
}

func CheckPasswordHash(password, hash string) (bool, string) {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	check := true
	msg := ""
	if err != nil {
		msg = "invalid credentials"
		check = false
	}
	return check, msg
}
