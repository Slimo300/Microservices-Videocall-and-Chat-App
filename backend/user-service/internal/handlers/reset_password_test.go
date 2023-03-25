package handlers_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Slimo300/MicroservicesChatApp/backend/lib/apperrors"
	emails "github.com/Slimo300/chat-emailservice/pkg/client"
	dbmock "github.com/Slimo300/chat-userservice/internal/database/mock"
	"github.com/Slimo300/chat-userservice/internal/handlers"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

type ResetPasswordSuite struct {
	suite.Suite
	server handlers.Server
	ids    map[string]uuid.UUID
}

func (s *ResetPasswordSuite) SetupSuite() {
	s.ids = make(map[string]uuid.UUID)

	db := new(dbmock.MockUsersDB)
	db.On("NewResetPasswordCode", "host@net.com").Return(nil, nil, nil)

	db.On("ResetPassword", "invalidCode", mock.Anything).Return(apperrors.NewNotFound("reset code", "invalidCode"))
	db.On("ResetPassword", "validCode", mock.Anything).Return(nil)

	emailClient := new(emails.MockEmailClient)
	emailClient.On("SendResetPasswordEmail", mock.Anything).Return(nil)

	s.server = handlers.Server{
		DB:          db,
		EmailClient: emailClient,
	}
}

func (s *ResetPasswordSuite) TestForgotPassword() {
	gin.SetMode(gin.TestMode)

	testCases := []struct {
		desc               string
		email              string
		expectedStatusCode int
		expectedResponse   interface{}
	}{
		{
			desc:               "forgotInvalidEmail",
			email:              "host@net.com2",
			expectedStatusCode: http.StatusBadRequest,
			expectedResponse:   gin.H{"err": "not a valid email address"},
		},
		{
			desc:               "forgotSuccess",
			email:              "host@net.com",
			expectedStatusCode: http.StatusOK,
			expectedResponse:   gin.H{"message": "reset password email sent"},
		},
	}

	for _, tC := range testCases {

		req, _ := http.NewRequest(http.MethodGet, "/api/forgot-password?email="+tC.email, nil)

		w := httptest.NewRecorder()
		_, engine := gin.CreateTestContext(w)
		engine.Handle(http.MethodGet, "/api/forgot-password", s.server.ForgotPassword)
		engine.ServeHTTP(w, req)

		response := w.Result()
		defer response.Body.Close()

		s.Equal(tC.expectedStatusCode, response.StatusCode)

		var msg gin.H
		if err := json.NewDecoder(response.Body).Decode(&msg); err != nil {
			s.Fail(err.Error())
		}

		s.Equal(tC.expectedResponse, msg)
	}
}

func (s *ResetPasswordSuite) TestResetForgottenPassword() {
	gin.SetMode(gin.TestMode)

	testCases := []struct {
		desc               string
		code               string
		data               map[string]string
		expectedStatusCode int
		expectedResponse   interface{}
	}{
		{
			desc:               "resetInvalidPass",
			code:               "validCode",
			data:               map[string]string{"newPassword": "pass", "repeatPassword": "pass"},
			expectedStatusCode: http.StatusBadRequest,
			expectedResponse:   gin.H{"err": "not a valid password"},
		},
		{
			desc:               "resetPassDontMatch",
			code:               "validCode",
			data:               map[string]string{"newPassword": "password12", "repeatPassword": "password13"},
			expectedStatusCode: http.StatusBadRequest,
			expectedResponse:   gin.H{"err": "passwords don't match"},
		},
		{
			desc:               "resetCodeNotFound",
			code:               "invalidCode",
			data:               map[string]string{"newPassword": "password12", "repeatPassword": "password12"},
			expectedStatusCode: http.StatusNotFound,
			expectedResponse:   gin.H{"err": "resource: reset code with value: invalidCode not found"},
		},
		{
			desc:               "resetSuccess",
			code:               "validCode",
			data:               map[string]string{"newPassword": "password12", "repeatPassword": "password12"},
			expectedStatusCode: http.StatusOK,
			expectedResponse:   gin.H{"message": "password changed"},
		},
	}

	for _, tC := range testCases {
		s.Run(tC.desc, func() {
			requestBody, _ := json.Marshal(tC.data)
			req, _ := http.NewRequest(http.MethodPatch, "/api/reset-password/"+tC.code, bytes.NewBuffer(requestBody))

			w := httptest.NewRecorder()
			_, engine := gin.CreateTestContext(w)
			engine.Handle(http.MethodPatch, "/api/reset-password/:code", s.server.ResetForgottenPassword)
			engine.ServeHTTP(w, req)

			response := w.Result()
			defer response.Body.Close()

			s.Equal(tC.expectedStatusCode, response.StatusCode)

			var msg gin.H
			if err := json.NewDecoder(response.Body).Decode(&msg); err != nil {
				s.Fail(err.Error())
			}

			s.Equal(tC.expectedResponse, msg)
		})
	}
}

func TestResetPasswordSuite(t *testing.T) {
	suite.Run(t, &ResetPasswordSuite{})
}
