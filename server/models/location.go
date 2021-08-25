package models

import (
	"context"

	"example.com/country-roads/schemas"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type Location struct {
	ID      primitive.ObjectID `bson:"_id" json:"id"`
	Display string             `bson:"display" json:"display"`
	Parent  primitive.ObjectID `bson:"parent,omitempty" json:"parent"`
}

func GetSingleLocation(ctx context.Context, database *mongo.Database, objID primitive.ObjectID) (Location, error) {
	var location Location
	err := database.Collection("locations").FindOne(ctx, bson.M{"_id": objID}).Decode(&location)
	return location, err
}

func GetLocations(ctx context.Context, database *mongo.Database) ([]schemas.LocationSchema, error) {
	results := make([]schemas.LocationSchema, 0)

	collection := database.Collection("locations")
	cursor, err := collection.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}

	for cursor.Next(ctx) {
		var location schemas.LocationSchema
		if err := cursor.Decode(&location); err != nil {
			return nil, err
		}

		results = append(results, location)
	}

	return results, nil
}

func (location Location) String() string {
	return location.Display
}

func (location Location) JSON() map[string]interface{} {
	return map[string]interface{}{
		"id":      location.ID.Hex(),
		"display": location.Display,
		"parent":  location.Parent,
	}
}
