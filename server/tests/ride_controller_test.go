package tests

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"example.com/country-roads/controllers"
	"example.com/country-roads/mocks"
	"example.com/country-roads/models"
	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func TestGetRide(t *testing.T) {
	now := time.Now()

	tests := []struct {
		Url          string
		Prepare      func(finder *mocks.MockRideFinder)
		ExpectedCode int
		ExpectedBody gin.H
	}{
		{
			Url: "/rides/551137c2f9e1fac808a5f572",
			Prepare: func(finder *mocks.MockRideFinder) {
				objID, _ := primitive.ObjectIDFromHex("551137c2f9e1fac808a5f572")
				finder.EXPECT().FindMany(gomock.Any(), gomock.Any()).Return(models.Rides{
					models.Ride{
						ID:        objID,
						Direction: "to_campus",
						Type:      "offer",
						Date:      now,
						CreatedAt: now,
						Destination: models.Location{
							Display:   "Kadıköy",
							ID:        objID,
							Key:       "kadikoy",
							ParentKey: "istanbul_asia",
							Parent:    nil,
						},
					},
				}, nil)
			},
			ExpectedCode: http.StatusOK,
			ExpectedBody: gin.H{"results": gin.H{
				"id":        "551137c2f9e1fac808a5f572",
				"direction": "to_campus",
				"type":      "offer",
				"date":      now.Unix(),
				"createdAt": now.Unix(),
				"destination": gin.H{
					"id":        "551137c2f9e1fac808a5f572",
					"key":       "kadikoy",
					"parentKey": "istanbul_asia",
					"display":   "Kadıköy",
				},
			}},
		},
		{
			Url: "/rides/551137c2f9e1fac808a5f572",
			Prepare: func(finder *mocks.MockRideFinder) {
				finder.EXPECT().FindMany(gomock.Any(), gomock.Any()).Return(models.Rides{}, nil)
			},
			ExpectedCode: http.StatusNotFound,
			ExpectedBody: gin.H{"results": gin.H{}},
		},
		{
			Url: "/rides/tooshort",
			Prepare: func(finder *mocks.MockRideFinder) {
				finder.EXPECT().FindMany(gomock.Any(), gomock.Any()).Times(0)
				finder.EXPECT().FindOne(gomock.Any(), gomock.Any()).Times(0)
			},
			ExpectedCode: http.StatusBadRequest,
			ExpectedBody: gin.H{"error": "Invalid ride id"},
		},
	}

	for _, tt := range tests {
		testName := tt.Url
		t.Run(testName, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			mockedRideFinder := mocks.NewMockRideFinder(ctrl)

			if (tt.Prepare != nil) {
				tt.Prepare(mockedRideFinder)
			}

			recorder := httptest.NewRecorder()
			_, r := gin.CreateTestContext(recorder)
			r.GET("/rides/:id", controllers.GetRide(mockedRideFinder))

			request, err := http.NewRequest(http.MethodGet, tt.Url, nil)
			if err != nil {
				t.Fatal(err)
			}
			r.ServeHTTP(recorder, request)

			if bodyAssertion, err := IsBodyEqual(tt.ExpectedBody, recorder.Result().Body); err != nil {
				t.Fatal(err)
			} else if !bodyAssertion {
				t.Errorf("response bodies don't match")
			}

			if (recorder.Result().StatusCode != tt.ExpectedCode) {
				t.Errorf("got %v, want %v", tt.ExpectedCode, recorder.Result().StatusCode)
			}
		})
	}
}
