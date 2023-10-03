package handlers_test

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	mockdb "github.com/Slimo300/Microservices-Videocall-and-Chat-App/backend/user-service/database/mock"
	"github.com/Slimo300/Microservices-Videocall-and-Chat-App/backend/user-service/handlers"
	"github.com/Slimo300/Microservices-Videocall-and-Chat-App/backend/user-service/models"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/suite"
)

type SearchTestSuite struct {
	suite.Suite
	server handlers.Server
}

func (s *SearchTestSuite) SetupSuite() {

	db := new(mockdb.MockUsersDB)
	db.On("GetUserByUsername", "valid").Return(&models.User{UserName: "valid"}, nil)
	db.On("GetUserByUsername", "invalid").Return(nil, errors.New("no such user"))

	s.server = handlers.Server{
		DB: db,
	}
}

func (s *SearchTestSuite) TestSearchUser() {
	gin.SetMode(gin.TestMode)

	testCases := []struct {
		desc               string
		returnVal          bool
		username           string
		expectedStatusCode int
		expectedResponse   interface{}
	}{
		{
			desc:               "searchUserSuccess",
			returnVal:          true,
			username:           "valid",
			expectedStatusCode: http.StatusOK,
			expectedResponse:   models.User{UserName: "valid"},
		},
		{
			desc:               "searchUserNotFound",
			returnVal:          false,
			username:           "invalid",
			expectedStatusCode: http.StatusNotFound,
			expectedResponse:   gin.H{"err": "no such user"},
		},
	}

	for _, tC := range testCases {
		s.Run(tC.desc, func() {
			req, _ := http.NewRequest(http.MethodGet, fmt.Sprintf("/api/search/%s", tC.username), nil)
			w := httptest.NewRecorder()
			_, engine := gin.CreateTestContext(w)

			engine.Handle(http.MethodGet, "/api/search/:name", s.server.SearchUser)
			engine.ServeHTTP(w, req)
			response := w.Result()
			defer response.Body.Close()

			s.Equal(tC.expectedStatusCode, response.StatusCode)

			var respBody interface{}
			if tC.returnVal {
				var user models.User
				if err := json.NewDecoder(response.Body).Decode(&user); err != nil {
					s.Fail(err.Error())
				}
				respBody = user
			} else {
				var msg gin.H
				if err := json.NewDecoder(response.Body).Decode(&msg); err != nil {
					s.Fail(err.Error())
				}
				respBody = msg
			}

			s.Equal(tC.expectedResponse, respBody)
		})
	}
}

func TestSearchSuite(t *testing.T) {
	suite.Run(t, &ProfileTestSuite{})
}
