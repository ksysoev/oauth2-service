package user

import (
	"context"

	"go.mongodb.org/mongo-driver/mongo"
)

type UserRepo struct {
	mongoClient *mongo.Client
}

func (u *UserRepo) AddUser(ctx context.Context, user User) (string, error) {

}

func (u *UserRepo) GetUser(ctx context.Context, id string) (User, error) {

}
