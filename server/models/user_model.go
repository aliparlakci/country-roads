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

type ContactResponse struct {
	Phone string `json:"phone,omitempty"`
	Email string `json:"email,omitempty"`
}

type UserSchema struct {
	ID               primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	DisplayName      string             `bson:"displayName" json:"displayName"`
	Email            string             `bson:"email" json:"email"`
	Phone            string             `bson:"phone" json:"phone"`
	Verified         bool               `bson:"verified" json:"verified"`
	SignedUpAt       time.Time          `bson:"signedUpAt" json:"signedUpAt" time_format:"unix"`
	ContactWithPhone bool               `bson:"contactWithPhone" json:"contactWithPhone"`
}

type NewUserForm struct {
	DisplayName      string `form:"displayName" json:"displayName" binding:"required"`
	Email            string `form:"email" json:"email" binding:"required"`
	Phone            string `form:"phone" json:"phone" binding:"required"`
	ContactWithPhone bool   `form:"contact_with_phone" json:"contact_with_phone"`
}

type LoginRequestForm struct {
	Email string `form:"email" binding:"required"`
}

type VerifyRequestForm struct {
	Email string `form:"email" binding:"required"`
	OTP   string `form:"otp" binding:"required"`
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

func (u UserResponse) Jsonify() map[string]interface{} {
	return map[string]interface{}{
		"id":          u.ID,
		"displayName": u.DisplayName,
	}
}
