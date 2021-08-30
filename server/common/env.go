package common

import (
	"github.com/aliparlakci/country-roads/models"
	"github.com/aliparlakci/country-roads/validators"
	"github.com/go-redis/redis"
)

type Env struct {
	Repositories     *RepositoryContainer
	ValidatorFactory validators.IValidatorFactory
	Rdb              *redis.Client
}

type RepositoryContainer struct {
	models.RideRepository
	models.LocationRepository
	models.UserRepository
}
