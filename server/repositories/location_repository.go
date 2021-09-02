package repositories

import (
	"context"
	"github.com/aliparlakci/country-roads/models"
	"go.mongodb.org/mongo-driver/mongo"
)

type LocationCollection struct {
	Collection *mongo.Collection
}

type LocationFinder interface {
	FindOne(ctx context.Context, filter interface{}) (models.Location, error)
	FindMany(ctx context.Context, pipeline interface{}) (models.Locations, error)
}

type LocationInserter interface {
	InsertOne(ctx context.Context, candidate models.LocationSchema) (interface{}, error)
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

func (l *LocationCollection) FindOne(ctx context.Context, filter interface{}) (models.Location, error) {
	var location models.Location
	result := l.Collection.FindOne(ctx, filter)
	if err := result.Err(); err != nil {
		return location, err
	}
	err := result.Decode(&location)
	return location, err
}

func (l *LocationCollection) FindMany(ctx context.Context, pipeline interface{}) (models.Locations, error) {
	results := make([]models.Location, 0)

	cursor, err := l.Collection.Aggregate(ctx, pipeline)
	if err != nil {
		return nil, err
	}

	for cursor.Next(ctx) {
		var location models.Location
		if err := cursor.Decode(&location); err != nil {
			return nil, err
		}

		results = append(results, location)
	}

	return results, nil
}

func (l *LocationCollection) InsertOne(ctx context.Context, candidate models.LocationSchema) (interface{}, error) {
	result, err := l.Collection.InsertOne(ctx, candidate)
	if err != nil {
		return nil, err
	}

	return result.InsertedID, nil
}

