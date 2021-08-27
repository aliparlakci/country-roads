package validators

import (
	"context"
	"example.com/country-roads/models"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type LocationValidator struct {
	Dto models.LocationDTO
	LocationFinder models.LocationFinder
}

func (v *LocationValidator) SetDto(dto interface{}) {
	if d, ok := dto.(models.LocationDTO); ok {
		v.Dto = d
	} else {
		panic(fmt.Errorf("dto is not assignable to LocationDTO"))
	}
}

func (l LocationValidator) ValidateParent(ctx context.Context) bool {
	if l.Dto.ParentID == "" {
		return true
	}

	parentId, err := primitive.ObjectIDFromHex(l.Dto.ParentID)
	if err != nil {
		return false
	}

	if _, err := l.LocationFinder.FindOne(ctx, bson.M{"_id": parentId}); err != nil {
		return false
	}

	return true
}

func (l LocationValidator) Validate(ctx context.Context) (bool, error) {
	return true, nil
}
