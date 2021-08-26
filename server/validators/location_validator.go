package validators

import (
	"context"
	"example.com/country-roads/models"
	"fmt"
)

type LocationValidator struct {
	Dto models.LocationDTO
}

func (v *LocationValidator) SetDto(dto interface{}) {
	if d, ok := dto.(models.LocationDTO); ok {
		v.Dto = d
	} else {
		panic(fmt.Errorf("dto is not assignable to LocationDTO"))
	}
}

func (l LocationValidator) Validate(ctx context.Context) (bool, error) {
	return true, nil
}
