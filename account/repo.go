package account

import (
	"context"
	"errors"

	"github.com/go-kit/kit/log"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/bson"
)

var RepoErr = errors.New("Unable to handle Repo Request")

type repo struct {
	db     *mongo.Database
	logger log.Logger
}

func NewRepo(db *mongo.Database, logger log.Logger) Repository {
	return &repo{
		db: db,
		logger: log.With(logger, "repo", "mongodb"),
	}
}

func (repo *repo) CreateUser(ctx context.Context, user User) error {
	if user.Email == "" || user.Password == "" {
		return RepoErr
	}

	_, err := repo.db.Collection("users").InsertOne(ctx, User{user.ID, user.Email, user.Password})
	if err != nil {
		return err
	}
	return nil
}

func (repo *repo) GetUser(ctx context.Context, id string) (string, error) {
	var user User
	err := repo.db.Collection("users").FindOne(ctx, bson.M{"id": id}).Decode(&user)
	if err != nil {
		return "", RepoErr
	}

	return user.Email, nil
}
