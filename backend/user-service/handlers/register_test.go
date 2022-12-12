package handlers_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Slimo300/MicroservicesChatApp/backend/lib/apperrors"
	mockqueue "github.com/Slimo300/MicroservicesChatApp/backend/lib/msgqueue/mock"
	mockdb "github.com/Slimo300/MicroservicesChatApp/backend/user-service/database/mock"
	"github.com/Slimo300/MicroservicesChatApp/backend/user-service/email"
	"github.com/Slimo300/MicroservicesChatApp/backend/user-service/handlers"
	"github.com/Slimo300/MicroservicesChatApp/backend/user-service/models"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

type RegisterSuite struct {
	suite.Suite
	server handlers.Server
	ids    map[string]uuid.UUID
}

func (s *RegisterSuite) SetupSuite() {
	s.ids = make(map[string]uuid.UUID)
	s.ids["userOK"] = uuid.MustParse("1819f7ba-01fd-4d79-aa4c-09db1a481f94")

	db := new(mockdb.MockUsersDB)
	db.On("RegisterUser", models.User{Email: "host@net.com", UserName: "taken", Pass: "password12"}).Return(nil, nil, apperrors.NewConflict("username", "taken"))
	db.On("RegisterUser", models.User{Email: "taken@net.com", UserName: "slimo300", Pass: "password12"}).Return(nil, nil, apperrors.NewConflict("email", "taken@net.com"))
	db.On("RegisterUser", models.User{Email: "host@net.com", UserName: "slimo300", Pass: "password12"}).Return(nil, nil, nil)

	db.On("VerifyCode", "invalidCode").Return(nil, apperrors.NewNotFound("code", "invalidCode"))
	db.On("VerifyCode", "validCode").Return(&models.User{ID: s.ids["userOK"]}, nil)

	emitter := new(mockqueue.MockEmitter)
	emitter.On("Emit", mock.Anything).Return(nil)

	emailService := new(email.MockEmailService)
	emailService.On("SendVerificationEmail", mock.Anything).Return(nil)

	s.server = handlers.Server{
		DB:           db,
		Emitter:      emitter,
		EmailService: emailService,
	}
}

func (s RegisterSuite) TestRegister() {
	gin.SetMode(gin.TestMode)

	testCases := []struct {
		desc               string
		data               map[string]string
		expectedStatusCode int
		expectedResponse   interface{}
	}{
		{
			desc:               "registerInvalidEmail",
			data:               map[string]string{"username": "slimo300", "email": "host@net.com2", "password": "password12", "rpassword": "password12"},
			expectedStatusCode: http.StatusBadRequest,
			expectedResponse:   gin.H{"err": "not a valid email"},
		},
		{
			desc:               "registerInvalidUsername",
			data:               map[string]string{"username": "s", "email": "host@net.com", "password": "password12", "rpassword": "password12"},
			expectedStatusCode: http.StatusBadRequest,
			expectedResponse:   gin.H{"err": "not a valid username"},
		},
		{
			desc:               "registerInvalidPass",
			data:               map[string]string{"username": "slimo300", "email": "host@net.com", "password": "pass", "rpassword": "pass"},
			expectedStatusCode: http.StatusBadRequest,
			expectedResponse:   gin.H{"err": "not a valid password"},
		},
		{
			desc:               "registerPassDontMatch",
			data:               map[string]string{"username": "slimo300", "email": "host@net.com", "password": "password12", "rpassword": "password123"},
			expectedStatusCode: http.StatusBadRequest,
			expectedResponse:   gin.H{"err": "passwords don't match"},
		},
		{
			desc:               "registerusernametaken",
			data:               map[string]string{"username": "taken", "email": "host@net.com", "password": "password12", "rpassword": "password12"},
			expectedStatusCode: http.StatusConflict,
			expectedResponse:   gin.H{"err": "resource: username with value: taken already exists"},
		},
		{
			desc:               "registeremailtaken",
			data:               map[string]string{"username": "slimo300", "email": "taken@net.com", "password": "password12", "rpassword": "password12"},
			expectedStatusCode: http.StatusConflict,
			expectedResponse:   gin.H{"err": "resource: email with value: taken@net.com already exists"},
		},
		{
			desc:               "registersuccess",
			data:               map[string]string{"username": "slimo300", "email": "host@net.com", "password": "password12", "rpassword": "password12"},
			expectedStatusCode: http.StatusCreated,
			expectedResponse:   gin.H{"message": "success"},
		},
	}
	for _, tC := range testCases {
		s.Run(tC.desc, func() {

			requestBody, _ := json.Marshal(tC.data)

			req, _ := http.NewRequest("POST", "/api/register", bytes.NewBuffer(requestBody))

			w := httptest.NewRecorder()
			_, engine := gin.CreateTestContext(w)
			engine.Handle(http.MethodPost, "/api/register", s.server.RegisterUser)
			engine.ServeHTTP(w, req)
			response := w.Result()

			s.Equal(tC.expectedStatusCode, response.StatusCode)
			var msg gin.H
			json.NewDecoder(response.Body).Decode(&msg)
			s.Equal(tC.expectedResponse, msg)
		})
	}
}

func (s RegisterSuite) TestVerifyEmail() {
	gin.SetMode(gin.TestMode)

	testCases := []struct {
		desc               string
		code               string
		expectedStatusCode int
		expectedResponse   interface{}
	}{
		{
			desc:               "verifyCodeNotFound",
			code:               "invalidCode",
			expectedStatusCode: http.StatusNotFound,
			expectedResponse:   gin.H{"err": "resource: code with value: invalidCode not found"},
		},
		{
			desc:               "verifySuccess",
			code:               "validCode",
			expectedStatusCode: http.StatusOK,
			expectedResponse:   gin.H{"message": "code verified"},
		},
	}

	for _, tC := range testCases {
		s.Run(tC.desc, func() {
			req, _ := http.NewRequest(http.MethodGet, "/api/verify-account/"+tC.code, nil)

			w := httptest.NewRecorder()
			_, engine := gin.CreateTestContext(w)
			engine.Handle(http.MethodGet, "/api/verify-account/:code", s.server.VerifyCode)
			engine.ServeHTTP(w, req)

			response := w.Result()

			s.Equal(tC.expectedStatusCode, response.StatusCode)

			var msg gin.H
			json.NewDecoder(response.Body).Decode(&msg)
			s.Equal(tC.expectedResponse, msg)
		})
	}
}

func TestRegisterSuite(t *testing.T) {
	suite.Run(t, &RegisterSuite{})
}
