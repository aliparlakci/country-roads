package tests

import (
	"example.com/country-roads/controllers"
	"example.com/country-roads/mocks"
	"example.com/country-roads/models"
	"example.com/country-roads/validators"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"go.mongodb.org/mongo-driver/bson"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestPostLocation(t *testing.T) {
	tests := []struct {
		Body         multipart.Form
		Prepare      func(inserter *mocks.MockLocationInserter, locationFinder *mocks.MockLocationFinder)
		ExpectedCode int
		ExpectedBody gin.H
	}{
		{
			Body: multipart.Form{Value: map[string][]string{
				"key":       {"taksim"},
				"display":   {"Taksim"},
				"parentKey": {"istanbul_europe"},
			}},
			Prepare: func(inserter *mocks.MockLocationInserter, locationFinder *mocks.MockLocationFinder) {
				locationFinder.EXPECT().FindOne(gomock.Any(), bson.M{"key": "istanbul_europe"}).Return(models.Location{}, nil).MinTimes(1)
				locationFinder.EXPECT().FindOne(gomock.Any(), bson.M{"key": "taksim"}).Return(models.Location{}, nil).MinTimes(1)
				inserter.EXPECT().InsertOne(gomock.Any(), gomock.Any()).Times(0)
			},
			ExpectedCode: http.StatusBadRequest,
			ExpectedBody: gin.H{"error": "Location format was invalid"},
		},
		{
			Body: multipart.Form{Value: map[string][]string{
				"key":     {"taksim"},
				"display": {"Taksim"},
			}},
			Prepare: func(inserter *mocks.MockLocationInserter, locationFinder *mocks.MockLocationFinder) {
				locationFinder.EXPECT().FindOne(gomock.Any(), gomock.Any()).Times(0)
				locationFinder.EXPECT().FindOne(gomock.Any(), bson.M{"key": "taksim"}).Return(models.Location{}, fmt.Errorf("")).MinTimes(1)
				inserter.EXPECT().InsertOne(gomock.Any(), models.LocationSchema{
					Key:     "taksim",
					Display: "Taksim",
				}).Return("551137c2f9e1fac808a5f572", nil)
			},
			ExpectedCode: http.StatusCreated,
			ExpectedBody: gin.H{"id": "551137c2f9e1fac808a5f572"},
		},
		{
			Body: multipart.Form{Value: map[string][]string{
				"key":       {"taksim"},
				"display":   {"Taksim"},
				"parentKey": {"istanbul_europe"},
			}},
			Prepare: func(inserter *mocks.MockLocationInserter, locationFinder *mocks.MockLocationFinder) {
				locationFinder.EXPECT().FindOne(gomock.Any(), bson.M{"key": "istanbul_europe"}).Return(models.Location{}, nil).MinTimes(1)
				locationFinder.EXPECT().FindOne(gomock.Any(), bson.M{"key": "taksim"}).Return(models.Location{}, fmt.Errorf("")).MinTimes(1)
				inserter.EXPECT().InsertOne(gomock.Any(), models.LocationSchema{
					Key:       "taksim",
					Display:   "Taksim",
					ParentKey: "istanbul_europe",
				}).Return("551137c2f9e1fac808a5f572", nil)
			},
			ExpectedCode: http.StatusCreated,
			ExpectedBody: gin.H{"id": "551137c2f9e1fac808a5f572"},
		},
	}

	for i, tt := range tests {
		testName := fmt.Sprintf("[%v]", i)
		t.Run(testName, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			mockedLocationInserter := mocks.NewMockLocationInserter(ctrl)
			mockedLocationFinder := mocks.NewMockLocationFinder(ctrl)
			validator := validators.ValidatorFactory{LocationFinder: mockedLocationFinder}

			if tt.Prepare != nil {
				tt.Prepare(mockedLocationInserter, mockedLocationFinder)
			}

			recorder := httptest.NewRecorder()
			_, r := gin.CreateTestContext(recorder)
			r.POST("/rides", controllers.PostLocation(mockedLocationInserter, validator))

			request, err := http.NewRequest(http.MethodPost, "/rides", nil)
			request.MultipartForm = &tt.Body
			request.Header.Set("Content-Type", "multipart/form-data")

			if err != nil {
				t.Fatal(err)
			}

			r.ServeHTTP(recorder, request)

			if bodyAssertion, err := IsBodyEqual(tt.ExpectedBody, recorder.Result().Body); err != nil {
				t.Fatal(err)
			} else if !bodyAssertion {
				t.Errorf("response bodies don't match")
			}

			if recorder.Result().StatusCode != tt.ExpectedCode {
				t.Errorf("want %v, got %v", tt.ExpectedCode, recorder.Result().StatusCode)
			}
		})
	}
}

func TestGetAllLocations(t *testing.T) {
	tests := []struct {
		Prepare              func(finder *mocks.MockLocationFinder)
		ExpectedCode         int
		ExpectedResultLength int
	}{
		{
			Prepare: func(finder *mocks.MockLocationFinder) {
				finder.EXPECT().FindMany(gomock.Any(), gomock.Any()).Return(models.Locations{
					models.Location{},
					models.Location{},
					models.Location{},
				}, nil)
			},
			ExpectedCode:         http.StatusOK,
			ExpectedResultLength: 3,
		},
	}

	for i, tt := range tests {
		testName := fmt.Sprintf("[%v]", i)
		t.Run(testName, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			mockedLocationFinder := mocks.NewMockLocationFinder(ctrl)

			if tt.Prepare != nil {
				tt.Prepare(mockedLocationFinder)
			}

			recorder := httptest.NewRecorder()
			_, r := gin.CreateTestContext(recorder)
			r.GET("/locations", controllers.GetAllLocations(mockedLocationFinder))

			request, err := http.NewRequest(http.MethodGet, "/locations", nil)
			if err != nil {
				t.Fatal(err)
			}
			r.ServeHTTP(recorder, request)

			if bodyAssertion, err := IsResultsSameLength(tt.ExpectedResultLength, recorder.Result().Body); err != nil {
				t.Fatal(err)
			} else if !bodyAssertion {
				t.Errorf("response bodies don't match")
			}

			if recorder.Result().StatusCode != tt.ExpectedCode {
				t.Errorf("got %v, want %v", tt.ExpectedCode, recorder.Result().StatusCode)
			}
		})
	}
}
