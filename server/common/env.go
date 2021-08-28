package common

import (
	"example.com/country-roads/models"
	"example.com/country-roads/validators"
	"github.com/go-redis/redis"
)

type Env struct {
	Collections *CollectionContainer
	Validators  validators.IValidatorFactory
	Rdb         *redis.Client
}

type CollectionContainer struct {
	RideCollection     models.RideRepository
	LocationCollection models.LocationRepository
}
