package tests

import (
	"errors"
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
	"time"
)

func TestSignUp(t *testing.T) {
	tests := []struct {
		Body         multipart.Form
		Prepare      func(findInserter *mocks.MockUserFindInserter)
		ExpectedCode int
		ExpectedBody gin.H
	}{
		{
			Body: multipart.Form{Value: map[string][]string{
				"displayName": {"Ali"},
				"email":       {"aliparlakci@sabanciuniv.edu"},
				"phone":       {"+905423538751"},
			}},
			Prepare: func(findInserter *mocks.MockUserFindInserter) {
				findInserter.EXPECT().FindOne(gomock.Any(), bson.M{
					"email": "aliparlakci@sabanciuniv.edu",
				}).Return(models.User{}, errors.New(""))
				findInserter.EXPECT().InsertOne(gomock.Any(), GetUserSchemaMatcher(models.UserSchema{
					DisplayName: "Ali",
					Email:       "aliparlakci@sabanciuniv.edu",
					Phone:       "+905423538751",
					Verified:    false,
					SignedUpAt:  time.Now(),
				})).Return("551137c2f9e1fac808a5f572", nil)
			},
			ExpectedCode: http.StatusCreated,
			ExpectedBody: gin.H{"id": "551137c2f9e1fac808a5f572"},
		},
		{
			Body: multipart.Form{Value: map[string][]string{
				"displayName": {"Ali"},
				"email":       {"aliparlakci@sabanciuniv.edu"},
				"phone":       {"+905423538751"},
			}},
			Prepare: func(findInserter *mocks.MockUserFindInserter) {
				findInserter.EXPECT().FindOne(gomock.Any(), bson.M{
					"email": "aliparlakci@sabanciuniv.edu",
				}).Return(models.User{}, nil)
				findInserter.EXPECT().InsertOne(gomock.Any(), GetUserSchemaMatcher(models.UserSchema{
					DisplayName: "Ali",
					Email:       "aliparlakci@sabanciuniv.edu",
					Phone:       "+905423538751",
					Verified:    false,
					SignedUpAt:  time.Now(),
				})).Times(0)
			},
			ExpectedCode: http.StatusBadRequest,
			ExpectedBody: gin.H{"error": "user already exists"},
		},
		{
			Body: multipart.Form{Value: map[string][]string{
				"displayName": {"Ali"},
				"email":       {"@sabanciuniv.edu"},
				"phone":       {"+905423538751"},
			}},
			Prepare: func(findInserter *mocks.MockUserFindInserter) {
				findInserter.EXPECT().FindOne(gomock.Any(), bson.M{
					"email": "@sabanciuniv.edu",
				}).Return(models.User{}, errors.New("")).MaxTimes(1)
				findInserter.EXPECT().InsertOne(gomock.Any(), GetUserSchemaMatcher(models.UserSchema{
					DisplayName: "Ali",
					Email:       "@sabanciuniv.edu",
					Phone:       "+905423538751",
					Verified:    false,
					SignedUpAt:  time.Now(),
				})).Times(0)
			},
			ExpectedCode: http.StatusBadRequest,
			ExpectedBody: gin.H{"error": "user is not valid"},
		},
		{
			Body: multipart.Form{Value: map[string][]string{
				"displayName": {"Ali"},
				"email":       {"aliparlakci@sabanciuniv.edu"},
				"phone":       {"5423538751"},
			}},
			Prepare: func(findInserter *mocks.MockUserFindInserter) {
				findInserter.EXPECT().FindOne(gomock.Any(), bson.M{
					"email": "aliparlakci@sabanciuniv.edu",
				}).Return(models.User{}, nil).MaxTimes(1)
				findInserter.EXPECT().InsertOne(gomock.Any(), GetUserSchemaMatcher(models.UserSchema{
					DisplayName: "Ali",
					Email:       "aliparlakci@sabanciuniv.edu",
					Phone:       "+905423538751",
					Verified:    false,
					SignedUpAt:  time.Now(),
				})).Times(0)
			},
			ExpectedCode: http.StatusBadRequest,
			ExpectedBody: gin.H{"error": "user is not valid"},
		},
	}

	for i, tt := range tests {
		testName := fmt.Sprintf("[%v]", i)
		t.Run(testName, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			mockUserFindInserter := mocks.NewMockUserFindInserter(ctrl)
			validator := validators.ValidatorFactory{}

			if tt.Prepare != nil {
				tt.Prepare(mockUserFindInserter)
			}

			recorder := httptest.NewRecorder()
			_, r := gin.CreateTestContext(recorder)
			r.POST("/users", controllers.PostUser(mockUserFindInserter, validator))

			request, err := http.NewRequest(http.MethodPost, "/users", nil)
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
