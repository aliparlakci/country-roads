package validators

import (
	"context"
	"fmt"

	"example.com/country-roads/models"
)

//go:generate mockgen -destination=../mocks/mock_ride_validator.go -package=mocks example.com/country-roads/validators Validator,IValidatorFactory

type Validator interface {
	SetDto(dto interface{}) error
	Validate(ctx context.Context) (bool, error)
}

type IValidatorFactory interface {
	GetValidator(name string) (Validator, error)
}

type ValidatorFactory struct {
	models.LocationFinder
}

func (v ValidatorFactory) GetValidator(name string) (Validator, error) {
	switch name {
	case "rides":
		return &RideValidator{LocationFinder: v.LocationFinder}, nil
	case "locations":
		return &LocationValidator{LocationFinder: v.LocationFinder}, nil
	case "users":
		return &UserValidator{}, nil
	default:
		return nil, fmt.Errorf("no validator of type %v", name)
	}
}
