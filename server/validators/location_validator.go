package validators

import (
	"context"
	"example.com/country-roads/models"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
)

type LocationValidator struct {
	Dto            models.NewLocationForm
	LocationFinder models.LocationFinder
}

func (v *LocationValidator) SetDto(dto interface{}) error {
	if d, ok := dto.(models.NewLocationForm); ok {
		v.Dto = d
	} else {
		return fmt.Errorf("dto is not assignable to NewLocationForm")
	}
	return nil
}

func (l LocationValidator) ValidateDisplay() bool {
	return l.Dto.Display != ""
}

func (l LocationValidator) ValidateParent(ctx context.Context) (bool, error) {
	if l.Dto.ParentKey == "" {
		return true, nil
	}

	if _, err := l.LocationFinder.FindOne(ctx, bson.M{"key": l.Dto.ParentKey}); err != nil {
		return false, err
	}

	return true, nil
}

func (l LocationValidator) ValidateKey(ctx context.Context) bool {
	if _, err := l.LocationFinder.FindOne(ctx, bson.M{"key": l.Dto.Key}); err != nil {
		return true
	}
	return false
}

func (l LocationValidator) Validate(ctx context.Context) (bool, error) {
	if parent, err := l.ValidateParent(ctx); !parent && err != nil {
		return parent, err
	}
	if display := l.ValidateDisplay(); !display {
		return display, nil
	}
	if key := l.ValidateKey(ctx); !key {
		return key, nil
	}
	return true, nil
}
