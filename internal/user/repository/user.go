package repository

import (
	"context"
	"errors"
	"mile-app-test/domain"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type userRepository struct {
	cl *mongo.Client
}

func NewUserRepository(cl *mongo.Client) domain.UserRepository {
	return &userRepository{cl}
}

func (r *userRepository) GetUser(ctx context.Context, username string) (*domain.User, error) {
	coll := r.cl.Database("local").Collection("users")
	var u domain.User
	if err := coll.FindOne(ctx, bson.M{"username": username}).Decode(&u); err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, nil
		}
		return nil, err
	}
	return &u, nil
}
