package models

import (
	"context"
	"fmt"

	"example.com/country-roads/schemas"
	"go.mongodb.org/mongo-driver/bson"
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
	Display  string `bson:"display" json:"display"`
	ParentID string `bson:"parentId,omitempty" json:"parentId,omitempty"`
}

func GetSingleLocation(ctx context.Context, database *mongo.Database, objID primitive.ObjectID) (Location, error) {
	var location Location
	err := database.Collection("locations").FindOne(ctx, bson.M{"_id": objID}).Decode(&location)
	return location, err
}

func GetLocations(ctx context.Context, database *mongo.Database) ([]Location, error) {
	results := make([]Location, 0)

	collection := database.Collection("locations")
	cursor, err := collection.Find(ctx, bson.M{})
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

func RegisterLocation(ctx context.Context, database *mongo.Database, location LocationDTO) (interface{}, error) {
	var newLocation schemas.LocationSchema
	newLocation.Display = location.Display

	if location.ParentID != "" {
		parentId, err := primitive.ObjectIDFromHex(location.ParentID)
		if err != nil {
			return nil, err
		}

		newLocation.ParentID = parentId

		if _, err := GetSingleLocation(ctx, database, parentId); err != nil {
			return nil, fmt.Errorf("Location with id %v does not exist", parentId.Hex())
		}
	}

	result, err := database.Collection("locations").InsertOne(ctx, newLocation)
	if err != nil {
		return nil, err
	}

	return result.InsertedID, nil
}

func (l Location) String() string {
	return l.Display
}

func (lDto LocationDTO) Validate() (bool, error) {
	return true, nil
}

func (location Location) JSONify() map[string]interface{} {
	var parent map[string]interface{}
	if location.Parent != nil {
		parent = location.Parent.JSONify()
	} else {
		parent = nil
	}

	return map[string]interface{}{
		"id":       location.ID.Hex(),
		"display":  location.Display,
		"parentId": location.Parent,
		"parent":   parent,
	}
}
