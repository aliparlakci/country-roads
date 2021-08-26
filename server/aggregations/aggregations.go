package aggregations

import (
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

var RideWithDestination = mongo.Pipeline{
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