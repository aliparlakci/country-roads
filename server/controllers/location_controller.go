package controllers

import (
	"fmt"
	"net/http"

	"go.mongodb.org/mongo-driver/bson"

	"github.com/aliparlakci/country-roads/common"
	"github.com/aliparlakci/country-roads/models"
	"github.com/gin-gonic/gin"
)

func GetAllLocations(finder models.LocationFinder) gin.HandlerFunc {
	return func(c *gin.Context) {
		logger := common.LoggerWithRequestId(c)

		results, err := finder.FindMany(c.Copy(), bson.D{})
		if err != nil {
			logger.Errorf("models.LocationFinder.FindMany raised an error while fetching all locations: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"results": results})
	}
}

func PostLocation(repository models.LocationRepository) gin.HandlerFunc {
	return func(c *gin.Context) {
		logger := common.LoggerWithRequestId(c)

		var locationDto models.NewLocationForm
		if err := locationDto.Bind(c); err != nil {
			logger.Debugf("cannot bind request to locationDto: %v", err)
			c.JSON(http.StatusBadRequest, gin.H{"error": "Location format was incorrect"})
			return
		}

		// does the key already exist?
		if exists, err := repository.Exists(c.Copy(), bson.M{"key": locationDto.Key}); err != nil {
			logger.WithField("location_key", locationDto.Key).Errorf("models.LocationRepository.Exists raised an error when querying for location_key: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "cannot validate the key"})
			return
		} else if exists {
			logger.WithField("location_key", locationDto.Key).Debug("location_key already exists in the collection")
			c.JSON(http.StatusBadRequest, gin.H{"error": "key already exists"})
			return
		}

		var schema models.LocationSchema
		if locationDto.ParentKey != "" {
			if exists, err := repository.Exists(c.Copy(), bson.M{"key": locationDto.ParentKey}); err != nil {
				logger.WithField("location_key", locationDto.ParentKey).Errorf("models.LocationRepository.Exists raised an error when querying for parent_key: %v", err)
				c.JSON(http.StatusInternalServerError, gin.H{"error": "cannot validate the key"})
				return
			} else if !exists {
				logger.WithField("parent_key", locationDto.ParentKey).Debug("parent_key does not exist in the collection")
				c.JSON(http.StatusBadRequest, gin.H{"error": "parentKey does not exist"})
				return
			}
			schema = models.LocationSchema{Key: locationDto.Key, Display: locationDto.Display, ParentKey: locationDto.ParentKey}
		} else {
			schema = models.LocationSchema{Key: locationDto.Key, Display: locationDto.Display}
		}

		id, err := repository.InsertOne(c.Copy(), schema)
		if err != nil {
			logger.Errorf("models.LocationRepository.InsertOne raised an error while inserting a new location: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Location couldn't get created: %v", err)})
			return
		}

		logger.WithField("id", id).WithField("key", schema.Key).Info("new location with key and id is inserted")
		c.JSON(http.StatusCreated, gin.H{"id": id})
	}
}

func RegisterLocationController(router *gin.RouterGroup, env *common.Env) {
	router.GET("/locations", GetAllLocations(env.Repositories.LocationRepository))
	router.POST("/locations", PostLocation(env.Repositories.LocationRepository))
}
