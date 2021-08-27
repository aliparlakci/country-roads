package validators

import "context"

//go:generate mockgen -destination=../mocks/mock_ride_validator.go -package=mocks example.com/country-roads/validators Validator

type Validator interface {
	SetDto(dto interface{})
	Validate(ctx context.Context) (bool, error)
}
