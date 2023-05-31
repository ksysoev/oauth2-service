package repos

import (
	"context"
	"fmt"

	"github.com/ksysoev/oauth2-service/pkg/aggregates"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type UserRepo struct {
	db *mongo.Client
}

func NewUserRepo(db *mongo.Client) *UserRepo {
	return &UserRepo{db: db}
}

func (r *UserRepo) AddUser(ctx context.Context, user *aggregates.User) error {
	err := r.db.Ping(context.Background(), nil)
	if err != nil {
		fmt.Println("Failed to ping MongoDB:", err)
		return err
	}

	_, err = r.db.Database("oauth2").Collection("users").InsertOne(ctx, user)

	if err != nil {
		return err
	}

	return nil
}

func (r *UserRepo) GetUser(ctx context.Context, email string) (*aggregates.User, error) {
	err := r.db.Ping(context.Background(), nil)
	if err != nil {
		fmt.Println("Failed to ping MongoDB:", err)
		return nil, err
	}

	var user aggregates.User
	err = r.db.Database("oauth2").Collection("users").FindOne(ctx, bson.M{"email": email}).Decode(&user)

	if err != nil {
		return nil, err
	}

	return &user, nil
}
