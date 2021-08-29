package tests

import (
	"encoding/json"
	"example.com/country-roads/common"
	"example.com/country-roads/models"
	_ "fmt"
	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"io"
)

func BodyReader(reader io.Reader) ([]byte, error) {
	rawBody := make([]byte, 0)
	chunk := make([]byte, 8)
	for {
		n, err := reader.Read(chunk)
		if err == io.EOF {
			break
		} else if err != nil {
			return nil, err
		}
		rawBody = append(rawBody, chunk[:n]...)
	}

	return rawBody, nil
}

func IsBodyEqual(expected gin.H, actual io.Reader) (bool, error) {
	actualBytes, err := BodyReader(actual)
	if err != nil {
		return false, err
	}

	parsedExpected := expected

	var parsedActual gin.H
	err = json.Unmarshal(actualBytes, &parsedActual)
	if err != nil {
		return false, nil
	}

	marshalledExpected, _ := json.Marshal(parsedExpected)
	marshalledActual, _ := json.Marshal(parsedActual)

	return string(marshalledExpected) == string(marshalledActual), nil
}

func IsResultsSameLength(expectedLength int, actual io.Reader) (bool, error) {
	actualBytes, err := BodyReader(actual)
	if err != nil {
		return false, err
	}

	var parsedActual gin.H
	err = json.Unmarshal(actualBytes, &parsedActual)
	if err != nil {
		return false, nil
	}

	results, success := parsedActual["results"].([]interface{})
	if !success {
		return false, nil
	}

	return expectedLength == len(results), nil
}

type RideSchemaMatcher struct {
	Expected models.RideSchema
}

func GetRideSchemaMatcher(schema models.RideSchema) gomock.Matcher {
	return &RideSchemaMatcher{Expected: schema}
}

func (r RideSchemaMatcher) Matches(x interface{}) bool {
	actual, succ := x.(models.RideSchema)
	if !succ {
		return false
	}

	result := true
	result = result && r.Expected.ID == actual.ID
	result = result && r.Expected.Type == actual.Type
	result = result && r.Expected.Date == actual.Date
	result = result && r.Expected.Destination == actual.Destination
	result = result && r.Expected.Direction == actual.Direction
	result = result && r.Expected.CreatedAt.Unix() != common.MinDate
	return result
}

func (r RideSchemaMatcher) String() string {
	return ""
}


