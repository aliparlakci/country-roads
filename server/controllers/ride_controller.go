package controllers

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"net/http"
	"time"

	"github.com/aliparlakci/country-roads/aggregations"
	"github.com/aliparlakci/country-roads/validators"
	"go.mongodb.org/mongo-driver/bson"

	"github.com/aliparlakci/country-roads/common"
	"github.com/aliparlakci/country-roads/models"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func GetRide(finder models.RideFinder) gin.HandlerFunc {
	return func(c *gin.Context) {
		logger := common.LoggerWithRequestId(c.Copy())

		rawId := c.Param("id")
		objID, err := primitive.ObjectIDFromHex(rawId)
		if err != nil {
			logger.WithFields(logrus.Fields{"id": rawId}).Debug("id cannot get converted to primitive.ObjectID")
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ride id"})
			return
		}

		filter := bson.D{primitive.E{Key: "$match", Value: bson.M{"_id": objID}}}
		pipeline := aggregations.BuildAggregation([]bson.D{filter})

		rides, err := finder.FindMany(c.Copy(), pipeline)
		if err != nil {
			logger.WithField("error", err.Error()).Error("models.RideFinder.FindMany raised an error")
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		if len(rides) < 1 {
			logger.WithField("id", objID.Hex()).Info("ride with id does not exist")
			c.JSON(http.StatusNotFound, gin.H{"results": gin.H{}})
			return
		}

		c.JSON(http.StatusOK, gin.H{"results": rides[0].Jsonify()})
	}
}

func SearchRides(finder models.RideFinder) gin.HandlerFunc {
	return func(c *gin.Context) {
		logger := common.LoggerWithRequestId(c.Copy())

		// TODO: Queries should be validated before running against the database to prevent attacks
		var queries models.SearchRideQueries
		if err := c.BindQuery(&queries); err != nil {
			logger.WithFields(logrus.Fields{
				"queries": c.Request.URL.RawQuery,
				"error": err.Error(),
			}).Debug("cannot bind query parameters to models.SearchRideQueries")
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		pipeline := aggregations.BuildAggregation(aggregations.FilterRides(queries), aggregations.RideWithDestination)

		results := make([]map[string]interface{}, 0)
		if rides, err := finder.FindMany(c.Copy(), pipeline); err != nil {
			logger.WithFields(logrus.Fields{
				"error": err.Error(),
			}).Error("models.RideFinder.FindMany raised an error")
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		} else {
			for _, ride := range rides {
				results = append(results, ride.Jsonify())
			}
		}

		c.JSON(http.StatusOK, gin.H{"results": results})
	}
}

func PostRides(inserter models.RideInserter, validators validators.IValidatorFactory) gin.HandlerFunc {
	validator, err := validators.GetValidator("rides")
	if err != nil {
		panic(err)
	}
	return func(c *gin.Context) {
		logger := common.LoggerWithRequestId(c.Copy())

		var rideDto models.NewRideForm
		if err := c.Bind(&rideDto); err != nil {
			logger.WithField("error", err.Error()).Debug("cannot bind to models.NewRideForm")
			c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("Ride format was incorrect: %v", err)})
			return
		}

		validator.SetDto(rideDto)
		if isValid, err := validator.Validate(c.Copy()); !isValid || err != nil {
			logger.WithFields(logrus.Fields{
				"data": common.JsonMarshalNoError(rideDto),
				"error": err.Error(),
			}).Debug("rideDto (models.NewRideForm) is not a valid ride")
			c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("Ride format was incorrect: %v", err)})
			return
		}

		newRide := models.RideSchema{
			Type:        rideDto.Type,
			Date:        rideDto.Date,
			Destination: rideDto.Destination,
			Direction:   rideDto.Direction,
			CreatedAt:   time.Now(),
		}
		id, err := inserter.InsertOne(c.Copy(), newRide)
		if err != nil {
			logger.WithFields(logrus.Fields{"error": err.Error()}).Error("models.RideInserter.InsertOne raised an error")
			c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Ride couldn't get created: %v", err)})
			return
		}
		logger.WithFields(logrus.Fields{"id": id}).Info("ride with id is created")
		c.JSON(http.StatusCreated, gin.H{"id": id})
	}
}

func DeleteRides(deleter models.RideDeleter) gin.HandlerFunc {
	return func(c *gin.Context) {
		logger := common.LoggerWithRequestId(c.Copy())

		objID, err := primitive.ObjectIDFromHex(c.Param("id"))
		if err != nil {
			logger.WithFields(logrus.Fields{"id": c.Param("id")}).Debug("id cannot get converted to primitive.ObjectID")
			c.JSON(http.StatusBadRequest, gin.H{"error": "ride id is not valid"})
			return
		}

		deletedCount, err := deleter.DeleteOne(c.Copy(), bson.M{"_id": objID})
		if err != nil {
			logger.WithFields(logrus.Fields{"error": err.Error()}).Error("models.RideDeleter.DeleteOne raised an error")
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		if deletedCount == 0 {
			logger.WithFields(logrus.Fields{"id": objID.Hex()}).Debug("ride with id does not exist, thus cannot get deleted")
			c.JSON(http.StatusNotFound, gin.H{"error": fmt.Sprintf("Ride with ID %v does not exist", objID)})
			return
		}

		logger.WithFields(logrus.Fields{"id": objID.Hex()}).Info("ride with id is deleted")
		c.JSON(http.StatusOK, gin.H{})
	}
}

func RegisterRideController(router *gin.RouterGroup, env *common.Env) {
	router.GET("/rides/:id", GetRide(env.Repositories.RideRepository))
	router.GET("/rides", SearchRides(env.Repositories.RideRepository))
	router.POST("/rides", PostRides(
		env.Repositories.RideRepository,
		env.ValidatorFactory,
	))
	router.DELETE("/rides/:id", DeleteRides(env.Repositories.RideRepository))
}
