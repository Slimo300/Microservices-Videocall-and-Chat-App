package handlers_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Slimo300/Microservices-Videocall-and-Chat-App/backend/lib/apperrors"
	mockqueue "github.com/Slimo300/Microservices-Videocall-and-Chat-App/backend/lib/msgqueue/mock"
	mockdb "github.com/Slimo300/Microservices-Videocall-and-Chat-App/backend/user-service/database/mock"
	"github.com/Slimo300/Microservices-Videocall-and-Chat-App/backend/user-service/handlers"
	"github.com/Slimo300/Microservices-Videocall-and-Chat-App/backend/user-service/models"
	"github.com/Slimo300/Microservices-Videocall-and-Chat-App/backend/user-service/storage"
	limits "github.com/gin-contrib/size"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

type ProfileTestSuite struct {
	suite.Suite
	server handlers.Server
	ids    map[string]uuid.UUID
}

func (s *ProfileTestSuite) SetupSuite() {

	s.ids = make(map[string]uuid.UUID)
	s.ids["userOK"] = uuid.MustParse("48fe113f-da82-4eb2-9944-3212bfcce63e")
	s.ids["userNotFound"] = uuid.MustParse("56aabce7-1142-4e7e-a2a2-622226c1b0d5")
	s.ids["userNoImage"] = uuid.MustParse("5099885a-9c28-42f6-b8a4-f8feeceab579")

	db := new(mockdb.MockUsersDB)
	db.On("GetUserById", s.ids["userOK"]).Return(&models.User{ID: s.ids["userOK"]}, nil)
	db.On("GetUserById", s.ids["userNotFound"]).Return(&models.User{}, errors.New("no such user"))

	db.On("ChangePassword", s.ids["userNotFound"], mock.Anything, mock.Anything).Return(apperrors.NewAuthorization("User not in database"))
	db.On("ChangePassword", s.ids["userOK"], "password", mock.Anything).Return(apperrors.NewForbidden("Wrong Password"))
	db.On("ChangePassword", s.ids["userOK"], "password12", mock.Anything).Return(nil)

	db.On("DeleteProfilePicture", s.ids["userNotFound"]).Return("", apperrors.NewAuthorization("User not found"))
	db.On("DeleteProfilePicture", s.ids["userNoImage"]).Return("", apperrors.NewBadRequest("User has no profile picture"))
	db.On("DeleteProfilePicture", s.ids["userOK"]).Return("picuteURL", nil)

	db.On("GetProfilePictureURL", s.ids["userNotFound"]).Return("", false, apperrors.NewAuthorization("User not found"))
	db.On("GetProfilePictureURL", s.ids["userOK"]).Return("pictureURL", false, nil)

	imageStorage := new(storage.MockStorage)
	imageStorage.On("UploadFile", mock.Anything, mock.Anything).Return(nil)
	imageStorage.On("DeleteFile", mock.Anything).Return(nil)

	emiter := new(mockqueue.MockEmitter)
	emiter.On("Emit", mock.Anything).Return(nil)

	s.server = handlers.Server{
		DB:           db,
		ImageStorage: imageStorage,
		Emitter:      emiter,
	}
}

func (s *ProfileTestSuite) TestGetUser() {
	gin.SetMode(gin.TestMode)

	testCases := []struct {
		desc               string
		returnVal          bool
		userID             string
		expectedStatusCode int
		expectedResponse   interface{}
	}{
		{
			desc:               "getUserSuccess",
			returnVal:          true,
			userID:             s.ids["userOK"].String(),
			expectedStatusCode: http.StatusOK,
			expectedResponse:   models.User{ID: s.ids["userOK"]},
		},
		{
			desc:               "getUserNotFound",
			returnVal:          false,
			userID:             s.ids["userNotFound"].String(),
			expectedStatusCode: http.StatusNotFound,
			expectedResponse:   gin.H{"err": "no such user"},
		},
		{
			desc:               "getUserInvalidID",
			returnVal:          false,
			userID:             "1",
			expectedStatusCode: http.StatusBadRequest,
			expectedResponse:   gin.H{"err": "invalid ID"},
		},
	}

	for _, tC := range testCases {
		s.Run(tC.desc, func() {
			req, _ := http.NewRequest(http.MethodGet, "/api/user", nil)
			w := httptest.NewRecorder()
			_, engine := gin.CreateTestContext(w)
			engine.Use(func(c *gin.Context) {
				c.Set("userID", tC.userID)
			})

			engine.Handle(http.MethodGet, "/api/user", s.server.GetUser)
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

func (s *ProfileTestSuite) TestChangePassword() {
	gin.SetMode(gin.TestMode)

	testCases := []struct {
		desc               string
		userID             string
		data               map[string]interface{}
		expectedStatusCode int
		expectedResponse   interface{}
	}{
		{
			desc:               "changePasswordInvalidID",
			userID:             "1",
			data:               map[string]interface{}{"oldPassword": "password12", "newPassword": "password123", "repeatPassword": "password123"},
			expectedStatusCode: http.StatusBadRequest,
			expectedResponse:   gin.H{"err": "invalid ID"},
		},
		{
			desc:               "changePasswordPassTooShort",
			userID:             s.ids["userOK"].String(),
			data:               map[string]interface{}{"oldPassword": "password12", "newPassword": "pass", "repeatPassword": "pass"},
			expectedStatusCode: http.StatusBadRequest,
			expectedResponse:   gin.H{"err": "Password must be at least 6 characters long"},
		},
		{
			desc:               "changePasswordPassDontMatch",
			userID:             s.ids["userOK"].String(),
			data:               map[string]interface{}{"oldPassword": "password12", "newPassword": "password123", "repeatPassword": "password12"},
			expectedStatusCode: http.StatusBadRequest,
			expectedResponse:   gin.H{"err": "Passwords don't match"},
		},
		{
			desc:               "changePasswordNoUser",
			userID:             s.ids["userNotFound"].String(),
			data:               map[string]interface{}{"oldPassword": "password12", "newPassword": "password123", "repeatPassword": "password123"},
			expectedStatusCode: http.StatusUnauthorized,
			expectedResponse:   gin.H{"err": "User not in database"},
		},
		{
			desc:               "changePasswordWrongPassword",
			userID:             s.ids["userOK"].String(),
			data:               map[string]interface{}{"oldPassword": "password", "newPassword": "password123", "repeatPassword": "password123"},
			expectedStatusCode: http.StatusForbidden,
			expectedResponse:   gin.H{"err": "Wrong Password"},
		},
		{
			desc:               "changePasswordSuccess",
			userID:             s.ids["userOK"].String(),
			data:               map[string]interface{}{"oldPassword": "password12", "newPassword": "password123", "repeatPassword": "password123"},
			expectedStatusCode: http.StatusOK,
			expectedResponse:   gin.H{"message": "password changed"},
		},
	}

	for _, tC := range testCases {
		s.Run(tC.desc, func() {

			requestBody, _ := json.Marshal(tC.data)
			req, _ := http.NewRequest(http.MethodPut, "/api/change-password", bytes.NewBuffer(requestBody))
			w := httptest.NewRecorder()
			_, engine := gin.CreateTestContext(w)
			engine.Use(func(c *gin.Context) {
				c.Set("userID", tC.userID)
			})

			engine.Handle(http.MethodPut, "/api/change-password", s.server.ChangePassword)
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

func (s *ProfileTestSuite) TestDeleteProfilePicture() {
	gin.SetMode(gin.TestMode)

	testCases := []struct {
		desc               string
		userID             string
		expectedStatusCode int
		expectedResponse   interface{}
	}{
		{
			desc:               "DeleteProfilePicturePInvalidID",
			userID:             "1",
			expectedStatusCode: http.StatusBadRequest,
			expectedResponse:   gin.H{"err": "invalid ID"},
		},
		{
			desc:               "DeleteProfilePictureNoPicture",
			userID:             s.ids["userNoImage"].String(),
			expectedStatusCode: http.StatusBadRequest,
			expectedResponse:   gin.H{"err": "User has no profile picture"},
		},
		{
			desc:               "DeleteProfilePictureNoUser",
			userID:             s.ids["userNotFound"].String(),
			expectedStatusCode: http.StatusUnauthorized,
			expectedResponse:   gin.H{"err": "User not found"},
		},
		{
			desc:               "DeleteProfilePictureSuccess",
			userID:             s.ids["userOK"].String(),
			expectedStatusCode: http.StatusOK,
			expectedResponse:   gin.H{"message": "success"},
		},
	}

	for _, tC := range testCases {
		s.Run(tC.desc, func() {

			req, _ := http.NewRequest(http.MethodDelete, "/api/delete-image", nil)
			w := httptest.NewRecorder()
			_, engine := gin.CreateTestContext(w)
			engine.Use(func(c *gin.Context) {
				c.Set("userID", tC.userID)
			})

			engine.Handle(http.MethodDelete, "/api/delete-image", s.server.DeleteProfilePicture)
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

func (s *ProfileTestSuite) TestSetProfilePicture() {
	gin.SetMode(gin.TestMode)

	testCases := []struct {
		desc               string
		userID             string
		imageData          map[string]string
		setBodyLimiter     bool
		expectedStatusCode int
		expectedResponse   interface{}
	}{
		{
			desc:               "UpdateProfilePictureInvalidUserID",
			userID:             "1",
			imageData:          map[string]string{"Key": "avatarFile", "CType": "image/png"},
			setBodyLimiter:     false,
			expectedStatusCode: http.StatusBadRequest,
			expectedResponse:   gin.H{"err": "invalid ID"},
		},
		{
			desc:               "UpdateProfilePictureNoFile",
			userID:             s.ids["userOK"].String(),
			imageData:          map[string]string{"Key": "WrongFile", "CType": "image/png"},
			setBodyLimiter:     false,
			expectedStatusCode: http.StatusBadRequest,
			expectedResponse:   gin.H{"err": "http: no such file"},
		},
		{
			desc:               "UpdateProfilePictureWrongImageType",
			userID:             s.ids["userOK"].String(),
			imageData:          map[string]string{"Key": "avatarFile", "CType": "application/octet-stream"},
			setBodyLimiter:     false,
			expectedStatusCode: http.StatusBadRequest,
			expectedResponse:   gin.H{"err": "image extention not allowed"},
		},
		{
			desc:               "UpdateProfilePictureTooBig",
			userID:             s.ids["userOK"].String(),
			imageData:          map[string]string{"Key": "avatarFile", "CType": "image/png"},
			setBodyLimiter:     true,
			expectedStatusCode: http.StatusRequestEntityTooLarge,
			expectedResponse:   nil,
		},
		{
			desc:               "UpdateProfilePictureNoUser",
			userID:             s.ids["userNotFound"].String(),
			imageData:          map[string]string{"Key": "avatarFile", "CType": "image/png"},
			setBodyLimiter:     false,
			expectedStatusCode: http.StatusUnauthorized,
			expectedResponse:   gin.H{"err": "User not found"},
		},
		{
			desc:               "UpdateProfilePictureSuccess",
			userID:             s.ids["userOK"].String(),
			imageData:          map[string]string{"Key": "avatarFile", "CType": "image/png"},
			setBodyLimiter:     false,
			expectedStatusCode: http.StatusOK,
			expectedResponse:   gin.H{"newUrl": "pictureURL"},
		},
	}

	for _, tC := range testCases {
		s.Run(tC.desc, func() {

			body, writer, err := createTestFormFile(tC.imageData["Key"], tC.imageData["CType"])
			if err != nil {
				s.Fail("error when creating form file: ", err)
			}

			req, _ := http.NewRequest(http.MethodPut, "/api/set-image", body)
			req.Header.Add("Content-Type", writer.FormDataContentType())

			w := httptest.NewRecorder()
			_, engine := gin.CreateTestContext(w)
			engine.Use(func(c *gin.Context) {
				c.Set("userID", tC.userID)
			})

			if tC.setBodyLimiter {
				engine.Use(limits.RequestSizeLimiter(10))
			}
			engine.Handle(http.MethodPut, "/api/set-image", s.server.UpdateProfilePicture)
			engine.ServeHTTP(w, req)
			response := w.Result()
			defer response.Body.Close()

			s.Equal(tC.expectedStatusCode, response.StatusCode)

			var respBody interface{}
			if tC.setBodyLimiter {
				if n, err := response.Body.Read([]byte{}); err != nil || n != 0 {
					s.Fail("Response should be empty whe 413 status is returned")
				}

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

func TestProfileSuite(t *testing.T) {
	suite.Run(t, &ProfileTestSuite{})
}
