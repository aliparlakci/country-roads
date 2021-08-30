package models

//go:generate mockgen -destination=../mocks/mock_user_model.go -package=mocks example.com/country-roads/models UserRepository,UserFinder,UserInserter,UserUpdater,UserFindUpdater

import (
	"context"
	"go.mongodb.org/mongo-driver/mongo"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
	ID          primitive.ObjectID `bson:"_id" json:"id"`
	DisplayName string             `bson:"displayName" json:"displayName"`
	Email       string             `bson:"email" json:"email"`
	Phone       string             `bson:"phone" json:"phone"`
	Verified    bool               `bson:"verified" json:"verified"`
	SignedUpAt  time.Time          `bson:"signedUpAt" json:"signedUpAt" time_format:"unix"`
}

type UserSchema struct {
	ID          primitive.ObjectID `bson:"_id" json:"id"`
	DisplayName string             `bson:"displayName" json:"displayName"`
	Email       string             `bson:"email" json:"email"`
	Phone       string             `bson:"phone" json:"phone"`
	Verified    bool               `bson:"verified" json:"verified"`
	SignedUpAt  time.Time          `bson:"signedUpAt" json:"signedUpAt" time_format:"unix"`
}

type NewUserForm struct {
	DisplayName string
	Email       string
	Phone       string
}

type UserCollection struct {
	Collection *mongo.Collection
}

type UserRepository interface {
	UserFinder
	UserInserter
	UserUpdater
}

type UserFinder interface {
	FindOne(ctx context.Context, filter interface{}) (User, error)
}

type UserInserter interface {
	InsertOne(ctx context.Context, candidate UserSchema) (interface{}, error)
}

type UserUpdater interface {
	UpdateOne(ctx context.Context, filter interface{}, changes interface{}) error
}

type UserFindUpdater interface {
	UserFinder
	UserUpdater
}

func (u UserCollection) FindOne(ctx context.Context, filter interface{}) (User, error) {
	return User{}, nil
}

func (u UserCollection) InsertOne(ctx context.Context, candidate UserSchema) (interface{}, error) {
	return User{}, nil
}

func (u UserCollection) UpdateOne(ctx context.Context, filter interface{}, changes interface{}) error {
	return nil
}
