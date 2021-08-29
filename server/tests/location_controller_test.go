package tests

import (
	"encoding/json"
	"example.com/country-roads/controllers"
	"example.com/country-roads/mocks"
	"example.com/country-roads/models"
	"example.com/country-roads/validators"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
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
				"key": {"taksim"},
				"display": {"Taksim"},
				"parentKey": {"istanbul_europe"},
			}},
			Prepare: func(inserter *mocks.MockLocationInserter, locationFinder *mocks.MockLocationFinder) {
				locationFinder.EXPECT().FindOne(gomock.Any(), bson.M{"key": "istanbul_europe"}).Return(models.Location{}, nil).MinTimes(1)
				locationFinder.EXPECT().FindOne(gomock.Any(), bson.M{"key": "taksim"}).Return(models.Location{}, fmt.Errorf("")).MinTimes(1)
				inserter.EXPECT().InsertOne(gomock.Any(), models.LocationSchema{
					Key: "taksim",
					Display: "Taksim",
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

func TestGetLocationsSuccess(t *testing.T) {
	id, err := primitive.ObjectIDFromHex("612620d35d526dc43e342e30")
	if err != nil {
		t.Fatal(err)
	}
	want := models.Locations{
		models.Location{ID: id, Key: "levent4", ParentKey: "istanbul_europe", Display: "4. Levent"},
	}
	wantCode := http.StatusOK

	mockCtrl := gomock.NewController(t)
	mockLocationFinder := mocks.NewMockLocationFinder(mockCtrl)
	mockLocationFinder.EXPECT().FindMany(gomock.Any(), gomock.Any()).Return(want, nil)

	controller := controllers.GetAllLocations(mockLocationFinder)

	t.Run("should return 500 when not found", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)

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
			t.Errorf("got %v, want %v", w.Result().StatusCode, wantCode)
		}

		jsonResultBody, _ := json.Marshal(resultBody)
		jsonWantBody, _ := json.Marshal(gin.H{"results": want.Jsonify()})

		if string(jsonResultBody) != string(jsonWantBody) {
			t.Errorf("got %v, want %v", w.Result().StatusCode, wantCode)
		}
	})
}

func TestGetLocationsNotFound(t *testing.T) {
	want := http.StatusNotFound

	mockCtrl := gomock.NewController(t)
	mockLocationFinder := mocks.NewMockLocationFinder(mockCtrl)
	mockLocationFinder.EXPECT().FindMany(gomock.Any(), gomock.Any()).Return([]models.Location{}, nil)

	controller := controllers.GetAllLocations(mockLocationFinder)

	t.Run("should return 404 when not found", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)

		controller(c)
		result := w.Result()

		if result.StatusCode != want {
			t.Errorf("got %v, want %v", w.Result().StatusCode, want)
		}
	})
}

func TestGetLocationsError(t *testing.T) {
	want := http.StatusInternalServerError

	mockCtrl := gomock.NewController(t)
	mockLocationFinder := mocks.NewMockLocationFinder(mockCtrl)
	mockLocationFinder.EXPECT().FindMany(gomock.Any(), gomock.Any()).Return([]models.Location{}, fmt.Errorf(""))

	controller := controllers.GetAllLocations(mockLocationFinder)

	t.Run("should return 500 when error", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)

		controller(c)
		result := w.Result()

		if result.StatusCode != want {
			t.Errorf("got %v, want %v", w.Result().StatusCode, want)
		}
	})
}