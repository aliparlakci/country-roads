package middlewares

import (
	"errors"
	"fmt"
	"github.com/aliparlakci/country-roads/mocks"
	"github.com/aliparlakci/country-roads/services"
	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestSessionMiddlewareSetCookie(t *testing.T) {
	tests := []struct {
		Cookies string
		ExpectedSetCookie string
		Prepare func(repo *mocks.MockSessionRepository)
	}{
		{
			Cookies: "",
			Prepare: func(repo *mocks.MockSessionRepository) {
				newSession := services.Session{}
				repo.EXPECT().CreateSession(gomock.Any(), newSession).Return("agreathash", nil)
				repo.EXPECT().FetchSession(gomock.Any(), gomock.Any()).Times(0)
			},
			ExpectedSetCookie: "visitor=agreathash; Max-Age=7776000;",
		},
		{
			Cookies: "visitor=itsmeyaboy",
			Prepare: func(repo *mocks.MockSessionRepository) {
				repo.EXPECT().CreateSession(gomock.Any(), gomock.Any()).Times(0)
				repo.EXPECT().FetchSession(gomock.Any(), "itsmeyaboy").Return(services.Session{}, nil)
			},
			ExpectedSetCookie: "",
		},
		{
			Cookies: "visitor=idontknowu",
			Prepare: func(repo *mocks.MockSessionRepository) {
				repo.EXPECT().FetchSession(gomock.Any(), "idontknowu").Return(services.Session{}, errors.New("new phone who dis"))
				repo.EXPECT().CreateSession(gomock.Any(), services.Session{}).Return("itsmeyaboy", nil)
			},
			ExpectedSetCookie: "visitor=itsmeyaboy; Max-Age=7776000;",
		},
	}

	for i, tt := range tests {
		testName := fmt.Sprintf("[%v]", i)
		t.Run(testName, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			mockedSessionRepository := mocks.NewMockSessionRepository(ctrl)

			if tt.Prepare != nil {
				tt.Prepare(mockedSessionRepository)
			}

			recorder := httptest.NewRecorder()
			_, r := gin.CreateTestContext(recorder)
			r.Use(SessionMiddleware(mockedSessionRepository))

			request, err := http.NewRequest(http.MethodPost, "/", nil)
			if err != nil {
				t.Fatal(err)
			}
			request.Header.Set("cookie", tt.Cookies)

			r.ServeHTTP(recorder, request)

			if header := recorder.Result().Header.Get("Set-Cookie"); header != tt.ExpectedSetCookie {
				t.Errorf("want %v, got %v", tt.ExpectedSetCookie, header)
			}
		})
	}
}

func TestSessionMiddlewareUpdateContext(t *testing.T) {
	//TODO: Implement
}