package main

import (
	"os"

	"github.com/aliparlakci/country-roads/middlewares"
	"github.com/aliparlakci/country-roads/repositories"
	"github.com/aliparlakci/country-roads/services"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"

	"github.com/aliparlakci/country-roads/common"
	"github.com/aliparlakci/country-roads/controllers"
	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load()

	logrus.SetLevel(logrus.InfoLevel)
	logrus.SetReportCaller(true)
	logrus.SetFormatter(&logrus.TextFormatter{PadLevelText: true})
	//logrus.SetFormatter(&logrus.JSONFormatter{PrettyPrint: true})

	env := &common.Env{}
	{
		db, close := common.InitializeDb(
			os.Getenv("MDB_URI"),
			os.Getenv("MDB_DBNAME"),
			os.Getenv("MDB_USERNAME"),
			os.Getenv("MDB_PASSWORD"),
		)
		defer close()

		rdbUri := os.Getenv("RDB_URI")
		redis := common.RedisInitilizer(rdbUri, os.Getenv("RDB_PASSWORD"))
		env.Repositories = &common.RepositoryContainer{
			RideRepository:     &repositories.RideCollection{Collection: db.Collection("rides")},
			LocationRepository: &repositories.LocationCollection{Collection: db.Collection("locations")},
			UserRepository:     &repositories.UserCollection{Collection: db.Collection("users")},
		}
		env.Services = &common.ServiceContainer{
			SessionService: &services.SessionStore{Store: redis(0)},
			OTPService:     &services.OTPStore{Store: redis(1)},
		}
	}

	router := gin.Default()
	router.Use(func(c *gin.Context) {
		c.Set("request_id", uuid.New().String())
		c.Next()
	})
	router.Use(middlewares.Logger())
	router.Use(cors.New(cors.Config{
		AllowAllOrigins: true,
		AllowMethods:    []string{"GET", "POST", "DELETE"},
	}))

	router.Use(middlewares.AuthMiddleware(
		env.Repositories.UserRepository,
		env.Services.SessionService,
	))
	api := router.Group("api")
	{
		controllers.RegisterRideController(api, env)
		controllers.RegisterLocationController(api, env)
		controllers.RegisterUserController(api, env)
		controllers.RegisterAuthController(api, env)
	}

	router.Run(":5000")
}
