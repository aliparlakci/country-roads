package validators

import (
	"context"
	"example.com/country-roads/models"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type LocationValidator struct {
	Dto            models.LocationDTO
	LocationFinder models.LocationFinder
}

func (v *LocationValidator) SetDto(dto interface{}) {
	if d, ok := dto.(models.LocationDTO); ok {
		v.Dto = d
	} else {
		panic(fmt.Errorf("dto is not assignable to LocationDTO"))
	}
}

func (l LocationValidator) ValidateDisplay() bool {
	return l.Dto.Display != ""
}

func (l LocationValidator) ValidateParent(ctx context.Context) (bool, error) {
	if l.Dto.ParentID == "" {
		return true, nil
	}

	parentId, err := primitive.ObjectIDFromHex(l.Dto.ParentID)
	if err != nil {
		return false, err
	}

	if _, err := l.LocationFinder.FindOne(ctx, bson.M{"_id": parentId}); err != nil {
		return false, err
	}

	return true, nil
}

func (l LocationValidator) Validate(ctx context.Context) (bool, error) {
	if parent, err := l.ValidateParent(ctx); !parent && err != nil {
		return parent, err
	}
	if display := l.ValidateDisplay(); !display {
		return display, nil
	}
	return true, nil
}
