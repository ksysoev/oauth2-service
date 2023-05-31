package services

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"go.mongodb.org/mongo-driver/mongo"

	"github.com/ksysoev/oauth2-service/pkg/aggregates"
	"github.com/ksysoev/oauth2-service/pkg/repos"
)

type SignUpRequest struct {
	Name     string `form:"name" validate:"required,min=2,max=50"`
	Email    string `form:"email" validate:"required,email"`
	Password string `form:"password" validate:"required,min=8,max=50"`
}

func SignUp(c *gin.Context, request *SignUpRequest, mongoClient *mongo.Client) error {

	validate := validator.New()
	if err := validate.Struct(request); err != nil {
		return fmt.Errorf("Invalid user data: %v", err)
	}

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
