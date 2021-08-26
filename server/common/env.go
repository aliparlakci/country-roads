package common

import (
	"example.com/country-roads/interfaces"
	"example.com/country-roads/models"
	"github.com/go-redis/redis"
)

type Env struct {
	Collections CollectionContainer
	Validators ValidatorContainer
	Rdb                *redis.Client
}

type CollectionContainer struct {
	RideCollection     models.RideRepository
	LocationCollection models.LocationRepository
}

type ValidatorContainer struct {
	RideValidator func() interfaces.Validator
	LocationValidator func() interfaces.Validator
}
