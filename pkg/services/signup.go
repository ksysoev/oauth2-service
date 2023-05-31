package services

import (
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"

	"github.com/ksysoev/oauth2-service/pkg/aggregates"
	"github.com/ksysoev/oauth2-service/pkg/repos"
)

type SignUpRequest struct {
	Name     string `form:"name"`
	Email    string `form:"email"`
	Password string `form:"password"`
}

func SignUp(c *gin.Context, request *SignUpRequest, mongoClient *mongo.Client) error {

	user, err := aggregates.CreateUser(request.Email, request.Name, request.Password)

	if err != nil {
		return err
	}

	UserRepo := repos.NewUserRepo(mongoClient)

	err = UserRepo.AddUser(c, user)

	if err != nil {
		return err
	}

	return nil
}
