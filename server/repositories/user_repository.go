package repositories

//go:generate mockgen -destination=../mocks/mock_user_repository.go -package=mocks github.com/aliparlakci/country-roads/repositories UserRepository,UserFinder,UserInserter,UserUpdater,UserFindUpdater,UserFindInserter

import (
	"context"
	"github.com/aliparlakci/country-roads/models"
	"go.mongodb.org/mongo-driver/mongo"
)

type UserCollection struct {
	Collection *mongo.Collection
}

type UserRepository interface {
	UserFinder
	UserInserter
	UserUpdater
}

type UserFinder interface {
	FindOne(ctx context.Context, filter interface{}) (models.User, error)
}

type UserInserter interface {
	InsertOne(ctx context.Context, candidate models.UserSchema) (interface{}, error)
}

type UserUpdater interface {
	UpdateOne(ctx context.Context, filter interface{}, changes interface{}) error
}

type UserFindUpdater interface {
	UserFinder
	UserUpdater
}

type UserFindInserter interface {
	UserFinder
	UserInserter
}

func (c *UserCollection) FindOne(ctx context.Context, filter interface{}) (models.User, error) {
	var user models.User
	result := c.Collection.FindOne(ctx, filter)
	if err := result.Err(); err != nil {
		return user, err
	}
	err := result.Decode(&user)
	return user, err
}

func (c *UserCollection) InsertOne(ctx context.Context, candidate models.UserSchema) (interface{}, error) {
	result, err := c.Collection.InsertOne(ctx, candidate)
	if err != nil {
		return nil, err
	}

	return result.InsertedID, nil
}

func (c *UserCollection) UpdateOne(ctx context.Context, filter interface{}, changes interface{}) error {
	_, err := c.Collection.UpdateOne(ctx, filter, changes)
	return err
}
