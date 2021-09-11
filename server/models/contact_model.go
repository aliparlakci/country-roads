package models

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"regexp"
)

type Contact struct {
	ID       primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Owner    primitive.ObjectID `bson:"owner" json:"owner"`
	Phone    string             `bson:"phone,omitempty" json:"phone,omitempty"`
	Whatsapp string             `bson:"whatsapp,omitempty" json:"whatsapp,omitempty"`
}

type UpdateContactForm struct {
	DisplayName string             `form:"displayName,omitempty"`
	Phone       string             `form:"phone,omitempty"`
	Whatsapp    string             `form:"whatsapp,omitempty"`
}

type ContactResponse struct {
	Name     string `json:"name,omitempty"`
	Email    string `json:"email,omitempty"`
	Phone    string `json:"phone,omitempty"`
	Whatsapp string `json:"whatsapp,omitempty"`
}

func (u UpdateContactForm) validatePhone() bool {
	return regexp.MustCompile(`^(\+|)([0-9]{1,3})([0-9]{10})$`).MatchString(u.Phone)
}

func (u UpdateContactForm) validateWhatsapp() bool {
	return u.validatePhone()
}

func (u UpdateContactForm) validate() (bool, error) {
	result := true
	result = result && u.validatePhone()
	result = result && u.validateWhatsapp()
	return result, nil
}

func (u *UpdateContactForm) Bind(c *gin.Context) error {
	if err := c.Bind(u); err != nil {
		return fmt.Errorf(err.Error())
	}

	if result, err := u.validate(); err != nil {
		return err
	} else if !result {
		return errors.New("UpdateContactForm is not valid")
	}
	return nil
}
