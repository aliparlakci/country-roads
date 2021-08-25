package models

import (
	"context"
	"fmt"
	"time"

	"example.com/country-roads/schemas"
	"go.mongodb.org/mongo-driver/bson"
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
	Type        string    `bson:"type" json:"type" form:"type" binding:"required,validridetype"`
	Date        time.Time `bson:"date" json:"date" form:"date" time_format:"unix" binding:"required,futuredate"`
	Direction   string    `bson:"direction" json:"direction" form:"direction" binding:"required,validdirection"`
	Destination string    `bson:"destination" json:"destination" form:"destination" binding:"required"`
}

var rideAggregationPipeline = mongo.Pipeline{
	bson.D{
		primitive.E{
			Key: "$lookup",
			Value: bson.D{
				primitive.E{Key: "from", Value: "locations"},
				primitive.E{Key: "localField", Value: "destination"},
				primitive.E{Key: "foreignField", Value: "_id"}, primitive.E{Key: "as", Value: "destination"},
			},
		},
	},
	bson.D{
		primitive.E{
			Key: "$unwind",
			Value: bson.D{
				primitive.E{Key: "path", Value: "$destination"},
				primitive.E{Key: "preserveNullAndEmptyArrays", Value: false},
			},
		},
	},
}

func GetSingleRide(ctx context.Context, database *mongo.Database, objID primitive.ObjectID) (Ride, error) {
	var ride Ride
	result := database.Collection("rides").FindOne(ctx, bson.M{"_id": objID})
	if err := result.Err(); err != nil {
		return Ride{}, err
	}
	err := result.Decode(&ride)
	return ride, err
}

func GetRides(ctx context.Context, database *mongo.Database) ([]Ride, error) {
	results := make([]Ride, 0)

	pipeline := rideAggregationPipeline

	collection := database.Collection("rides")
	cursor, err := collection.Aggregate(ctx, pipeline)
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

func CreateRide(ctx context.Context, database *mongo.Database, newRide RideDTO) (interface{}, error) {
	destination_id, err := primitive.ObjectIDFromHex(newRide.Destination)
	if err != nil {
		return nil, err
	}

	if _, err := GetSingleLocation(ctx, database, destination_id); err != nil {
		return nil, fmt.Errorf("Location with id %v does not exist", destination_id.Hex())
	}

	result, err := database.Collection("rides").InsertOne(ctx, schemas.RideSchema{
		Type:        newRide.Type,
		Date:        newRide.Date,
		Destination: destination_id,
		Direction:   newRide.Direction,
		CreatedAt:   time.Now(),
	})
	if err != nil {
		return nil, err
	}

	return result.InsertedID, nil
}

func DeleteRide(ctx context.Context, database *mongo.Database, objID primitive.ObjectID) (int64, error) {
	collection := database.Collection("rides")
	result, err := collection.DeleteOne(ctx, bson.M{"_id": objID})
	if err != nil {
		return 0, err
	}

	return result.DeletedCount, nil
}

func (r Ride) JSON() map[string]interface{} {
	return map[string]interface{}{
		"id":          r.ID.Hex(),
		"type":        r.Type,
		"date":        fmt.Sprint(r.Date.Unix()),
		"direction":   r.Direction,
		"destination": r.Destination.JSON(),
		"createdAt":   fmt.Sprint(r.CreatedAt.Unix()),
	}
}
