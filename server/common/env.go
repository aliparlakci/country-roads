package common

import (
	"github.com/aliparlakci/country-roads/models"
	"github.com/aliparlakci/country-roads/services"
	"github.com/aliparlakci/country-roads/validators"
)

type Env struct {
	Repositories     *RepositoryContainer
	Services         *ServiceContainer
	ValidatorFactory validators.IValidatorFactory
}

type RepositoryContainer struct {
	models.RideRepository
	models.LocationRepository
	models.UserRepository
}

type ServiceContainer struct {
	services.SessionService
	services.OTPService
}
