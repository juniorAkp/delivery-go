package utils

import (
	"context"
	"errors"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/juniorAkp/delivery-go/database"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type SignedClaims struct {
	Sub   string
	Email string
	jwt.RegisteredClaims
}

type TokenType string

const (
	AccessToken  TokenType = "accessToken"
	RefreshToken TokenType = "refreshToken"
)

var AccessSecretKey string = os.Getenv("ACCESS_SECRET_KEY")
var RefreshSecretKey string = os.Getenv("REFRESH_SECRET_KEY")

func GenerateAllTokens(id, email string) (signedToken, signedRefreshToken string, err error) {
	claims := &SignedClaims{
		Email: email,
		Sub:   id,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    "delivery-go",
			Subject:   "accessToken",
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err = token.SignedString([]byte(AccessSecretKey))
	if err != nil {
		return "", "", err
	}

	refreshClaims := &SignedClaims{
		Email: email,
		Sub:   id,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    "delivery-go",
			Subject:   "refreshToken",
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour * 7)),
		},
	}
	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims)
	signedRefreshToken, err = refreshToken.SignedString([]byte(RefreshSecretKey))
	if err != nil {
		return "", "", err
	}
	return signedToken, signedRefreshToken, nil
}

func UpdateAllTokens(userId, token, refreshToken string, client *mongo.Client) error {
	var ctx, cancel = context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	updateAt, _ := time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))

	updateData := bson.M{
		"$set": bson.M{
			"token":        token,
			"refreshToken": refreshToken,
			"updatedAt":    updateAt,
		},
	}
	var customersCollection *mongo.Collection = database.OpenCollection("Customers", client)

	_, err := customersCollection.UpdateOne(ctx, bson.M{"_id": userId}, updateData)
	if err != nil {
		return err
	}

	return nil
}

func ValidateToken(tokenString string, tokenType TokenType) (*SignedClaims, error) {
	claims := &SignedClaims{}

	var secretKey string

	switch tokenType {
	case AccessToken:
		secretKey = AccessSecretKey
	case RefreshToken:
		secretKey = RefreshSecretKey
	default:
		return nil, errors.New("invalid token type")
	}

	if secretKey == "" {
		return nil, errors.New("secret key not found")
	}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (any, error) {
		return []byte(secretKey), nil
	})
	if err != nil {
		return nil, err
	}

	if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
		return nil, err
	}

	if !token.Valid {
		return nil, errors.New("invalid token")
	}

	if claims.ExpiresAt.Time.Before(time.Now()) {
		return nil, errors.New("token is expired")
	}

	return claims, nil
}

func GetUserId(c *gin.Context) (string, error) {
	userId, exists := c.Get("userId")
	if !exists {
		return "", errors.New("userId not found in context")
	}
	id, ok := userId.(string)
	if !ok {
		return "", errors.New("userId has an invalid type")
	}
	if id == "" {
		return "", errors.New("userId is empty")
	}
	return id, nil
}
