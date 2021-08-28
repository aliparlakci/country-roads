package tests

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"example.com/country-roads/controllers"
	"example.com/country-roads/mocks"
	"example.com/country-roads/models"
	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func TestGetRidesInvalidID(t *testing.T) {
	tests := []struct {
		Params gin.Params
		Want   int
	}{
		{Params: gin.Params{gin.Param{Key: "id", Value: "tooshort"}}, Want: http.StatusBadRequest},
		{Params: gin.Params{}, Want: http.StatusBadRequest},
	}

	controller := controllers.GetRides(nil)

	for _, tt := range tests {
		testId, _ := tt.Params.Get("id")
		testname := fmt.Sprintf("%v, %v", testId, tt.Want)
		t.Run(testname, func(t *testing.T) {
			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)

			c.Params = tt.Params
			controller(c)

			if w.Result().StatusCode != tt.Want {
				t.Errorf("got %v, want %v", w.Result().StatusCode, tt.Want)
			}
		})
	}
}

func TestGetRidesSuccess(t *testing.T) {
	param := gin.Param{Key: "id", Value: "551137c2f9e1fac808a5f572"}
	id, err := primitive.ObjectIDFromHex(param.Value)
	if err != nil {
		t.Fatal(err)
	}
	want := models.Ride{ID: id}

	mockCtrl := gomock.NewController(t)
	mockedRideFinder := mocks.NewMockRideFinder(mockCtrl)
	mockedRideFinder.EXPECT().FindOne(gomock.Any(), gomock.Any()).Return(want, nil)

	controller := controllers.GetRides(mockedRideFinder)

	testname := fmt.Sprint(param.Value)
	t.Run(testname, func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)

		c.Params = gin.Params{param}
		controller(c)

		result := w.Result()

		rawBody, err := BodyReader(result.Body)
		if err != nil {
			t.Fatal(err)
		}

		var resultBody gin.H
		err = json.Unmarshal(rawBody, &resultBody)
		if err != nil {
			t.Fatal(err)
		}

		if result.StatusCode != http.StatusOK {
			t.Errorf("got %v, want %v", w.Result().StatusCode, want)
		}

		jsonResultBody, _ := json.Marshal(resultBody)
		jsonWantBody, _ := json.Marshal(want.Jsonify())

		if string(jsonResultBody) != string(jsonWantBody) {
			t.Errorf("got %v, want %v", w.Result().StatusCode, want)
		}
	})
}

func TestGetRidesNotFound(t *testing.T) {
	param := gin.Param{Key: "id", Value: "551137c2f9e1fac808a5f572"}
	want := http.StatusNotFound

	mockCtrl := gomock.NewController(t)
	mockedRideFinder := mocks.NewMockRideFinder(mockCtrl)
	mockedRideFinder.EXPECT().FindOne(gomock.Any(), gomock.Any()).Return(models.Ride{}, fmt.Errorf(""))

	controller := controllers.GetRides(mockedRideFinder)

	testname := fmt.Sprint(param.Value)
	t.Run(testname, func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)

		c.Params = gin.Params{param}
		controller(c)
		result := w.Result()

		if result.StatusCode != want {
			t.Errorf("got %v, want %v", w.Result().StatusCode, want)
		}
	})
}
