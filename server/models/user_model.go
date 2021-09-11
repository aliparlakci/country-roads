package models

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/mail"
	"regexp"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
	ID               primitive.ObjectID `bson:"_id" json:"id"`
	DisplayName      string             `bson:"displayName" json:"displayName"`
	Email            string             `bson:"email" json:"email"`
	Phone            string             `bson:"phone" json:"phone"`
	Verified         bool               `bson:"verified" json:"verified"`
	SignedUpAt       time.Time          `bson:"signedUpAt" json:"signedUpAt" time_format:"unix"`
	ContactWithPhone bool               `bson:"contactWithPhone" json:"contactWithPhone"`
}

type UserResponse struct {
	ID          primitive.ObjectID `bson:"_id" json:"id"`
	DisplayName string             `bson:"displayName" json:"displayName"`
}

type UserSchema struct {
	ID               primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	DisplayName      string             `bson:"displayName" json:"displayName"`
	Email            string             `bson:"email" json:"email"`
	Verified         bool               `bson:"verified" json:"verified"`
	SignedUpAt       time.Time          `bson:"signedUpAt" json:"signedUpAt" time_format:"unix"`
}

type NewUserForm struct {
	DisplayName      string `form:"displayName" json:"displayName" binding:"required"`
	Email            string `form:"email" json:"email" binding:"required"`
}

type SignInRequestForm struct {
	Email string `form:"email" binding:"required"`
}

type VerifyRequestForm struct {
	Email string `form:"email" binding:"required"`
	OTP   string `form:"otp" binding:"required"`
}

func (u NewUserForm) validateDisplayName() bool {
	return len(u.DisplayName) > 1
}

func (u NewUserForm) validateEmail() (string, bool) {
	parsedEmail, err := mail.ParseAddress(u.Email)
	if err != nil {
		return "", false
	}

	if regexp.MustCompile(`^.+@sabanciuniv\.edu$`).MatchString(parsedEmail.Address) {
		return parsedEmail.Address, true
	} else {
		return "", false
	}
}

func (u *NewUserForm) validate() (bool, error) {
	if isValidName := u.validateDisplayName(); !isValidName {
		return false, fmt.Errorf("display name is not valid")
	}
	if validEmail, isValidEmail := u.validateEmail(); !isValidEmail {
		return false, fmt.Errorf("email is not valid")
	} else {
		u.Email = validEmail
	}

	return true, nil
}

func (u *NewUserForm) Bind(c *gin.Context) error {
	if err := c.Bind(u); err != nil {
		return fmt.Errorf(err.Error())
	}

	if result, err := u.validate(); err != nil {
		return err
	} else if !result {
		return fmt.Errorf("")
	}
	return nil
}

func (u UserResponse) Jsonify() map[string]interface{} {
	return map[string]interface{}{
		"id":          u.ID,
		"displayName": u.DisplayName,
	}
}
