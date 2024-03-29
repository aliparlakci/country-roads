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
	if queries.From != "" {
		keys = append(keys, primitive.E{Key: "from", Value: queries.From})
	}
	if queries.To != "" {
		keys = append(keys, primitive.E{Key: "to", Value: queries.To})
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

var RideResponseAggregation = []bson.D{
	bson.D{
		primitive.E{
			Key: "$lookup",
			Value: bson.D{
				primitive.E{Key: "from", Value: "locations"},
				primitive.E{Key: "localField", Value: "from"},
				primitive.E{Key: "foreignField", Value: "key"},
				primitive.E{Key: "as", Value: "from"},
			},
		},
	},
	bson.D{
		primitive.E{
			Key: "$lookup",
			Value: bson.D{
				primitive.E{Key: "from", Value: "locations"},
				primitive.E{Key: "localField", Value: "to"},
				primitive.E{Key: "foreignField", Value: "key"},
				primitive.E{Key: "as", Value: "to"},
			},
		},
	},
	bson.D{
		primitive.E{
			Key: "$lookup",
			Value: bson.D{
				primitive.E{Key: "from", Value: "users"},
				primitive.E{Key: "localField", Value: "owner"},
				primitive.E{Key: "foreignField", Value: "_id"},
				primitive.E{Key: "as", Value: "owner"},
			},
		},
	},
	bson.D{
		primitive.E{
			Key: "$unwind",
			Value: bson.D{
				primitive.E{Key: "path", Value: "$owner"},
				primitive.E{Key: "preserveNullAndEmptyArrays", Value: false},
			},
		},
	},
	bson.D{
		primitive.E{
			Key: "$unwind",
			Value: bson.D{
				primitive.E{Key: "path", Value: "$from"},
				primitive.E{Key: "preserveNullAndEmptyArrays", Value: false},
			},
		},
	}, bson.D{
		primitive.E{
			Key: "$unwind",
			Value: bson.D{
				primitive.E{Key: "path", Value: "$to"},
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
