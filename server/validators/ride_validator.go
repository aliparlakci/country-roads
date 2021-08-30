package validators

import (
	"context"
	"fmt"
	"time"

	"github.com/aliparlakci/country-roads/models"
	"go.mongodb.org/mongo-driver/bson"
)

type RideValidator struct {
	Dto            models.NewRideForm
	LocationFinder models.LocationFinder
}

func (v *RideValidator) SetDto(dto interface{}) error {
	switch t := dto.(type) {
	case models.NewRideForm:
		v.Dto = dto.(models.NewRideForm)
	default:
		return fmt.Errorf("%v is not assignable to NewRideForm", t)
	}
	return nil
}

func (v RideValidator) ValidateDate() bool {
	date := v.Dto.Date
	today := time.Now()
	if today.Year() < date.Year() {
		return true
	} else if today.Year() == date.Year() {
		if today.Month() < date.Month() {
			return true
		} else if today.Month() == date.Month() {
			if today.Day() <= date.Day() {
				return true
			}
		}
	}
	return false
}

func (v RideValidator) ValidateType() bool {
	switch v.Dto.Type {
	case "offer":
		return true
	case "request":
		return true
	case "taxi":
		return true
	default:
		return false
	}
}

func (v RideValidator) ValidateDirection() bool {
	switch v.Dto.Direction {
	case "to_campus":
		return true
	case "from_campus":
		return true
	default:
		return false
	}
}

func (v RideValidator) ValidateDestination(ctx context.Context) bool {
	if _, err := v.LocationFinder.FindOne(ctx, bson.M{"key": v.Dto.Destination}); err != nil {
		return false
	}
	return true
}

func (v RideValidator) Validate(ctx context.Context) (bool, error) {
	if !v.ValidateDate() {
		return false, fmt.Errorf("ride date is not valid")
	}
	if !v.ValidateDirection() {
		return false, fmt.Errorf("ride direction is not valid")
	}
	if !v.ValidateType() {
		return false, fmt.Errorf("ride type is not valid")
	}
	if !v.ValidateDestination(ctx) {
		return false, fmt.Errorf("ride destination is not valid")
	}

	return true, nil
}
