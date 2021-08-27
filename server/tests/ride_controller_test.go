package tests

import (
	"example.com/country-roads/controllers"
	"example.com/country-roads/models"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"net/http/httptest"
	"testing"

	"example.com/country-roads/common"
	"example.com/country-roads/mocks"
	"github.com/golang/mock/gomock"
)

func TestGetRideValidId(t *testing.T) {
	tests := []struct {
		ID   string
		Want int
	}{
		{"5c0a7922c9d89830f4911426", http.StatusOK},
		{"tooshort", http.StatusBadRequest},
	}

	mockCtrl := gomock.NewController(t)
	mockedRideRepository := mocks.NewMockRideRepository(mockCtrl)
	mockedRideRepository.EXPECT().FindOne(gomock.Any(), gomock.Any()).Return(models.Ride{}, nil)
	mockEnv := common.Env{
		Collections: common.CollectionContainer{
			RideCollection: mockedRideRepository,
		},
	}

	controller := controllers.GetRide(&mockEnv)

	for _, tt := range tests {
		testname := fmt.Sprintf("%v, %v", tt.ID, tt.Want)
		t.Run(testname, func(t *testing.T) {
			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)

			req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("?id=%v", tt.ID), nil)
			if err != nil {
				panic(err)
			}
			c.Request = req
			c.Params = gin.Params{
				gin.Param{"id", tt.ID},
			}
			controller(c)

			if w.Result().StatusCode != tt.Want {
				t.Errorf("got %v, want %v", w.Result().StatusCode, tt.Want)
			}
		})
	}
}
