package services

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/ksysoev/oauth2-service/pkg/repos"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
)

type SignInRequest struct {
	Email    string `form:"email" validate:"required,email"`
	Password string `form:"password" validate:"required,min=8,max=50"`
}

func SignIn(c *gin.Context, request *SignInRequest, mongoClient *mongo.Client) error {
	validate := validator.New()
	if err := validate.Struct(request); err != nil {
		return fmt.Errorf("Invalid user data: %v", err)
	}

	UserRepo := repos.NewUserRepo(mongoClient)

	user, err := UserRepo.GetUser(c, request.Email)

	if err != nil {
		return err
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(request.Password))

	if err != nil {
		return err
	}

	return nil
}
