package controllers

import (
	"fmt"
	"net/http"
	"time"

	"example.com/country-roads/aggregations"
	"example.com/country-roads/schemas"
	"example.com/country-roads/validators"
	"go.mongodb.org/mongo-driver/bson"

	"example.com/country-roads/common"
	"example.com/country-roads/models"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func GetRide(env *common.Env) gin.HandlerFunc {
	return func(c *gin.Context) {
		rawId := c.Param("id")
		objID, err := primitive.ObjectIDFromHex(rawId)
		if err != nil {
			c.String(http.StatusBadRequest, err.Error())
			return
		}

		var repo models.RideRepository = env.Collections.RideCollection

		ride, err := repo.FindOne(c, bson.M{"_id": objID})
		if err != nil {
			c.String(http.StatusInternalServerError, err.Error())
			return
		}

		c.JSON(http.StatusOK, ride.Jsonify())
	}
}

func GetAllRides(env *common.Env) gin.HandlerFunc {
	return func(c *gin.Context) {
		var results []map[string]interface{}

		var finder models.RideFinder = env.Collections.RideCollection
		if rides, err := finder.FindMany(c, aggregations.RideWithDestination); err != nil {
			c.String(http.StatusInternalServerError, err.Error())
			return
		} else {
			for _, ride := range rides {
				results = append(results, ride.Jsonify())
			}
		}

		c.JSON(http.StatusOK, results)
	}
}

func PostRides(env *common.Env) gin.HandlerFunc {
	return func(c *gin.Context) {
		var rideDto models.RideDTO

		if err := c.Bind(&rideDto); err != nil {
			c.JSON(http.StatusBadRequest, fmt.Sprintf("Ride format was incorrect: %v", err))
			return
		}

		destinationId, err := primitive.ObjectIDFromHex(rideDto.Destination)
		if err != nil {
			c.JSON(http.StatusBadRequest, fmt.Sprintf("Ride format was incorrect: %v", err))
			return
		}

		var validator validators.Validator = env.Validators.RideValidator()
		validator.SetDto(rideDto)
		if isValid, err := validator.Validate(c); !isValid || err != nil {
			c.JSON(http.StatusBadRequest, fmt.Sprintf("Ride format was invalid: %v", err))
		}

		var inserter models.RideInserter = env.Collections.RideCollection
		id, err := inserter.InsertOne(c, schemas.RideSchema{
			Type:        rideDto.Type,
			Date:        rideDto.Date,
			Destination: destinationId,
			Direction:   rideDto.Direction,
			CreatedAt:   time.Now(),
		})
		if err != nil {
			c.String(http.StatusInternalServerError, fmt.Sprintf("Ride couldn't get created: %v", err))
			return
		}

		c.JSON(http.StatusCreated, id)
	}
}

func DeleteRides(env *common.Env) gin.HandlerFunc {
	return func(c *gin.Context) {
		objID, err := primitive.ObjectIDFromHex(c.Param("id"))
		if err != nil {
			c.String(http.StatusInternalServerError, err.Error())
			return
		}

		var deleter models.RideDeleter = env.Collections.RideCollection
		deletedCount, err := deleter.DeleteOne(c, bson.M{"_id": objID})
		if err != nil {
			c.String(http.StatusInternalServerError, err.Error())
			return
		}

		if deletedCount == 0 {
			c.JSON(http.StatusNotFound, fmt.Sprintf("Ride with ID %v does not exist", objID))
		}

		c.JSON(http.StatusOK, "")
	}
}

func RegisterRideController(router *gin.RouterGroup, env *common.Env) {
	router.GET("/rides/:id", GetRide(env))
	router.GET("/rides", GetAllRides(env))
	router.POST("/rides", PostRides(env))
	router.DELETE("/rides/:id", DeleteRides(env))
}
