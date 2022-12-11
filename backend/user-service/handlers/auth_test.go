package handlers_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
	"time"

	"github.com/Slimo300/MicroservicesChatApp/backend/lib/auth"
	"github.com/Slimo300/MicroservicesChatApp/backend/lib/auth/pb"
	mockdb "github.com/Slimo300/MicroservicesChatApp/backend/user-service/database/mock"
	"github.com/Slimo300/MicroservicesChatApp/backend/user-service/handlers"
	"github.com/Slimo300/MicroservicesChatApp/backend/user-service/models"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
	"gorm.io/gorm"
)

func TestSignIn(t *testing.T) {
	gin.SetMode(gin.TestMode)

	TokenClientMock := new(auth.MockTokenClient)
	TokenClientMock.On("NewPairFromUserID", mock.Anything).Return(&pb.TokenPair{
		AccessToken:  "validAccessToken",
		RefreshToken: "validRefreshToken",
		Error:        "",
	}, nil)

	DBMock := new(mockdb.MockUsersDB)
	DBMock.On("GetUserByEmail", "mal.zein@email.com").Return(models.User{ID: uuid.MustParse("c5904224-deec-4275-83bd-56e4cdeba1ae"), Pass: "$2a$10$6BSuuiaPdRJJF2AygYAfnOGkrKLY2o0wDWbEpebn.9Rk0O95D3hW."}, nil)
	DBMock.On("GetUserByEmail", "mal1.zein@email.com").Return(models.User{}, gorm.ErrRecordNotFound)
	DBMock.On("SignInUser", uuid.MustParse("c5904224-deec-4275-83bd-56e4cdeba1ae")).Return(nil)

	s := handlers.Server{
		DB:           DBMock,
		TokenService: TokenClientMock,
	}

	testCases := []struct {
		desc               string
		data               map[string]string
		expectedStatusCode int
		expectedResponse   interface{}
	}{
		{
			desc:               "loginsuccess",
			data:               map[string]string{"email": "mal.zein@email.com", "password": "test"},
			expectedStatusCode: http.StatusOK,
			expectedResponse:   gin.H{"accessToken": "validAccessToken"},
		},
		{
			desc:               "loginnosuchemail",
			data:               map[string]string{"email": "mal1.zein@email.com", "password": "test"},
			expectedStatusCode: http.StatusBadRequest,
			expectedResponse:   gin.H{"err": "wrong email or password"},
		},
		{
			desc:               "logininvalidpass",
			data:               map[string]string{"email": "mal.zein@email.com", "password": "t2est"},
			expectedStatusCode: http.StatusBadRequest,
			expectedResponse:   gin.H{"err": "wrong email or password"},
		},
	}

	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			requestBody, _ := json.Marshal(tC.data)
			req := httptest.NewRequest(http.MethodPost, "/api/login", bytes.NewBuffer(requestBody))
			w := httptest.NewRecorder()
			_, engine := gin.CreateTestContext(w)
			engine.Handle(http.MethodPost, "/api/login", s.SignIn)
			engine.ServeHTTP(w, req)
			response := w.Result()

			if response.StatusCode != tC.expectedStatusCode {
				t.Errorf("Received Status code %d does not match expected status %d", response.StatusCode, tC.expectedStatusCode)
			}
			var respBody gin.H
			json.NewDecoder(response.Body).Decode(&respBody)
			if !reflect.DeepEqual(respBody, tC.expectedResponse) {
				t.Errorf("Received HTTP response body %+v does not match expected HTTP response Body %+v", respBody, tC.expectedResponse)
			}
		})
	}
}

func TestSignOut(t *testing.T) {
	gin.SetMode(gin.TestMode)

	TokenClientMock := auth.NewMockTokenClient()
	TokenClientMock.On("DeleteUserToken", mock.Anything).Return(nil)

	DBMock := new(mockdb.MockUsersDB)
	DBMock.On("SignOutUser", uuid.MustParse("1c4dccaf-a341-4920-9003-f24e0412f8e0")).Return(nil)
	DBMock.On("SignOutUser", uuid.MustParse("2f8fd072-29d4-470a-9359-b3b0e056bf65")).Return(errors.New("No user with id: 2f8fd072-29d4-470a-9359-b3b0e056bf65"))

	s := handlers.Server{
		DB:           DBMock,
		TokenService: TokenClientMock,
	}

	testCases := []struct {
		desc               string
		id                 string
		expectedStatusCode int
		expectedResponse   interface{}
	}{
		{
			desc:               "logoutsuccess",
			id:                 "1c4dccaf-a341-4920-9003-f24e0412f8e0",
			expectedStatusCode: http.StatusOK,
			expectedResponse:   gin.H{"message": "success"},
		},
		{
			desc:               "logoutnouser",
			id:                 "2f8fd072-29d4-470a-9359-b3b0e056bf65",
			expectedStatusCode: http.StatusBadRequest,
			expectedResponse:   gin.H{"err": "No user with id: 2f8fd072-29d4-470a-9359-b3b0e056bf65"},
		},
	}

	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {

			req := httptest.NewRequest(http.MethodPost, "/api/signout", nil)
			req.AddCookie(&http.Cookie{Name: "refreshToken", Value: "validRefreshToken", Path: "/", Expires: time.Now().Add(time.Hour * 24), Domain: "localhost"})
			w := httptest.NewRecorder()

			_, engine := gin.CreateTestContext(w)
			engine.Use(func(c *gin.Context) {
				c.Set("userID", tC.id)
			})

			engine.Handle(http.MethodPost, "/api/signout", s.SignOutUser)
			engine.ServeHTTP(w, req)
			response := w.Result()

			if response.StatusCode != tC.expectedStatusCode {
				t.Errorf("Received Status code %d does not match expected status %d", response.StatusCode, tC.expectedStatusCode)
			}

			var respBody gin.H
			json.NewDecoder(response.Body).Decode(&respBody)
			if !reflect.DeepEqual(respBody, tC.expectedResponse) {
				t.Errorf("Received HTTP response body %+v does not match expected HTTP response Body %+v", respBody, tC.expectedResponse)
			}
		})
	}
}

func TestRefresh(t *testing.T) {

	gin.SetMode(gin.TestMode)
	TokenClientMock := auth.NewMockTokenClient()
	s := handlers.Server{
		DB:           new(mockdb.MockUsersDB),
		TokenService: TokenClientMock,
	}

	testCases := []struct {
		desc               string
		withCookie         bool
		prepare            func(m *mock.Mock)
		expectedStatusCode int
		expectedResponse   interface{}
	}{
		{
			desc:               "refreshNoCookie",
			withCookie:         false,
			prepare:            func(m *mock.Mock) {},
			expectedStatusCode: http.StatusBadRequest,
			expectedResponse:   gin.H{"err": "No token provided"},
		},
		{
			desc:       "refreshTokenBlacklisted",
			withCookie: true,
			prepare: func(m *mock.Mock) {
				m.On("NewPairFromRefresh", mock.Anything).Return(&pb.TokenPair{
					AccessToken:  "",
					RefreshToken: "",
					Error:        "Token Blacklisted",
				}, nil).Once()
			},
			expectedStatusCode: http.StatusForbidden,
			expectedResponse:   gin.H{"err": "Token Blacklisted"},
		},
		{
			desc:       "refreshOK",
			withCookie: true,
			prepare: func(m *mock.Mock) {
				m.On("NewPairFromRefresh", mock.Anything).Return(&pb.TokenPair{
					AccessToken:  "validAccessToken",
					RefreshToken: "validRefreshToken",
					Error:        "",
				}, nil).Once()
			},
			expectedStatusCode: http.StatusOK,
			expectedResponse:   gin.H{"accessToken": "validAccessToken"},
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			tC.prepare(&TokenClientMock.Mock)

			req := httptest.NewRequest(http.MethodPost, "/api/refresh", nil)
			if tC.withCookie {
				req.AddCookie(&http.Cookie{Name: "refreshToken", Value: "validRefreshToken", Path: "/", Expires: time.Now().Add(time.Hour * 24), Domain: "localhost"})
			}
			w := httptest.NewRecorder()

			_, engine := gin.CreateTestContext(w)

			engine.Handle(http.MethodPost, "/api/refresh", s.RefreshToken)
			engine.ServeHTTP(w, req)
			response := w.Result()

			if response.StatusCode != tC.expectedStatusCode {
				t.Errorf("Received Status code %d does not match expected status %d", response.StatusCode, tC.expectedStatusCode)
			}

			var respBody gin.H
			json.NewDecoder(response.Body).Decode(&respBody)
			if !reflect.DeepEqual(respBody, tC.expectedResponse) {
				t.Errorf("Received HTTP response body %+v does not match expected HTTP response Body %+v", respBody, tC.expectedResponse)
			}
		})
	}
}
