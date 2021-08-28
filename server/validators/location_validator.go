package validators

import (
	"context"
	"example.com/country-roads/models"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
)

type LocationValidator struct {
	Dto            models.NewLocationFrom
	LocationFinder models.LocationFinder
}

func (v *LocationValidator) SetDto(dto interface{}) {
	if d, ok := dto.(models.NewLocationFrom); ok {
		v.Dto = d
	} else {
		panic(fmt.Errorf("dto is not assignable to NewLocationFrom"))
	}
}

func (l LocationValidator) validateDisplay() bool {
	return l.Dto.Display != ""
}

func (l LocationValidator) validateParent(ctx context.Context) (bool, error) {
	if l.Dto.ParentKey == "" {
		return true, nil
	}

	if _, err := l.LocationFinder.FindOne(ctx, bson.M{"key": l.Dto.ParentKey}); err != nil {
		return false, err
	}

	return true, nil
}

func (l LocationValidator) validateKey(ctx context.Context) bool {
	if _, err := l.LocationFinder.FindOne(ctx, bson.M{"key": l.Dto.Key}); err != nil {
		return true
	}
	return false
}

func (l LocationValidator) Validate(ctx context.Context) (bool, error) {
	if parent, err := l.validateParent(ctx); !parent && err != nil {
		return parent, err
	}
	if display := l.validateDisplay(); !display {
		return display, nil
	}
	if key := l.validateKey(ctx); !key {
		return key, nil
	}
	return true, nil
}
