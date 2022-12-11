package handlers_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	mockdb "github.com/Slimo300/MicroservicesChatApp/backend/user-service/database/mock"
	"github.com/Slimo300/MicroservicesChatApp/backend/user-service/email"
	"github.com/Slimo300/MicroservicesChatApp/backend/user-service/handlers"
	"github.com/Slimo300/MicroservicesChatApp/backend/user-service/models"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
	"github.com/thanhpk/randstr"
)

func TestRegister(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockEmailService := email.NewMockEmailService()
	mockEmailService.On("SendVerificationEmail", mock.Anything).Return(nil)

	mockDB := new(mockdb.MockUsersDB)
	mockDB.On("IsUsernameInDatabase", "johnny").Return(true)
	mockDB.On("IsUsernameInDatabase", "johnny1").Return(false)
	mockDB.On("IsEmailInDatabase", "johnny@net.com").Return(true)
	mockDB.On("IsEmailInDatabase", "johnny1@net.com").Return(false)
	mockDB.On("RegisterUser", mock.Anything).Return(models.User{ID: uuid.New(), Email: "johnny@net.pl", UserName: "johnny", Pass: "password", Verified: false}, nil)
	mockDB.On("NewVerificationCode", mock.Anything, mock.Anything).Return(models.VerificationCode{UserID: uuid.New(), ActivationCode: randstr.String(10)}, nil)

	s := handlers.Server{
		DB:           mockDB,
		EmailService: mockEmailService,
	}

	testCases := []struct {
		desc               string
		data               map[string]string
		expectedStatusCode int
		expectedResponse   interface{}
		prepare            func(dbMock *mock.Mock, emailMock *mock.Mock)
	}{
		{
			desc:               "registersuccess",
			data:               map[string]string{"username": "johnny1", "email": "johnny1@net.com", "password": "password"},
			expectedStatusCode: http.StatusCreated,
			expectedResponse:   gin.H{"message": "success"},
		},
		{
			desc:               "registerusernametaken",
			data:               map[string]string{"username": "johnny", "email": "johnny@net.com", "password": "password"},
			expectedStatusCode: http.StatusConflict,
			expectedResponse:   gin.H{"err": "username taken"},
		},
		{
			desc:               "registeremailtaken",
			data:               map[string]string{"username": "johnny1", "email": "johnny@net.com", "password": "password"},
			expectedStatusCode: http.StatusConflict,
			expectedResponse:   gin.H{"err": "email already in database"},
		},
		{
			desc:               "registerinvalidpass",
			data:               map[string]string{"username": "johnny1", "email": "johnny1@net.com", "password": ""},
			expectedStatusCode: http.StatusBadRequest,
			expectedResponse:   gin.H{"err": "not a valid password"},
		},
		{
			desc:               "registerinvalidemail",
			data:               map[string]string{"username": "johnny1", "email": "johnny1@net.com2", "password": "password"},
			expectedStatusCode: http.StatusBadRequest,
			expectedResponse:   gin.H{"err": "not a valid email"},
		},
		{
			desc:               "registerinvalidusername",
			data:               map[string]string{"username": "j", "email": "johnny1@net.com", "password": "password"},
			expectedStatusCode: http.StatusBadRequest,
			expectedResponse:   gin.H{"err": "not a valid username"},
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {

			requestBody, _ := json.Marshal(tC.data)

			req, _ := http.NewRequest("POST", "/api/register", bytes.NewBuffer(requestBody))

			w := httptest.NewRecorder()
			_, engine := gin.CreateTestContext(w)
			engine.Handle(http.MethodPost, "/api/register", s.RegisterUser)
			engine.ServeHTTP(w, req)
			response := w.Result()

			if response.StatusCode != tC.expectedStatusCode {
				t.Errorf("Received Status code %d does not match expected status %d", response.StatusCode, tC.expectedStatusCode)
			}
			var errmsg gin.H
			json.NewDecoder(response.Body).Decode(&errmsg)
			if !reflect.DeepEqual(errmsg, tC.expectedResponse) {
				t.Errorf("Received HTTP response body %+v does not match expected HTTP response Body %+v", errmsg, tC.expectedResponse)
			}
		})
	}
}
