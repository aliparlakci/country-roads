package models

//go:generate mockgen -destination=../mocks/mock_user_model.go -package=mocks example.com/country-roads/models UserRepository,UserFinder,UserInserter,UserUpdater,UserFindUpdater,UserFindInserter

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
	ID          primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	DisplayName string             `bson:"displayName" json:"displayName"`
	Email       string             `bson:"email" json:"email"`
	Phone       string             `bson:"phone" json:"phone"`
	Verified    bool               `bson:"verified" json:"verified"`
	SignedUpAt  time.Time          `bson:"signedUpAt" json:"signedUpAt" time_format:"unix"`
}

type NewUserForm struct {
	DisplayName string `form:"displayName" binding:"required"`
	Email       string `form:"email" binding:"required"`
	Phone       string `form:"phone" binding:"required"`
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

type UserFindInserter interface {
	UserFinder
	UserInserter
}

func (u UserCollection) FindOne(ctx context.Context, filter interface{}) (User, error) {
	var user User
	result := u.Collection.FindOne(ctx, filter)
	if err := result.Err(); err != nil {
		return user, err
	}
	err := result.Decode(&user)
	return user, err
}

func (u UserCollection) InsertOne(ctx context.Context, candidate UserSchema) (interface{}, error) {
	result, err := u.Collection.InsertOne(ctx, candidate)
	if err != nil {
		return nil, err
	}

	return result.InsertedID, nil
}

func (u UserCollection) UpdateOne(ctx context.Context, filter interface{}, changes interface{}) error {
	//TODO: Implement UpdateOne
	return nil
}
