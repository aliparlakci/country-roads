package tests

import (
	"encoding/json"
	"fmt"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"testing"

	"example.com/country-roads/controllers"
	"example.com/country-roads/mocks"
	"example.com/country-roads/models"
	"example.com/country-roads/validators"
	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func TestGetLocationsSuccess(t *testing.T) {
	id, err := primitive.ObjectIDFromHex("612620d35d526dc43e342e30")
	if err != nil {
		t.Fatal(err)
	}
	parentId, err := primitive.ObjectIDFromHex("61261e7d5d526dc43e342e21")
	if err != nil {
		t.Fatal(err)
	}
	want := models.Locations{
		models.Location{ID: id, ParentID: parentId, Display: "4. Levent"},
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
			t.Errorf("got %v, want %v", w.Result().StatusCode, want)
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

func TestPostLocation(t *testing.T) {
	tests := []struct {
		Form multipart.Form
		Want int
	}{
		{Form: multipart.Form{Value: map[string][]string{"display": {"Ankara"}}}, Want: http.StatusCreated},
		{Form: multipart.Form{Value: map[string][]string{"display": {"Ankara"}, "parentId": {"612620d35d526dc43e342e30"}}}, Want: http.StatusCreated},
	}

	mockCtrl := gomock.NewController(t)
	mockLocationInserter := mocks.NewMockLocationInserter(mockCtrl)
	mockLocationInserter.EXPECT().InsertOne(gomock.Any(), gomock.Any()).Return(nil, nil).AnyTimes()

	mockLocationValidator := mocks.NewMockValidator(mockCtrl)
	mockLocationValidator.EXPECT().Validate(gomock.Any()).Return(true, nil).AnyTimes()
	mockLocationValidator.EXPECT().SetDto(gomock.Any()).Return().AnyTimes()

	controller := controllers.PostLocation(mockLocationInserter, func() validators.Validator { return mockLocationValidator })

	for _, tt := range tests {
		testId := tt.Form
		testname := fmt.Sprintf("%v, %v", testId, tt.Want)
		t.Run(testname, func(t *testing.T) {
			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)

			req := &http.Request{MultipartForm: &tt.Form, Header: map[string][]string{"Content-Type": {"multipart/form-data"}}}
			c.Request = req
			controller(c)

			if w.Result().StatusCode != tt.Want {
				t.Errorf("got %v, want %v", w.Result().StatusCode, tt.Want)
			}
		})
	}
}
