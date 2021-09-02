package common

import (
	"github.com/aliparlakci/country-roads/repositories"
	"github.com/aliparlakci/country-roads/services"
)

type Env struct {
	Repositories     *RepositoryContainer
	Services         *ServiceContainer
}

type RepositoryContainer struct {
	repositories.RideRepository
	repositories.LocationRepository
	repositories.UserRepository
}

type ServiceContainer struct {
	services.SessionService
	services.OTPService
}
