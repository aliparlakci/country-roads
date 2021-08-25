package controllers

import (
	"fmt"
	"net/http"
	"time"

	"example.com/country-roads/common"
	"example.com/country-roads/models"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func getRide(env *common.Env) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var ride models.Ride
		collection := env.Db.Database("country-roads").Collection("rides")

		objID, err := primitive.ObjectIDFromHex(ctx.Param("id"))
		if err != nil {
			ctx.String(http.StatusInternalServerError, err.Error())
			return
		}

		err = collection.FindOne(ctx, bson.M{"_id": objID}).Decode(&ride)
		if err != nil {
			ctx.String(http.StatusInternalServerError, err.Error())
			return
		}

		ctx.JSON(http.StatusOK, ride)
	}
}

func getAllRides(env *common.Env) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		results := make([]models.Ride, 0)

		collection := env.Db.Database("country-roads").Collection("rides")

		cursor, err := collection.Find(ctx, bson.M{})
		if err != nil {
			ctx.String(http.StatusInternalServerError, err.Error())
			return
		}

		for cursor.Next(ctx) {
			var ride models.Ride
			if err := cursor.Decode(&ride); err != nil {
				ctx.String(http.StatusInternalServerError, err.Error())
				return
			}

			results = append(results, ride)
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

		collection := env.Db.Database("country-roads").Collection("rides")
		result, err := collection.InsertOne(ctx, models.Ride{
			Type:        ride.Type,
			Date:        ride.Date,
			Destination: ride.Destination,
			Direction:   ride.Direction,
			CreatedAt:   time.Now(),
		})

		if err != nil {
			ctx.String(http.StatusInternalServerError, fmt.Sprintf("Ride couldn't get created: %v", err))
			return
		}

		ctx.JSON(http.StatusCreated, result.InsertedID)
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

		count, err := collection.DeleteOne(ctx, bson.M{"_id": objID})
		if err != nil {
			ctx.String(http.StatusInternalServerError, err.Error())
			return
		}

		ctx.JSON(http.StatusOK, count)
	}
}

func RegisterRideController(router *gin.RouterGroup, env *common.Env) {
	router.GET("/rides/:id", getRide(env))
	router.GET("/rides", getAllRides(env))
	router.POST("/rides", postRides(env))
	router.DELETE("/rides/:id", deleteRides(env))
}