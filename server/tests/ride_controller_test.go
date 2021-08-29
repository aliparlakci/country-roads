package tests

import (
	"example.com/country-roads/aggregations"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
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
		{
			Url: "/rides/551137c2f9e1fac808a5f572",
			Prepare: func(finder *mocks.MockRideFinder) {
				finder.EXPECT().FindMany(gomock.Any(), gomock.Any()).Return(nil, fmt.Errorf("cikolatayi severim"))
			},
			ExpectedCode: http.StatusInternalServerError,
			ExpectedBody: gin.H{"error": "cikolatayi severim"},
		},
	}

	for _, tt := range tests {
		testName := tt.Url
		t.Run(testName, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			mockedRideFinder := mocks.NewMockRideFinder(ctrl)

			if tt.Prepare != nil {
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

			if recorder.Result().StatusCode != tt.ExpectedCode {
				t.Errorf("got %v, want %v", tt.ExpectedCode, recorder.Result().StatusCode)
			}
		})
	}
}

func TestSearchRidesInvalid(t *testing.T) {
	tests := []struct {
		Url          string
		Prepare      func(finder *mocks.MockRideFinder)
		ExpectedCode int
	}{
		{
			Url: "/rides?start_date=cannot&end_date=bind",
			Prepare: func(finder *mocks.MockRideFinder) {
				finder.EXPECT().FindMany(gomock.Any(), gomock.Any()).Times(0)
			},
			ExpectedCode: http.StatusBadRequest,
		},
		{
			Url: "/rides",
			Prepare: func(finder *mocks.MockRideFinder) {
				finder.EXPECT().FindMany(gomock.Any(), gomock.Any()).Return(nil, fmt.Errorf(""))
			},
			ExpectedCode: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		testName := tt.Url
		t.Run(testName, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			mockedRideFinder := mocks.NewMockRideFinder(ctrl)

			if tt.Prepare != nil {
				tt.Prepare(mockedRideFinder)
			}

			recorder := httptest.NewRecorder()
			_, r := gin.CreateTestContext(recorder)
			r.GET("/rides", controllers.SearchRides(mockedRideFinder))

			request, err := http.NewRequest(http.MethodGet, tt.Url, nil)
			if err != nil {
				t.Fatal(err)
			}
			r.ServeHTTP(recorder, request)

			if recorder.Result().StatusCode != tt.ExpectedCode {
				t.Errorf("got %v, want %v", tt.ExpectedCode, recorder.Result().StatusCode)
			}
		})
	}
}

func TestSearchRidesMany(t *testing.T) {
	tests := []struct {
		Url                  string
		Prepare              func(finder *mocks.MockRideFinder)
		ExpectedCode         int
		ExpectedResultLength int
	}{
		{
			Url: "/rides?type=offer&direction=to_campus&destination=levent4&start_date=1630227365&end_date=1630227365",
			Prepare: func(finder *mocks.MockRideFinder) {
				queries := models.SearchRideQueries{
					Type:        "offer",
					Direction:   "to_campus",
					Destination: "levent4",
					StartDate:   time.Unix(1630227365, 0),
					EndDate:     time.Unix(1630227365, 0),
				}
				pipeline := aggregations.BuildAggregation(aggregations.FilterRides(queries), aggregations.RideWithDestination)
				finder.EXPECT().FindMany(gomock.Any(), pipeline).Return(models.Rides{
					models.Ride{},
					models.Ride{},
					models.Ride{},
				}, nil)
			},
			ExpectedCode:         http.StatusOK,
			ExpectedResultLength: 3,
		},
		{
			Url: "/rides?it=should&ignore=this&queries=also",
			Prepare: func(finder *mocks.MockRideFinder) {
				queries := models.SearchRideQueries{}
				pipeline := aggregations.BuildAggregation(aggregations.FilterRides(queries), aggregations.RideWithDestination)
				finder.EXPECT().FindMany(gomock.Any(), pipeline).Return(models.Rides{
					models.Ride{},
					models.Ride{},
				}, nil)
			},
			ExpectedCode:         http.StatusOK,
			ExpectedResultLength: 2,
		},
		{
			Url: "/rides",
			Prepare: func(finder *mocks.MockRideFinder) {
				queries := models.SearchRideQueries{}
				pipeline := aggregations.BuildAggregation(aggregations.FilterRides(queries), aggregations.RideWithDestination)
				finder.EXPECT().FindMany(gomock.Any(), pipeline).Return(models.Rides{
					models.Ride{},
					models.Ride{},
					models.Ride{},
					models.Ride{},
				}, nil)
			},
			ExpectedCode:         http.StatusOK,
			ExpectedResultLength: 4,
		},
	}

	for _, tt := range tests {
		testName := tt.Url
		t.Run(testName, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			mockedRideFinder := mocks.NewMockRideFinder(ctrl)

			if tt.Prepare != nil {
				tt.Prepare(mockedRideFinder)
			}

			recorder := httptest.NewRecorder()
			_, r := gin.CreateTestContext(recorder)
			r.GET("/rides", controllers.SearchRides(mockedRideFinder))

			request, err := http.NewRequest(http.MethodGet, tt.Url, nil)
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

func TestDeleteRide(t *testing.T) {
	tests := []struct {
		Url          string
		Prepare      func(finder *mocks.MockRideDeleter)
		ExpectedCode int
		ExpectedBody gin.H
	}{
		{
			Url: "/rides/551137c2f9e1fac808a5f572",
			Prepare: func(finder *mocks.MockRideDeleter) {
				objID, _ := primitive.ObjectIDFromHex("551137c2f9e1fac808a5f572")
				finder.EXPECT().DeleteOne(gomock.Any(), bson.M{"_id": objID}).Return(int64(1), nil)
			},
			ExpectedCode: http.StatusOK,
			ExpectedBody: gin.H{},
		},
		{
			Url: "/rides/tooshort",
			Prepare: func(finder *mocks.MockRideDeleter) {
				finder.EXPECT().DeleteOne(gomock.Any(), gomock.Any()).Times(0)
			},
			ExpectedCode: http.StatusBadRequest,
			ExpectedBody: gin.H{"error": "ride id is not valid"},
		},
		{
			Url: "/rides/551137c2f9e1fac808a5f572",
			Prepare: func(finder *mocks.MockRideDeleter) {
				finder.EXPECT().DeleteOne(gomock.Any(), gomock.Any()).Return(int64(0), fmt.Errorf("cikolatayi severim"))
			},
			ExpectedCode: http.StatusInternalServerError,
			ExpectedBody: gin.H{"error": "cikolatayi severim"},
		},
		{
			Url: "/rides/551137c2f9e1fac808a5f572",
			Prepare: func(finder *mocks.MockRideDeleter) {
				finder.EXPECT().DeleteOne(gomock.Any(), gomock.Any()).Return(int64(0), nil)
			},
			ExpectedCode: http.StatusNotFound,
			ExpectedBody: gin.H{"error": "Ride with ID ObjectID(\"551137c2f9e1fac808a5f572\") does not exist"},
		},
	}

	for i, tt := range tests {
		testName := fmt.Sprintf("[%v] %v", i, tt.Url)
		t.Run(testName, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			mockedRideDeleter := mocks.NewMockRideDeleter(ctrl)

			if tt.Prepare != nil {
				tt.Prepare(mockedRideDeleter)
			}

			recorder := httptest.NewRecorder()
			_, r := gin.CreateTestContext(recorder)
			r.GET("/rides/:id", controllers.DeleteRides(mockedRideDeleter))

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

			if recorder.Result().StatusCode != tt.ExpectedCode {
				t.Errorf("got %v, want %v", tt.ExpectedCode, recorder.Result().StatusCode)
			}
		})
	}
}