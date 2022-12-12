package handlers_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/Slimo300/MicroservicesChatApp/backend/lib/apperrors"
	"github.com/Slimo300/MicroservicesChatApp/backend/lib/auth"
	"github.com/Slimo300/MicroservicesChatApp/backend/lib/auth/pb"
	mockdb "github.com/Slimo300/MicroservicesChatApp/backend/user-service/database/mock"
	"github.com/Slimo300/MicroservicesChatApp/backend/user-service/handlers"
	"github.com/Slimo300/MicroservicesChatApp/backend/user-service/models"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/stretchr/testify/suite"
)

type AuthTestSuite struct {
	suite.Suite
	server handlers.Server
	ids    map[string]uuid.UUID
}

func (s *AuthTestSuite) SetupSuite() {
	s.ids = make(map[string]uuid.UUID)
	s.ids["userOK"] = uuid.MustParse("f2085c84-dadb-4362-accc-44898dedde7a")

	db := new(mockdb.MockUsersDB)
	db.On("SignIn", "host@net.pl", "password12").Return(models.User{ID: s.ids["userOK"]}, nil)
	db.On("SignIn", "host2@net.pl", "password12").Return(models.User{}, apperrors.NewBadRequest("wrong email or password"))
	db.On("SignIn", "host@net.pl", "password").Return(models.User{}, apperrors.NewBadRequest("wrong email or password"))

	tokenService := new(auth.MockTokenClient)

	tokenService.On("NewPairFromUserID", s.ids["userOK"]).Return(&pb.TokenPair{
		AccessToken:  "validAccessToken",
		RefreshToken: "validRefreshToken",
		Error:        "",
	}, nil)
	tokenService.On("DeleteUserToken", "validRefreshToken").Return(nil)
	tokenService.On("DeleteUserToken", "invalidRefreshToken").Return(errors.New("invalid refresh token"))
	tokenService.On("NewPairFromRefresh", "validRefreshToken").Return(&pb.TokenPair{
		AccessToken:  "validAccessToken",
		RefreshToken: "validRefreshToken",
		Error:        "",
	}, nil)
	tokenService.On("NewPairFromRefresh", "expiredRefreshToken").Return(&pb.TokenPair{
		AccessToken:  "",
		RefreshToken: "",
		Error:        "Token Expired",
	}, nil)
	tokenService.On("NewPairFromRefresh", "blacklistedRefreshToken").Return(&pb.TokenPair{
		AccessToken:  "",
		RefreshToken: "",
		Error:        "Token Blacklisted",
	}, nil)

	s.server = handlers.Server{
		DB:           db,
		TokenService: tokenService,
	}

}

func (s *AuthTestSuite) TestSignIn() {
	gin.SetMode(gin.TestMode)

	testCases := []struct {
		desc               string
		data               map[string]string
		expectedStatusCode int
		expectedResponse   interface{}
	}{
		{
			desc:               "loginsuccess",
			data:               map[string]string{"email": "host@net.pl", "password": "password12"},
			expectedStatusCode: http.StatusOK,
			expectedResponse:   gin.H{"accessToken": "validAccessToken"},
		},
		{
			desc:               "loginnosuchemail",
			data:               map[string]string{"email": "host2@net.pl", "password": "password12"},
			expectedStatusCode: http.StatusBadRequest,
			expectedResponse:   gin.H{"err": "Bad request. Reason: wrong email or password"},
		},
		{
			desc:               "logininvalidpass",
			data:               map[string]string{"email": "host@net.pl", "password": "password"},
			expectedStatusCode: http.StatusBadRequest,
			expectedResponse:   gin.H{"err": "Bad request. Reason: wrong email or password"},
		},
	}

	for _, tC := range testCases {
		s.Run(tC.desc, func() {
			requestBody, _ := json.Marshal(tC.data)
			req := httptest.NewRequest(http.MethodPost, "/api/login", bytes.NewBuffer(requestBody))
			w := httptest.NewRecorder()
			_, engine := gin.CreateTestContext(w)
			engine.Handle(http.MethodPost, "/api/login", s.server.SignIn)
			engine.ServeHTTP(w, req)
			response := w.Result()

			s.Equal(tC.expectedStatusCode, response.StatusCode)

			var respBody gin.H
			json.NewDecoder(response.Body).Decode(&respBody)

			s.Equal(tC.expectedResponse, respBody)
		})
	}
}

func (s AuthTestSuite) TestSignOut() {
	gin.SetMode(gin.TestMode)

	testCases := []struct {
		desc               string
		cookiePresent      bool
		cookieValue        string
		expectedStatusCode int
		expectedResponse   interface{}
	}{
		{
			desc:               "logoutsuccess",
			cookiePresent:      true,
			cookieValue:        "validRefreshToken",
			expectedStatusCode: http.StatusOK,
			expectedResponse:   gin.H{"message": "success"},
		},
		{
			desc:               "logoutinvalidtoken",
			cookiePresent:      true,
			cookieValue:        "invalidRefreshToken",
			expectedStatusCode: http.StatusBadRequest,
			expectedResponse:   gin.H{"err": "invalid refresh token"},
		},
		{
			desc:               "logoutnocookie",
			cookiePresent:      false,
			expectedStatusCode: http.StatusBadRequest,
			expectedResponse:   gin.H{"err": "No token to invalidate"},
		},
	}

	for _, tC := range testCases {
		s.Run(tC.desc, func() {

			req := httptest.NewRequest(http.MethodPost, "/api/signout", nil)
			if tC.cookiePresent {
				req.AddCookie(&http.Cookie{Name: "refreshToken", Value: tC.cookieValue, Path: "/", Expires: time.Now().Add(time.Hour * 24), Domain: "localhost"})
			}
			w := httptest.NewRecorder()

			_, engine := gin.CreateTestContext(w)

			engine.Handle(http.MethodPost, "/api/signout", s.server.SignOutUser)
			engine.ServeHTTP(w, req)
			response := w.Result()

			s.Equal(tC.expectedStatusCode, response.StatusCode)

			var respBody gin.H
			json.NewDecoder(response.Body).Decode(&respBody)
			s.Equal(tC.expectedResponse, respBody)
		})
	}
}

func (s AuthTestSuite) TestRefresh() {
	gin.SetMode(gin.TestMode)

	testCases := []struct {
		desc               string
		withCookie         bool
		cookieValue        string
		expectedStatusCode int
		expectedResponse   interface{}
	}{
		{
			desc:               "refreshNoCookie",
			withCookie:         false,
			expectedStatusCode: http.StatusBadRequest,
			expectedResponse:   gin.H{"err": "No token provided"},
		},
		{
			desc:               "refreshTokenExpired",
			withCookie:         true,
			cookieValue:        "expiredRefreshToken",
			expectedStatusCode: http.StatusBadRequest,
			expectedResponse:   gin.H{"err": "Token Expired"},
		},
		{
			desc:               "refreshTokenBlacklisted",
			withCookie:         true,
			cookieValue:        "blacklistedRefreshToken",
			expectedStatusCode: http.StatusForbidden,
			expectedResponse:   gin.H{"err": "Token Blacklisted"},
		},
		{
			desc:               "refreshOK",
			withCookie:         true,
			cookieValue:        "validRefreshToken",
			expectedStatusCode: http.StatusOK,
			expectedResponse:   gin.H{"accessToken": "validAccessToken"},
		},
	}
	for _, tC := range testCases {
		s.Run(tC.desc, func() {

			req := httptest.NewRequest(http.MethodPost, "/api/refresh", nil)
			if tC.withCookie {
				req.AddCookie(&http.Cookie{Name: "refreshToken", Value: tC.cookieValue, Path: "/", Expires: time.Now().Add(time.Hour * 24), Domain: "localhost"})
			}
			w := httptest.NewRecorder()

			_, engine := gin.CreateTestContext(w)

			engine.Handle(http.MethodPost, "/api/refresh", s.server.RefreshToken)
			engine.ServeHTTP(w, req)
			response := w.Result()

			s.Equal(tC.expectedStatusCode, response.StatusCode)

			var respBody gin.H
			json.NewDecoder(response.Body).Decode(&respBody)
			s.Equal(tC.expectedResponse, respBody)
		})
	}
}

func TestUserSuite(t *testing.T) {
	suite.Run(t, &AuthTestSuite{})
}
