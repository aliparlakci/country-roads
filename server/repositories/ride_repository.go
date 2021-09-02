package repositories

import (
	"context"
	"github.com/aliparlakci/country-roads/models"
	"go.mongodb.org/mongo-driver/mongo"
)

type RideCollection struct {
	Collection *mongo.Collection
}

type RideFinder interface {
	FindOne(ctx context.Context, filter interface{}) (models.Ride, error)
	FindMany(ctx context.Context, pipeline interface{}) (models.Rides, error)
}

type RideInserter interface {
	InsertOne(ctx context.Context, candidate models.RideSchema) (interface{}, error)
}

type RideDeleter interface {
	DeleteOne(ctx context.Context, filter interface{}) (int64, error)
}

type RideRepository interface {
	RideFinder
	RideInserter
	RideDeleter
}

func (r *RideCollection) FindOne(ctx context.Context, filter interface{}) (models.Ride, error) {
	var ride models.Ride
	result := r.Collection.FindOne(ctx, filter)
	if err := result.Err(); err != nil {
		return ride, err
	}
	err := result.Decode(&ride)
	return ride, err
}

func (r *RideCollection) FindMany(ctx context.Context, pipeline interface{}) (models.Rides, error) {
	results := make([]models.Ride, 0)

	cursor, err := r.Collection.Aggregate(ctx, pipeline)
	if err != nil {
		return nil, err
	}

	for cursor.Next(ctx) {
		var ride models.Ride
		if err := cursor.Decode(&ride); err != nil {
			return nil, err
		}

		results = append(results, ride)
	}

	return results, nil
}

func (r *RideCollection) InsertOne(ctx context.Context, candidate models.RideSchema) (interface{}, error) {
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
