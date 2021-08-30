package validators

import (
	"context"
	"errors"
	"fmt"
	"net/mail"
	"regexp"

	"github.com/aliparlakci/country-roads/models"
)

type UserValidator struct {
	Dto *models.NewUserForm
	models.UserFinder
}

func (u *UserValidator) SetDto(dto interface{}) error {
	if d, ok := dto.(*models.NewUserForm); ok {
		u.Dto = d
	} else {
		return fmt.Errorf("dto is not assignable to NewUserForm")
	}
	return nil
}

func ValidateDisplayName(name string) bool {
	return len(name) > 1
}

func ValidatePhone(phone string) bool {
	return regexp.MustCompile(`^(\+|)([0-9]{1,3})([0-9]{10})$`).MatchString(phone)
}

func ValidateEmail(email string) (string, bool) {
	parsedEmail, err := mail.ParseAddress(email)
	if err != nil {
		return "", false
	}

	if regexp.MustCompile(`^.+@sabanciuniv\.edu$`).MatchString(parsedEmail.Address) {
		return parsedEmail.Address, true
	} else {
		return "", false
	}
}

func (u *UserValidator) Validate(ctx context.Context) (bool, error) {
	if isValidName := ValidateDisplayName(u.Dto.DisplayName); !isValidName {
		return false, errors.New("display name is not valid")
	}
	if validEmail, isValidEmail := ValidateEmail(u.Dto.Email); !isValidEmail {
		return false, errors.New("email is not valid")
	} else {
		u.Dto.Email = validEmail
	}
	if isValidPhone := ValidatePhone(u.Dto.Phone); !isValidPhone {
		return false, errors.New("phone is not valid")
	}

	return true, nil
}
