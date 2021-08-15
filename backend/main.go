package main

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

type Ride struct {
	ID   uint32    `json:"id"`
	Type string    `json:"type"`
	Date time.Time `json:"date"`
	From string    `json:"from"`
	To   string    `json:"to"`
}

type RideDTO struct {
	Type string `json:"type"`
	Date int64  `json:"date"`
	From string `json:"from"`
	To   string `json:"to"`
}

var rides = []Ride{
	{ID: 1, Type: "taxi", Date: time.Now(), From: "Campus", To: "Sabiha Gokcen Airport"},
	{ID: 2, Type: "offer", Date: time.Now(), From: "Campus", To: "Kadikoy"},
	{ID: 3, Type: "request", Date: time.Now(), From: "Taksim", To: "Kampus"},
}

var ridesId uint32 = 3

func getRide(ctx *gin.Context) {
	var id uint32
	if id_, err := strconv.ParseUint(ctx.Param("id"), 10, 32); err != nil {
		ctx.JSON(http.StatusBadRequest, "Bad ID")
		return
	} else {
		id = uint32(id_)
	}

	for _, ride := range rides {
		if ride.ID == id {
			ctx.JSON(http.StatusOK, ride)
			return
		}
	}

	ctx.JSON(http.StatusNotFound, fmt.Sprintf("Ride with id %v not found", id))
}

func getAllRides(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, rides)
}

func postRides(ctx *gin.Context) {
	var ride RideDTO

	if err := ctx.BindJSON(&ride); err != nil {
		ctx.JSON(http.StatusBadRequest, fmt.Sprintf("Ride format was incorrect: %v", err))
		return
	}

	newRide := Ride{
		ID:   ridesId + 1,
		Type: ride.Type,
		Date: time.Unix(ride.Date, 0),
		From: ride.From,
		To:   ride.To,
	}
	ridesId += 1

	rides = append(rides, newRide)
	ctx.JSON(http.StatusCreated, newRide)
}

func deleteRide(ctx *gin.Context) {
	var id uint32
	if id_, err := strconv.ParseUint(ctx.Param("id"), 10, 32); err != nil {
		ctx.JSON(http.StatusBadRequest, "Bad ID")
		return
	} else {
		id = uint32(id_)
	}

	index := -1
	for i, ride := range rides {
		if ride.ID == id {
			index = i
		}
	}

	if index != -1 {
		rides[index] = rides[len(rides)-1]
		rides = rides[:len(rides)-1]
		ctx.Status(http.StatusNoContent)
	} else {
		ctx.JSON(http.StatusNotFound, fmt.Sprintf("Ride with id %v not found", id))
	}
}

func main() {
	router := gin.Default()
	v1 := router.Group("v1")
	{
		v1.GET("/ride/:id", getRide)
		v1.GET("/rides", getAllRides)
		v1.POST("/rides", postRides)
		v1.DELETE("/ride/:id", deleteRide)
	}
	router.Run(":8080")
}
