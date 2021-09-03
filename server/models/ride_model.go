package models

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type RideResponse struct {
	ID          primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	Type        string             `bson:"type" json:"type"`
	Date        time.Time          `bson:"date" json:"date" time_format:"unix"`
	Direction   string             `bson:"direction" json:"direction"`
	Destination Location           `bson:"destination" json:"destination"`
	CreatedAt   time.Time          `bson:"createdAt" json:"createdAt" time_format:"unix"`
	Owner       UserResponse       `bson:"owner" json:"owner"`
}

type Rides []RideResponse

type NewRideForm struct {
	Type        string    `bson:"type" json:"type" form:"type" binding:"required"`
	Date        time.Time `bson:"date" json:"date" form:"date" time_format:"unix" binding:"required"`
	Direction   string    `bson:"direction" json:"direction" form:"direction" binding:"required"`
	Destination string    `bson:"destination" json:"destination" form:"destination" binding:"required"`
}

type RideSchema struct {
	ID          primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	Type        string             `bson:"type" json:"type"`
	Date        time.Time          `bson:"date" json:"date" time_format:"unix"`
	Direction   string             `bson:"direction" json:"direction"`
	Destination string             `bson:"destination" json:"destination"`
	CreatedAt   time.Time          `bson:"createdAt" json:"createdAt,omitempty" time_format:"unix"`
	Owner       primitive.ObjectID `bson:"owner" json:"owner"`
}

type SearchRideQueries struct {
	Type        string    `form:"type" json:"type"`
	StartDate   time.Time `form:"start_date" json:"start_date" time_format:"unix"`
	EndDate     time.Time `form:"end_date" json:"end_date" time_format:"unix"`
	Direction   string    `form:"direction" json:"direction"`
	Destination string    `form:"destination" json:"destination"`
}

func (n NewRideForm) validateDate() bool {
	date := n.Date
	today := time.Now()
	if today.Year() < date.Year() {
		return true
	} else if today.Year() == date.Year() {
		if today.Month() < date.Month() {
			return true
		} else if today.Month() == date.Month() {
			if today.Day() <= date.Day() {
				return true
			}
		}
	}
	return false
}

func (n NewRideForm) validateType() bool {
	switch n.Type {
	case "offer":
		return true
	case "request":
		return true
	case "taxi":
		return true
	default:
		return false
	}
}

func (n NewRideForm) validateDirection() bool {
	switch n.Direction {
	case "to_campus":
		return true
	case "from_campus":
		return true
	default:
		return false
	}
}

func (n NewRideForm) validate() (bool, error) {
	if !n.validateDate() {
		return false, fmt.Errorf("ride date is not valid")
	}
	if !n.validateDirection() {
		return false, fmt.Errorf("ride direction is not valid")
	}
	if !n.validateType() {
		return false, fmt.Errorf("ride type is not valid")
	}

	return true, nil
}

func (n *NewRideForm) Bind(c *gin.Context) error {
	if err := c.Bind(n); err != nil {
		return fmt.Errorf(err.Error())
	}

	if result, err := n.validate(); err != nil {
		return err
	} else if !result {
		return errors.New("")
	}
	return nil
}

func (r RideResponse) Jsonify() map[string]interface{} {
	return map[string]interface{}{
		"id":          r.ID.Hex(),
		"type":        r.Type,
		"date":        r.Date.Unix(),
		"direction":   r.Direction,
		"destination": r.Destination.Jsonify(),
		"createdAt":   r.CreatedAt.Unix(),
		"owner":       r.Owner.Jsonify(),
	}
}

func (r Rides) Jsonify() []map[string]interface{} {
	result := make([]map[string]interface{}, 0)
	for _, location := range r {
		result = append(result, location.Jsonify())
	}
	return result
}
