package services

import (
	"context"
	"fmt"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
)

type SignUpRequest struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

func SignUp(c *gin.Context, request *SignUpRequest, mongoClient *mongo.Client) error {
	err := mongoClient.Ping(context.Background(), nil)
	if err != nil {
		fmt.Println("Failed to ping MongoDB:", err)
		return err
	}

	_, err = mongoClient.Database("oauth2").Collection("users").InsertOne(c, request)

	if err != nil {
		return err
	}

	return nil
}
