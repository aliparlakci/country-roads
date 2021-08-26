package controllers

import (
	"example.com/country-roads/aggregations"
	"example.com/country-roads/schemas"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"net/http"
	"time"

	"example.com/country-roads/common"
	"example.com/country-roads/models"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func getRide(env *common.Env) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		objID, err := primitive.ObjectIDFromHex(ctx.Param("id"))
		if err != nil {
			ctx.String(http.StatusInternalServerError, err.Error())
			return
		}

		var finder models.RideFinder = env.Collections.RideCollection

		ride, err := finder.FindOne(ctx, bson.M{"_id": objID})
		if err != nil {
			ctx.String(http.StatusInternalServerError, err.Error())
			return
		}

		ctx.JSON(http.StatusOK, ride.Jsonify())
	}
}

func getAllRides(env *common.Env) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var results []map[string]interface{}

		var finder models.RideFinder = env.Collections.RideCollection
		if rides, err := finder.FindMany(ctx, aggregations.RideWithDestination); err != nil {
			ctx.String(http.StatusInternalServerError, err.Error())
			return
		} else {
			for _, ride := range rides {
				results = append(results, ride.Jsonify())
			}
		}

		ctx.JSON(http.StatusOK, results)
	}
}

func postRides(env *common.Env) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var rideDto models.RideDTO

		if err := ctx.Bind(&rideDto); err != nil {
			ctx.JSON(http.StatusBadRequest, fmt.Sprintf("Ride format was incorrect: %v", err))
			return
		}

		destinationId, err := primitive.ObjectIDFromHex(rideDto.Destination)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, fmt.Sprintf("Ride format was incorrect: %v", err))
			return
		}

		validator := env.Validators.RideValidator()
		validator.SetDto(rideDto)
		if isValid, err := validator.Validate(ctx); !isValid || err != nil {
			ctx.JSON(http.StatusBadRequest, fmt.Sprintf("Ride format was invalid: %v", err))
		}

		var inserter models.RideInserter = env.Collections.RideCollection
		id, err := inserter.InsertOne(ctx, schemas.RideSchema{
			Type:        rideDto.Type,
			Date:        rideDto.Date,
			Destination: destinationId,
			Direction:   rideDto.Direction,
			CreatedAt:   time.Now(),
		})
		if err != nil {
			ctx.String(http.StatusInternalServerError, fmt.Sprintf("Ride couldn't get created: %v", err))
			return
		}

		ctx.JSON(http.StatusCreated, id)
	}
}

func deleteRides(env *common.Env) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		objID, err := primitive.ObjectIDFromHex(ctx.Param("id"))
		if err != nil {
			ctx.String(http.StatusInternalServerError, err.Error())
			return
		}

		var deleter models.RideDeleter = env.Collections.RideCollection
		deletedCount, err := deleter.DeleteOne(ctx, bson.M{"_id": objID})
		if err != nil {
			ctx.String(http.StatusInternalServerError, err.Error())
			return
		}

		if deletedCount == 0 {
			ctx.JSON(http.StatusNotFound, fmt.Sprintf("Ride with ID %v does not exist", objID))
		}

		ctx.JSON(http.StatusOK, "")
	}
}

func RegisterRideController(router *gin.RouterGroup, env *common.Env) {
	router.GET("/rides/:id", getRide(env))
	router.GET("/rides", getAllRides(env))
	router.POST("/rides", postRides(env))
	router.DELETE("/rides/:id", deleteRides(env))
}
