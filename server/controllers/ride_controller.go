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

func GetRides(finder models.RideFinder) gin.HandlerFunc {
	return func(c *gin.Context) {
		rawId := c.Param("id")
		objID, err := primitive.ObjectIDFromHex(rawId)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ride id"})
			return
		}

		ride, err := finder.FindOne(c, bson.M{"_id": objID})
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, ride.Jsonify())
	}
}

func GetAllRides(finder models.RideFinder) gin.HandlerFunc {
	return func(c *gin.Context) {
		var results []map[string]interface{}

		if rides, err := finder.FindMany(c, aggregations.RideWithDestination); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		} else {
			for _, ride := range rides {
				results = append(results, ride.Jsonify())
			}
		}

		if len(results) == 0 {
			c.JSON(http.StatusNotFound, gin.H{"error": "No matching ride was found"})
		}

		c.JSON(http.StatusOK, gin.H{"payload": results})
	}
}

func PostRides(inserter models.RideInserter, getValidator func() validators.Validator) gin.HandlerFunc {
	validator := getValidator()
	return func(c *gin.Context) {
		var rideDto models.RideDTO

		if err := c.Bind(&rideDto); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("Ride format was incorrect: %v", err)})
			return
		}

		destinationId, err := primitive.ObjectIDFromHex(rideDto.Destination)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("Ride format was incorrect: %v", err)})
			return
		}

		validator.SetDto(rideDto)
		if isValid, err := validator.Validate(c); !isValid || err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("Ride format was incorrect: %v", err)})
			return
		}

		id, err := inserter.InsertOne(c, schemas.RideSchema{
			Type:        rideDto.Type,
			Date:        rideDto.Date,
			Destination: destinationId,
			Direction:   rideDto.Direction,
			CreatedAt:   time.Now(),
		})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Ride couldn't get created: %v", err)})
			return
		}

		c.JSON(http.StatusCreated, gin.H{"id": id})
	}
}

func DeleteRides(deleter models.RideDeleter) gin.HandlerFunc {
	return func(c *gin.Context) {
		objID, err := primitive.ObjectIDFromHex(c.Param("id"))
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		deletedCount, err := deleter.DeleteOne(c, bson.M{"_id": objID})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		if deletedCount == 0 {
			c.JSON(http.StatusNotFound, gin.H{"error": fmt.Sprintf("Ride with ID %v does not exist", objID)})
		}

		c.JSON(http.StatusOK, gin.H{})
	}
}

func RegisterRideController(router *gin.RouterGroup, env *common.Env) {
	router.GET("/rides", GetRides(env.Collections.RideCollection))
	router.POST("/rides", PostRides(
		env.Collections.RideCollection,
		env.Validators.RideValidator,
	))
	router.DELETE("/rides/:id", DeleteRides(env.Collections.RideCollection))
}
