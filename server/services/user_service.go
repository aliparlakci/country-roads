package services

import (
	"context"
	"errors"
	"fmt"
	"github.com/aliparlakci/country-roads/models"
	"github.com/aliparlakci/country-roads/repositories"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"regexp"
	"time"
)

var (
	ErrUserAlreadyExist = errors.New("user already exists")
)

type UserServiceStruct struct {
	Repo repositories.UserRepository
}

type UserService interface {
	CreateUser(c context.Context, email string) (interface{}, error)
	FindUser(c context.Context, email string) (models.User, error)
	Exists(c context.Context, email string) (bool, error)
	FindUserById(c context.Context, id string) (models.User, error)
}

func (u UserServiceStruct) CreateUser(c context.Context, email string) (interface{}, error) {
	if _, err := u.Repo.FindOne(c, bson.M{"email": email}); err == nil {
		return nil, ErrUserAlreadyExist
	}

	name := regexp.
		MustCompile(`(.*)@sabanciuniv.edu`).
		FindAllStringSubmatch(email, -1)[0][1]

	return u.Repo.InsertOne(c, models.UserSchema{
		DisplayName: name,
		Email:       email,
		Verified:    false,
		SignedUpAt:  time.Now(),
	})
}

func (u UserServiceStruct) FindUserById(c context.Context, id string) (models.User, error) {
	return u.Repo.FindOne(c, bson.M{"_id": id})
}

func (u UserServiceStruct) FindUser(c context.Context, email string) (models.User, error) {
	return u.Repo.FindOne(c, bson.M{"email": email})
}

func (u UserServiceStruct) Exists(c context.Context, email string) (bool, error) {
	if _, err := u.Repo.FindOne(c, bson.M{"email": email}); err == mongo.ErrNoDocuments {
		return false, nil
	} else if err != nil {
		return false, fmt.Errorf("UserRepository.FindOne raised an error: %v", err)
	}
	return true, nil
}
