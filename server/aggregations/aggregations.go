package aggregations

import (
	"github.com/aliparlakci/country-roads/common"
	"github.com/aliparlakci/country-roads/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func FilterRides(queries models.SearchRideQueries) []bson.D {
	keys := bson.D{}
	if queries.Type != "" {
		keys = append(keys, primitive.E{Key: "type", Value: queries.Type})
	}
	if queries.Direction != "" {
		keys = append(keys, primitive.E{Key: "direction", Value: queries.Direction})
	}
	if queries.Destination != "" {
		keys = append(keys, primitive.E{Key: "destination", Value: queries.Destination})
	}
	if queries.StartDate.Unix() != common.MinDate || queries.EndDate.Unix() != common.MinDate {
		dateRange := bson.D{}
		if queries.StartDate.Unix() != common.MinDate {
			dateRange = append(dateRange, primitive.E{Key: "$gte", Value: queries.StartDate})
		}
		if queries.EndDate.Unix() != common.MinDate {
			dateRange = append(dateRange, primitive.E{Key: "$lte", Value: queries.EndDate})
		}
		keys = append(keys, primitive.E{Key: "date", Value: dateRange})
	}

	return []bson.D{
		bson.D{
			primitive.E{Key: "$match", Value: keys},
		},
	}
}

var RideWithDestination = []bson.D{
	bson.D{
		primitive.E{
			Key: "$lookup",
			Value: bson.D{
				primitive.E{Key: "from", Value: "locations"},
				primitive.E{Key: "localField", Value: "destination"},
				primitive.E{Key: "foreignField", Value: "key"},
				primitive.E{Key: "as", Value: "destination"},
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

func BuildAggregation(queries ...[]bson.D) mongo.Pipeline {
	var pipeline mongo.Pipeline

	for _, group := range queries {
		for _, query := range group {
			pipeline = append(pipeline, query)
		}
	}

	return pipeline
}
