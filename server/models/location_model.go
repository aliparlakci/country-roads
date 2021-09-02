package models

//go:generate mockgen -destination=../mocks/mock_location_model.go -package=mocks github.com/aliparlakci/country-roads/models LocationFinder,LocationInserter,LocationRepository

import (
	"context"
	"errors"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type Location struct {
	ID        primitive.ObjectID `bson:"_id" json:"id"`
	Key       string             `bson:"key" json:"key"`
	Display   string             `bson:"display" json:"display"`
	ParentKey string             `bson:"parentKey,omitempty" json:"parentKey,omitempty"`
	Parent    *Location          `bson:"parent,omitempty" json:"parent,omitempty"`
}

type Locations []Location

type NewLocationForm struct {
	Key       string `bson:"key" json:"key" form:"key" binding:"required"`
	Display   string `bson:"display" json:"display" form:"display" binding:"required"`
	ParentKey string `bson:"parentKey,omitempty" json:"parentKey,omitempty" form:"parentKey,omitempty"`
}

type LocationSchema struct {
	ID        primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Key       string             `bson:"key" json:"key" form:"key"`
	Display   string             `bson:"display" json:"display"`
	ParentKey string             `bson:"parentKey,omitempty" json:"parentKey"`
}

func (n NewLocationForm) validateDisplay() bool {
	return n.Display != ""
}

func (n NewLocationForm) validate() (bool, error) {
	if display := n.validateDisplay(); !display {
		return display, nil
	}
	return true, nil
}

func (n *NewLocationForm) Bind(c *gin.Context) error {
	if err := c.Bind(n); err != nil {
		return err
	}

	if result, err := n.validate(); err != nil {
		return err
	} else if !result {
		return errors.New("")
	}
	return nil
}


type LocationCollection struct {
	Collection *mongo.Collection
}

type LocationFinder interface {
	FindOne(ctx context.Context, filter interface{}) (Location, error)
	FindMany(ctx context.Context, pipeline interface{}) (Locations, error)
}

type LocationInserter interface {
	InsertOne(ctx context.Context, candidate LocationSchema) (interface{}, error)
}

type LocationRepository interface {
	LocationFinder
	LocationInserter
	Exists(c context.Context, filter interface{}) (bool, error)
}

func (l *LocationCollection) Exists(c context.Context, filter interface{}) (bool, error) {
	if _, err := l.FindOne(c, filter); err == mongo.ErrNoDocuments {
		return false, nil
	} else if err != nil {
		return false, err
	}
	return true, nil
}

func (l *LocationCollection) FindOne(ctx context.Context, filter interface{}) (Location, error) {
	var location Location
	result := l.Collection.FindOne(ctx, filter)
	if err := result.Err(); err != nil {
		return location, err
	}
	err := result.Decode(&location)
	return location, err
}

func (l *LocationCollection) FindMany(ctx context.Context, pipeline interface{}) (Locations, error) {
	results := make([]Location, 0)

	cursor, err := l.Collection.Aggregate(ctx, pipeline)
	if err != nil {
		return nil, err
	}

	for cursor.Next(ctx) {
		var location Location
		if err := cursor.Decode(&location); err != nil {
			return nil, err
		}

		results = append(results, location)
	}

	return results, nil
}

func (l *LocationCollection) InsertOne(ctx context.Context, candidate LocationSchema) (interface{}, error) {
	result, err := l.Collection.InsertOne(ctx, candidate)
	if err != nil {
		return nil, err
	}

	return result.InsertedID, nil
}

func (l Location) String() string {
	return l.Display
}

func (l Location) Jsonify() map[string]interface{} {
	if l.Parent != nil {
		return map[string]interface{}{
			"id":        l.ID.Hex(),
			"key":       l.Key,
			"display":   l.Display,
			"parentKey": l.ParentKey,
			"parent":    l.Parent.Jsonify(),
		}
	} else {
		return map[string]interface{}{
			"id":        l.ID.Hex(),
			"key":       l.Key,
			"display":   l.Display,
			"parentKey": l.ParentKey,
		}
	}
}

func (l Locations) Jsonify() []map[string]interface{} {
	result := make([]map[string]interface{}, 0)

	for _, location := range l {
		result = append(result, location.Jsonify())
	}

	return result
}
