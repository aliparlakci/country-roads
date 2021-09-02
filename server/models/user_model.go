package models

//go:generate mockgen -destination=../mocks/mock_user_model.go -package=mocks github.com/aliparlakci/country-roads/models UserRepository,UserFinder,UserInserter,UserUpdater,UserFindUpdater,UserFindInserter

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/mail"
	"regexp"
	"time"

	"go.mongodb.org/mongo-driver/mongo"

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
	DisplayName string `form:"displayName" json:"displayName" binding:"required"`
	Email       string `form:"email" json:"email" binding:"required"`
	Phone       string `form:"phone" json:"phone" binding:"required"`
}

type LoginRequestForm struct {
	Email string `form:"email" binding:"required"`
}

type VerifyRequestForm struct {
	Email string `form:"email" binding:"required"`
	OTP string `form:"otp" binding:"required"`
}

func (n NewUserForm) validateDisplayName() bool {
	return len(n.DisplayName) > 1
}

func (n NewUserForm) validatePhone() bool {
	return regexp.MustCompile(`^(\+|)([0-9]{1,3})([0-9]{10})$`).MatchString(n.Phone)
}

func (n NewUserForm) validateEmail() (string, bool) {
	parsedEmail, err := mail.ParseAddress(n.Email)
	if err != nil {
		return "", false
	}

	if regexp.MustCompile(`^.+@sabanciuniv\.edu$`).MatchString(parsedEmail.Address) {
		return parsedEmail.Address, true
	} else {
		return "", false
	}
}

func (n *NewUserForm) validate() (bool, error) {
	if isValidName := n.validateDisplayName(); !isValidName {
		return false, fmt.Errorf("display name is not valid")
	}
	if validEmail, isValidEmail := n.validateEmail(); !isValidEmail {
		return false, fmt.Errorf("email is not valid")
	} else {
		n.Email = validEmail
	}
	if isValidPhone := n.validatePhone(); !isValidPhone {
		return false, fmt.Errorf("phone is not valid")
	}

	return true, nil
}

func (n *NewUserForm) Bind(c *gin.Context) error {
	if err := c.Bind(n); err != nil {
		return fmt.Errorf(err.Error())
	}

	if result, err := n.validate(); err != nil {
		return err
	} else if !result {
		return fmt.Errorf("")
	}
	return nil
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

func (u *UserCollection) FindOne(ctx context.Context, filter interface{}) (User, error) {
	var user User
	result := u.Collection.FindOne(ctx, filter)
	if err := result.Err(); err != nil {
		return user, err
	}
	err := result.Decode(&user)
	return user, err
}

func (u *UserCollection) InsertOne(ctx context.Context, candidate UserSchema) (interface{}, error) {
	result, err := u.Collection.InsertOne(ctx, candidate)
	if err != nil {
		return nil, err
	}

	return result.InsertedID, nil
}

func (u *UserCollection) UpdateOne(ctx context.Context, filter interface{}, changes interface{}) error {
	//TODO: Implement UpdateOne
	return nil
}
