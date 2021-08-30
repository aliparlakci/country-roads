package models

//go:generate mockgen -destination=../mocks/mock_ride_model.go -package=mocks example.com/country-roads/models RideRepository,RideFinder,RideInserter,RideDeleter

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type Ride struct {
	ID          primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	Type        string             `bson:"type" json:"type"`
	Date        time.Time          `bson:"date" json:"date" time_format:"unix"`
	Direction   string             `bson:"direction" json:"direction"`
	Destination Location           `bson:"destination" json:"destination"`
	CreatedAt   time.Time          `bson:"createdAt" json:"createdAt" time_format:"unix"`
	// Author      User               `bson:"author" json:"author"`
}

type Rides []Ride

type NewRideForm struct {
	Type         string    `bson:"type" json:"type" form:"type" binding:"required"`
	Date         time.Time `bson:"date" json:"date" form:"date" time_format:"unix" binding:"required"`
	Direction    string    `bson:"direction" json:"direction" form:"direction" binding:"required"`
	Destination  string    `bson:"destination" json:"destination" form:"destination" binding:"required"`
}

type RideSchema struct {
	ID          primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	Type        string             `bson:"type" json:"type"`
	Date        time.Time          `bson:"date" json:"date" time_format:"unix"`
	Direction   string             `bson:"direction" json:"direction"`
	Destination string             `bson:"destination" json:"destination"`
	CreatedAt   time.Time          `bson:"createdAt" json:"createdAt,omitempty" time_format:"unix"`
}

type SearchRideQueries struct {
	Type        string    `form:"type"`
	StartDate   time.Time `form:"start_date" time_format:"unix"`
	EndDate     time.Time `form:"end_date" time_format:"unix"`
	Direction   string    `form:"direction"`
	Destination string    `form:"destination"`
}

type RideCollection struct {
	Collection *mongo.Collection
}

type RideFinder interface {
	FindOne(ctx context.Context, filter interface{}) (Ride, error)
	FindMany(ctx context.Context, pipeline interface{}) (Rides, error)
}

type RideInserter interface {
	InsertOne(ctx context.Context, candidate RideSchema) (interface{}, error)
}

type RideDeleter interface {
	DeleteOne(ctx context.Context, filter interface{}) (int64, error)
}

type RideRepository interface {
	RideFinder
	RideInserter
	RideDeleter
}

func (r *RideCollection) FindOne(ctx context.Context, filter interface{}) (Ride, error) {
	var ride Ride
	result := r.Collection.FindOne(ctx, filter)
	if err := result.Err(); err != nil {
		return ride, err
	}
	err := result.Decode(&ride)
	return ride, err
}

func (r *RideCollection) FindMany(ctx context.Context, pipeline interface{}) (Rides, error) {
	results := make([]Ride, 0)

	cursor, err := r.Collection.Aggregate(ctx, pipeline)
	if err != nil {
		return nil, err
	}

	for cursor.Next(ctx) {
		var ride Ride
		if err := cursor.Decode(&ride); err != nil {
			return nil, err
		}

		results = append(results, ride)
	}

	return results, nil
}

func (r *RideCollection) InsertOne(ctx context.Context, candidate RideSchema) (interface{}, error) {
	result, err := r.Collection.InsertOne(ctx, candidate)
	if err != nil {
		return nil, err
	}

	return result.InsertedID, nil
}

func (r *RideCollection) DeleteOne(ctx context.Context, filter interface{}) (int64, error) {
	result, err := r.Collection.DeleteOne(ctx, filter)
	if err != nil {
		return 0, err
	}

	return result.DeletedCount, nil
}

func (r Ride) Jsonify() map[string]interface{} {
	return map[string]interface{}{
		"id":          r.ID.Hex(),
		"type":        r.Type,
		"date":        r.Date.Unix(),
		"direction":   r.Direction,
		"destination": r.Destination.Jsonify(),
		"createdAt":   r.CreatedAt.Unix(),
	}
}

func (r Rides) Jsonify() []map[string]interface{} {
	result := make([]map[string]interface{}, 0)
	for _, location := range r {
		result = append(result, location.Jsonify())
	}
	return result
}
