package controllers

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"example.com/country-roads/common"
	"example.com/country-roads/models"
	"github.com/gin-gonic/gin"
)

func getRides(env *common.Env) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var id uint32
		if id_, err := strconv.ParseUint(ctx.Param("id"), 10, 32); err != nil {
			ctx.JSON(http.StatusBadRequest, "Bad ID")
			return
		} else {
			id = uint32(id_)
		}

		for _, ride := range models.Rides {
			if ride.ID == id {
				ctx.JSON(http.StatusOK, ride)
				return
			}
		}

		ctx.JSON(http.StatusNotFound, fmt.Sprintf("Ride with id %v not found", id))
	}
}

func getAllRides(env *common.Env) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, models.Rides)
	}
}

func putRides(env *common.Env) gin.HandlerFunc {
	return func(ctx *gin.Context) {}
}

func postRides(env *common.Env) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var ride models.RideDTO

		if err := ctx.BindJSON(&ride); err != nil {
			ctx.JSON(http.StatusBadRequest, fmt.Sprintf("Ride format was incorrect: %v", err))
			return
		}

		newRide := models.Ride{
			ID:   models.RidesId + 1,
			Type: ride.Type,
			Date: time.Unix(ride.Date, 0),
			From: ride.From,
			To:   ride.To,
		}
		models.RidesId += 1

		models.Rides = append(models.Rides, newRide)
		ctx.JSON(http.StatusCreated, newRide)
	}
}

func deleteRides(env *common.Env) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var id uint32
		if id_, err := strconv.ParseUint(ctx.Param("id"), 10, 32); err != nil {
			ctx.JSON(http.StatusBadRequest, "Bad ID")
			return
		} else {
			id = uint32(id_)
		}

		index := -1
		for i, ride := range models.Rides {
			if ride.ID == id {
				index = i
			}
		}

		if index != -1 {
			models.Rides[index] = models.Rides[len(models.Rides)-1]
			models.Rides = models.Rides[:len(models.Rides)-1]
			ctx.Status(http.StatusNoContent)
		} else {
			ctx.JSON(http.StatusNotFound, fmt.Sprintf("Ride with id %v not found", id))
		}
	}
}

func RegisterRideController(router *gin.RouterGroup, env *common.Env) {
	router.GET("/rides/:id", getRides(env))
	router.GET("/rides", getAllRides(env))
	router.POST("/rides", postRides(env))
	router.PUT("/rides/:id", putRides(env))
	router.DELETE("/rides/:id", deleteRides(env))
}
