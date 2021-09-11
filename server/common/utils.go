package common

import (
	"encoding/json"
	"github.com/sirupsen/logrus"
	"io"

	"github.com/aliparlakci/country-roads/models"
	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
)

var MinDate int64 = -62135596800

func NoError(err error, callback func()) {
	if err != nil {
		callback()
	}
}

func JsonMarshalNoError(v interface{}) string {
	if bytes, err := json.Marshal(v); err != nil {
		return "cannot marshal"
	} else {
		return string(bytes)
	}
}

func LoggerWithRequestId(c *gin.Context) *logrus.Entry {
	return logrus.WithFields(logrus.Fields{"request_id": c.GetString("request_id"), "ip": c.Request.RemoteAddr})
}

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
	result = result && r.Expected.From == actual.From
	result = result && r.Expected.To == actual.To
	result = result && r.Expected.CreatedAt.Unix() != MinDate
	return result
}

func (r RideSchemaMatcher) String() string {
	return ""
}

type UserSchemaMatcher struct {
	Expected models.UserSchema
}

func GetUserSchemaMatcher(schema models.UserSchema) gomock.Matcher {
	return &UserSchemaMatcher{Expected: schema}
}

func (r UserSchemaMatcher) Matches(x interface{}) bool {
	actual, succ := x.(models.UserSchema)
	if !succ {
		return false
	}

	result := true
	result = result && r.Expected.ID == actual.ID
	result = result && r.Expected.DisplayName == actual.DisplayName
	result = result && r.Expected.Email == actual.Email
	result = result && r.Expected.Verified == actual.Verified
	result = result && r.Expected.SignedUpAt.Unix() != MinDate
	return result
}

func (r UserSchemaMatcher) String() string {
	return ""
}
