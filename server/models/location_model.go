package models

import (
	"context"
	"example.com/country-roads/schemas"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type Location struct {
	ID       primitive.ObjectID `bson:"_id" json:"id"`
	Display  string             `bson:"display" json:"display"`
	ParentID primitive.ObjectID `bson:"parentId,omitempty" json:"parentId,omitempty"`
	Parent   *Location          `bson:"parent,omitempty" json:"parent,omitempty"`
}

type LocationDTO struct {
	Display  string `bson:"display" json:"display" form:"display"`
	ParentID string `bson:"parentId,omitempty" json:"parentId,omitempty" form:"parentId,omitempty"`
}

type LocationCollection struct {
	Collection *mongo.Collection
}

type LocationFinder interface {
	FindOne(ctx context.Context, filter interface{}) (Location, error)
	FindMany(ctx context.Context, pipeline interface{}) ([]Location, error)
}

type LocationInserter interface {
	InsertOne(ctx context.Context, candidate schemas.LocationSchema) (interface{}, error)
}

type LocationRepository interface {
	LocationFinder
	LocationInserter
}

func (l *LocationCollection) FindOne(ctx context.Context, filter interface{}) (Location, error) {
	var location Location
	err := l.Collection.FindOne(ctx, filter).Decode(&location)
	return location, err
}

func (l *LocationCollection) FindMany(ctx context.Context, pipeline interface{}) ([]Location, error) {
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

func (l *LocationCollection) InsertOne(ctx context.Context, candidate schemas.LocationSchema) (interface{}, error) {
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
	var parent map[string]interface{}
	if l.Parent != nil {
		parent = l.Parent.Jsonify()
	} else {
		parent = nil
	}

	return map[string]interface{}{
		"id":       l.ID.Hex(),
		"display":  l.Display,
		"parentId": l.Parent,
		"parent":   parent,
	}
}
