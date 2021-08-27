package models

//go:generate mockgen -destination=../mocks/mock_ride_model.go -package=mocks example.com/country-roads/models RideRepository

import (
	"context"
	"fmt"
	"time"

	"example.com/country-roads/schemas"
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

type RideDTO struct {
	Type        string    `bson:"type" json:"type" form:"type" binding:"required"`
	Date        time.Time `bson:"date" json:"date" form:"date" time_format:"unix" binding:"required"`
	Direction   string    `bson:"direction" json:"direction" form:"direction" binding:"required"`
	Destination string    `bson:"destination" json:"destination" form:"destination" binding:"required"`
}

type RideCollection struct {
	Collection *mongo.Collection
}

type RideFinder interface {
	FindOne(ctx context.Context, filter interface{}) (Ride, error)
	FindMany(ctx context.Context, pipeline interface{}) ([]Ride, error)
}

type RideInserter interface {
	InsertOne(ctx context.Context, candidate schemas.RideSchema) (interface{}, error)
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

func (r *RideCollection) FindMany(ctx context.Context, pipeline interface{}) ([]Ride, error) {
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

func (r *RideCollection) InsertOne(ctx context.Context, candidate schemas.RideSchema) (interface{}, error) {
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
		"date":        fmt.Sprint(r.Date.Unix()),
		"direction":   r.Direction,
		"destination": r.Destination.Jsonify(),
		"createdAt":   fmt.Sprint(r.CreatedAt.Unix()),
	}
}
