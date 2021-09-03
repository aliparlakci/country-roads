package models

import (
	"errors"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
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
