package controllers

import (
	"fmt"
	"github.com/aliparlakci/country-roads/repositories"
	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/mongo"
	"net/http"
	"time"

	"github.com/aliparlakci/country-roads/aggregations"
	"go.mongodb.org/mongo-driver/bson"

	"github.com/aliparlakci/country-roads/common"
	"github.com/aliparlakci/country-roads/models"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func GetRide(finder repositories.RideFinder) gin.HandlerFunc {
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
			logger.Errorf("RideFinder.FindMany raised an error: %v", err.Error())
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

func SearchRides(finder repositories.RideFinder) gin.HandlerFunc {
	return func(c *gin.Context) {
		logger := common.LoggerWithRequestId(c.Copy())

		// TODO: Queries should be validated before running against the database to prevent attacks
		var queries models.SearchRideQueries
		if err := c.BindQuery(&queries); err != nil {
			logger.WithField("queries", c.Request.URL.RawQuery).Debugf("cannot bind query parameters to models.SearchRideQueries: %v", err.Error())
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		pipeline := aggregations.BuildAggregation(aggregations.FilterRides(queries), aggregations.RideResponseAggregation)

		results := make([]map[string]interface{}, 0)
		if rides, err := finder.FindMany(c.Copy(), pipeline); err != nil {
			logger.Errorf("RideFinder.FindMany raised an error: %v", err.Error())
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

func PostRides(rideInserter repositories.RideInserter, locationFinder repositories.LocationFinder) gin.HandlerFunc {
	return func(c *gin.Context) {
		logger := common.LoggerWithRequestId(c.Copy())

		var user models.User
		if u, exists := c.Get("user"); !exists {
			c.JSON(http.StatusUnauthorized, gin.H{})
			return
		} else {
			user = u.(models.User)
		}

		var rideDto models.NewRideForm
		if err := rideDto.Bind(c); err != nil {
			logger.Debugf("cannot bind to models.NewRideForm: %v", err.Error())
			c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("Ride format was incorrect: %v", err)})
			return
		}

		if _, err := locationFinder.FindOne(c.Copy(), bson.M{"key": rideDto.Destination}); err == mongo.ErrNoDocuments {
			logger.WithField("destination", rideDto.Destination).Debug("no location with the destination key exists")
			c.JSON(http.StatusBadRequest, gin.H{"error": "destination does not exist"})
		} else if err != nil {
			logger.Errorf("locaitionFinder.FindOne raised an error while querying for destination: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "cannot validate the destination"})
			return
		}

		newRide := models.RideSchema{
			Type:        rideDto.Type,
			Date:        rideDto.Date,
			Destination: rideDto.Destination,
			Direction:   rideDto.Direction,
			CreatedAt:   time.Now(),
			Owner:       user.ID,
		}
		id, err := rideInserter.InsertOne(c.Copy(), newRide)
		if err != nil {
			logger.Errorf("RideInserter.InsertOne raised an error %v", err.Error())
			c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Ride couldn't get created: %v", err)})
			return
		}
		logger.WithFields(logrus.Fields{"id": id}).Info("ride with id is created")
		c.JSON(http.StatusCreated, gin.H{"id": id})
	}
}

func DeleteRides(repository repositories.RideRepository) gin.HandlerFunc {
	return func(c *gin.Context) {
		logger := common.LoggerWithRequestId(c.Copy())

		var user models.User
		if u, exists := c.Get("user"); !exists {
			c.JSON(http.StatusUnauthorized, gin.H{})
			return
		} else {
			user = u.(models.User)
		}

		objID, err := primitive.ObjectIDFromHex(c.Param("id"))
		if err != nil {
			logger.WithFields(logrus.Fields{"id": c.Param("id")}).Debug("id cannot get converted to primitive.ObjectID")
			c.JSON(http.StatusBadRequest, gin.H{"error": "ride id is not valid"})
			return
		}

		ride, err := repository.FindOne(c.Copy(), bson.M{"_id": objID})
		if err == mongo.ErrNoDocuments {
			logger.WithField("id", objID.Hex()).Debug("cannot find ride with id")
			c.JSON(http.StatusBadRequest, gin.H{"error": "ride does not exist"})
			return
		} else if err != nil {
			logger.WithField("id", objID.Hex()).Errorf("RideRepository.FindOne raised an error while fetching ride with id: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "cannot fetch the ride with id"})
			return
		}

		if ride.Owner != user.ID {
			logger.WithFields(logrus.Fields{
				"requester": user.ID.Hex(),
				"owner": ride.Owner.Hex(),
			}).Debug("requester cannot delete ride of owner")
			c.JSON(http.StatusUnauthorized, gin.H{"error": "you cannot delete someone else's post"})
			return
		}

		deletedCount, err := repository.DeleteOne(c.Copy(), bson.M{"_id": objID})
		if err != nil {
			logger.Errorf("RideDeleter.DeleteOne raised an error: %v", err.Error())
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
		env.Repositories.LocationRepository,
	))
	router.DELETE("/rides/:id", DeleteRides(env.Repositories.RideRepository))
}
