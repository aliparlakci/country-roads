package controllers

import (
	"fmt"
	"net/http"

	"example.com/country-roads/common"
	"example.com/country-roads/models"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func getRide(env *common.Env) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		database := env.Db.Database("country-roads")

		objID, err := primitive.ObjectIDFromHex(ctx.Param("id"))
		if err != nil {
			ctx.String(http.StatusInternalServerError, err.Error())
			return
		}

		ride, err := models.GetSingleRide(ctx, database, objID)
		if err != nil {
			ctx.String(http.StatusInternalServerError, err.Error())
			return
		}

		ctx.JSON(http.StatusOK, ride.JSON())
	}
}

func getAllRides(env *common.Env) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var results []map[string]interface{}

		if rides, err := models.GetRides(ctx, env.Db.Database("country-roads")); err != nil {
			ctx.String(http.StatusInternalServerError, err.Error())
			return
		} else {
			for _, ride := range rides {
				results = append(results, ride.JSON())
			}
		}

		ctx.JSON(http.StatusOK, results)
	}
}

func postRides(env *common.Env) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var ride models.RideDTO

		if err := ctx.Bind(&ride); err != nil {
			ctx.JSON(http.StatusBadRequest, fmt.Sprintf("Ride format was incorrect: %v", err))
			return
		}

		id, err := models.CreateRide(ctx, env.Db.Database("country-roads"), ride)
		if err != nil {
			ctx.String(http.StatusInternalServerError, fmt.Sprintf("Ride couldn't get created: %v", err))
			return
		}

		ctx.JSON(http.StatusCreated, id)
	}
}

func deleteRides(env *common.Env) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		collection := env.Db.Database("country-roads").Collection("rides")

		objID, err := primitive.ObjectIDFromHex(ctx.Param("id"))
		if err != nil {
			ctx.String(http.StatusInternalServerError, err.Error())
			return
		}

		result, err := collection.DeleteOne(ctx, bson.M{"_id": objID})
		if err != nil {
			ctx.String(http.StatusInternalServerError, err.Error())
			return
		}

		if result.DeletedCount == 0 {
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
