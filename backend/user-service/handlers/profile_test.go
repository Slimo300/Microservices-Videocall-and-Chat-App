package handlers_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/Slimo300/MicroservicesChatApp/backend/lib/storage"
	mockdb "github.com/Slimo300/MicroservicesChatApp/backend/user-service/database/mock"
	"github.com/Slimo300/MicroservicesChatApp/backend/user-service/handlers"
	"github.com/Slimo300/MicroservicesChatApp/backend/user-service/models"
	limits "github.com/gin-contrib/size"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
	"gorm.io/gorm"
)

func TestChangePassword(t *testing.T) {
	gin.SetMode(gin.TestMode)

	DBMock := new(mockdb.DBLayerMock)
	DBMock.On("GetUserById", uuid.MustParse("c71b4a02-85fb-4092-9e98-051302728eaf")).Return(models.User{
		ID:   uuid.MustParse("c71b4a02-85fb-4092-9e98-051302728eaf"),
		Pass: "$2a$10$6BSuuiaPdRJJF2AygYAfnOGkrKLY2o0wDWbEpebn.9Rk0O95D3hW."}, nil)
	DBMock.On("GetUserById", uuid.MustParse("5fc8ab07-cc10-40cc-a84e-9c850309b038")).Return(models.User{}, errors.New("No user with id: 5fc8ab07-cc10-40cc-a84e-9c850309b038"))
	DBMock.On("SetPassword", uuid.MustParse("c71b4a02-85fb-4092-9e98-051302728eaf"), mock.AnythingOfType("string")).Return(nil)
	s := handlers.Server{
		DB: DBMock,
	}

	testCases := []struct {
		desc               string
		userID             string
		data               map[string]interface{}
		expectedStatusCode int
		expectedResponse   interface{}
	}{
		{
			desc:               "changePasswordInvalidID",
			userID:             "1c4dccaf-a341-4920-9003-f4e0412f8e0",
			data:               map[string]interface{}{"oldPassword": "test", "newPassword": "test12"},
			expectedStatusCode: http.StatusBadRequest,
			expectedResponse:   gin.H{"err": "invalid ID"},
		},
		{
			desc:               "changePasswordPassTooShort",
			userID:             "c71b4a02-85fb-4092-9e98-051302728eaf",
			data:               map[string]interface{}{"oldPassword": "test", "newPassword": "test1"},
			expectedStatusCode: http.StatusBadRequest,
			expectedResponse:   gin.H{"err": "Password must be at least 6 characters long"},
		},
		{
			desc:               "changePasswordNoUser",
			userID:             "5fc8ab07-cc10-40cc-a84e-9c850309b038",
			data:               map[string]interface{}{"oldPassword": "test", "newPassword": "test12"},
			expectedStatusCode: http.StatusBadRequest,
			expectedResponse:   gin.H{"err": "No user with id: 5fc8ab07-cc10-40cc-a84e-9c850309b038"},
		},
		{
			desc:               "changePasswordPassDontMatch",
			userID:             "c71b4a02-85fb-4092-9e98-051302728eaf",
			data:               map[string]interface{}{"oldPassword": "test1", "newPassword": "test12"},
			expectedStatusCode: http.StatusForbidden,
			expectedResponse:   gin.H{"err": "Wrong password"},
		},
		{
			desc:               "changePasswordSuccess",
			userID:             "c71b4a02-85fb-4092-9e98-051302728eaf",
			data:               map[string]interface{}{"oldPassword": "test", "newPassword": "test12"},
			expectedStatusCode: http.StatusOK,
			expectedResponse:   gin.H{"message": "password changed"},
		},
	}

	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {

			requestBody, _ := json.Marshal(tC.data)
			req, _ := http.NewRequest(http.MethodPut, "/api/change-password", bytes.NewBuffer(requestBody))
			w := httptest.NewRecorder()
			_, engine := gin.CreateTestContext(w)
			engine.Use(func(c *gin.Context) {
				c.Set("userID", tC.userID)
			})

			engine.Handle(http.MethodPut, "/api/change-password", s.ChangePassword)
			engine.ServeHTTP(w, req)
			response := w.Result()

			if response.StatusCode != tC.expectedStatusCode {
				t.Errorf("Received Status code %d does not match expected status %d", response.StatusCode, tC.expectedStatusCode)
			}

			var msg gin.H
			json.NewDecoder(response.Body).Decode(&msg)

			if !reflect.DeepEqual(msg, tC.expectedResponse) {
				t.Errorf("Received HTTP response body %+v does not match expected HTTP response Body %+v", msg, tC.expectedResponse)
			}
		})
	}
}

func TestDeleteProfilePicture(t *testing.T) {
	gin.SetMode(gin.TestMode)

	DBMock := new(mockdb.DBLayerMock)
	DBMock.On("GetProfilePictureURL", uuid.MustParse("0ef41409-24b0-43e6-80a3-cf31a4b1a684")).Return("", nil)
	DBMock.On("GetProfilePictureURL", uuid.MustParse("f586fa1a-af84-4a2e-9fc6-1a4ada270fe4")).Return("", gorm.ErrRecordNotFound)
	DBMock.On("GetProfilePictureURL", uuid.MustParse("1c4dccaf-a341-4920-9003-f24e0412f8e0")).Return("url", nil)
	DBMock.On("SetProfilePicture", uuid.MustParse("1c4dccaf-a341-4920-9003-f24e0412f8e0"), mock.AnythingOfType("string")).Return(nil)

	s := handlers.Server{
		DB:           DBMock,
		ImageStorage: new(storage.MockStorage),
	}

	testCases := []struct {
		desc               string
		userID             string
		expectedStatusCode int
		expectedResponse   interface{}
	}{
		{
			desc:               "DeleteProfilePicturePInvalidID",
			userID:             "1c4dccaf-a341-4920-9003-f4e12f8e0",
			expectedStatusCode: http.StatusBadRequest,
			expectedResponse:   gin.H{"err": "invalid ID"},
		},
		{
			desc:               "DeleteProfilePictureNoPicture",
			userID:             "0ef41409-24b0-43e6-80a3-cf31a4b1a684",
			expectedStatusCode: http.StatusBadRequest,
			expectedResponse:   gin.H{"err": "user has no image to delete"},
		},
		{
			desc:               "DeleteProfilePictureNoUser",
			userID:             "f586fa1a-af84-4a2e-9fc6-1a4ada270fe4",
			expectedStatusCode: http.StatusBadRequest,
			expectedResponse:   gin.H{"err": "User not found"},
		},
		{
			desc:               "DeleteProfilePictureSuccess",
			userID:             "1c4dccaf-a341-4920-9003-f24e0412f8e0",
			expectedStatusCode: http.StatusOK,
			expectedResponse:   gin.H{"message": "success"},
		},
	}

	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {

			req, _ := http.NewRequest(http.MethodDelete, "/api/delete-image", nil)
			w := httptest.NewRecorder()
			_, engine := gin.CreateTestContext(w)
			engine.Use(func(c *gin.Context) {
				c.Set("userID", tC.userID)
			})

			engine.Handle(http.MethodDelete, "/api/delete-image", s.DeleteProfilePicture)
			engine.ServeHTTP(w, req)
			response := w.Result()

			if response.StatusCode != tC.expectedStatusCode {
				t.Errorf("Received Status code %d does not match expected status %d", response.StatusCode, tC.expectedStatusCode)
			}

			var msg gin.H
			json.NewDecoder(response.Body).Decode(&msg)

			if !reflect.DeepEqual(msg, tC.expectedResponse) {
				t.Errorf("Received HTTP response body %+v does not match expected HTTP response Body %+v", msg, tC.expectedResponse)
			}
		})
	}
}

func TestSetProfilePicture(t *testing.T) {
	gin.SetMode(gin.TestMode)

	DBMock := new(mockdb.DBLayerMock)
	DBMock.On("GetProfilePictureURL", uuid.MustParse("1c4dccaf-a341-4920-9003-f24e0412f8e0")).Return("someUrl", nil)

	s := handlers.Server{
		DB:           DBMock,
		ImageStorage: new(storage.MockStorage),
	}

	testCases := []struct {
		desc               string
		userID             string
		imageData          map[string]string
		setBodyLimiter     bool
		returnVal          bool
		expectedStatusCode int
		expectedResponse   interface{}
	}{
		{
			desc:               "UpdateProfilePictureInvalidUserID",
			userID:             "1c4dccaf-a341-4920-9003-f4e0412f8e0",
			imageData:          map[string]string{"Key": "avatarFile", "CType": "image/png"},
			setBodyLimiter:     false,
			returnVal:          false,
			expectedStatusCode: http.StatusBadRequest,
			expectedResponse:   gin.H{"err": "invalid ID"},
		},
		{
			desc:               "UpdateProfilePictureNoFile",
			userID:             "1c4dccaf-a341-4920-9003-f24e0412f8e0",
			imageData:          map[string]string{"Key": "WrongFile", "CType": "image/png"},
			setBodyLimiter:     false,
			returnVal:          false,
			expectedStatusCode: http.StatusBadRequest,
			expectedResponse:   gin.H{"err": "http: no such file"},
		},
		{
			desc:               "UpdateProfilePictureWrongImageType",
			userID:             "1c4dccaf-a341-4920-9003-f24e0412f8e0",
			imageData:          map[string]string{"Key": "avatarFile", "CType": "application/octet-stream"},
			setBodyLimiter:     false,
			returnVal:          false,
			expectedStatusCode: http.StatusBadRequest,
			expectedResponse:   gin.H{"err": "image extention not allowed"},
		},
		{
			desc:               "UpdateProfilePictureTooBig",
			userID:             "1c4dccaf-a341-4920-9003-f24e0412f8e0",
			imageData:          map[string]string{"Key": "avatarFile", "CType": "image/png"},
			setBodyLimiter:     true,
			returnVal:          false,
			expectedStatusCode: http.StatusRequestEntityTooLarge,
			expectedResponse:   gin.H{"err": "too large"},
		},
		{
			desc:               "UpdateProfilePictureSuccess",
			userID:             "1c4dccaf-a341-4920-9003-f24e0412f8e0",
			imageData:          map[string]string{"Key": "avatarFile", "CType": "image/png"},
			setBodyLimiter:     false,
			returnVal:          true,
			expectedStatusCode: http.StatusOK,
			expectedResponse:   nil,
		},
	}

	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {

			body, writer, err := createTestFormFile(tC.imageData["Key"], tC.imageData["CType"])
			if err != nil {
				t.Errorf("error when creating form file: %v", err)
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
			engine.Handle(http.MethodPut, "/api/set-image", s.UpdateProfilePicture)
			engine.ServeHTTP(w, req)
			response := w.Result()

			if response.StatusCode != tC.expectedStatusCode {
				t.Errorf("Received Status code %d does not match expected status %d", response.StatusCode, tC.expectedStatusCode)
			}

			var msg gin.H
			json.NewDecoder(response.Body).Decode(&msg)

			if tC.setBodyLimiter {
				// expecting empty response
			} else if !tC.returnVal {
				if !reflect.DeepEqual(msg, tC.expectedResponse) {
					t.Errorf("Received HTTP response body %+v does not match expected HTTP response Body %+v", msg, tC.expectedResponse)
				}
			} else {
				if msg["newUrl"] == "" {
					t.Errorf("Received HTTP response body %+v is not set", tC.expectedResponse)
				}
			}
		})
	}
}
