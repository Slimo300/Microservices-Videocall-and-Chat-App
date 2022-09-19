package handlers_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	limits "github.com/gin-contrib/size"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func TestDeleteGroupProfilePicture(t *testing.T) {
	gin.SetMode(gin.TestMode)
	s := setupTestServer()

	badGroupID := uuid.NewString()

	testCases := []struct {
		desc               string
		userID             string
		groupID            string
		expectedStatusCode int
		expectedResponse   interface{}
	}{
		{
			desc:               "DeleteProfilePictureInvalidUserID",
			userID:             "1c4dccaf-a341-4920-9003-f4e0412f8e0",
			groupID:            "61fbd273-b941-471c-983a-0a3cd2c74747",
			expectedStatusCode: http.StatusBadRequest,
			expectedResponse:   gin.H{"err": "invalid ID"},
		},
		{
			desc:               "DeleteProfilePictureInvalidGroupID",
			userID:             "1c4dccaf-a341-4920-9003-f24e0412f8e0",
			groupID:            "61fbd273-b941-471c-983a-a3cd2c7477",
			expectedStatusCode: http.StatusBadRequest,
			expectedResponse:   gin.H{"err": "invalid group ID"},
		},
		{
			desc:               "DeleteProfilePictureNoMember",
			userID:             "1fa00013-89b1-49ad-af29-a79afea1f8b8",
			groupID:            "61fbd273-b941-471c-983a-0a3cd2c74747",
			expectedStatusCode: http.StatusForbidden,
			expectedResponse:   gin.H{"err": "no rights to set"},
		},
		{
			desc:               "DeleteProfilePictureNoRightsToSet",
			userID:             "0ef41409-24b0-43e6-80a3-cf31a4b1a684",
			groupID:            "61fbd273-b941-471c-983a-0a3cd2c74747",
			expectedStatusCode: http.StatusForbidden,
			expectedResponse:   gin.H{"err": "no rights to set"},
		},
		{
			desc:               "DeleteProfilePictureNoGroup",
			userID:             "0ef41409-24b0-43e6-80a3-cf31a4b1a684",
			groupID:            badGroupID,
			expectedStatusCode: http.StatusBadRequest,
			expectedResponse:   gin.H{"err": "record not found"},
		},
		{
			desc:               "DeleteProfilePictureNoPicture",
			userID:             "1c4dccaf-a341-4920-9003-f24e0412f8e0",
			groupID:            "87a0c639-e590-422e-9664-6aedd5ef85ba",
			expectedStatusCode: http.StatusBadRequest,
			expectedResponse:   gin.H{"err": "group has no image to delete"},
		},
		{
			desc:               "DeleteProfilePictureSuccess",
			userID:             "1c4dccaf-a341-4920-9003-f24e0412f8e0",
			groupID:            "61fbd273-b941-471c-983a-0a3cd2c74747",
			expectedStatusCode: http.StatusOK,
			expectedResponse:   gin.H{"message": "success"},
		},
	}

	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {

			req, _ := http.NewRequest(http.MethodDelete, "/api/group/"+tC.groupID+"/image", nil)
			w := httptest.NewRecorder()
			_, engine := gin.CreateTestContext(w)
			engine.Use(func(c *gin.Context) {
				c.Set("userID", tC.userID)
			})

			engine.Handle(http.MethodDelete, "/api/group/:groupID/image", s.DeleteGroupProfilePicture)
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

func TestSetGroupProfilePicture(t *testing.T) {
	gin.SetMode(gin.TestMode)
	s := setupTestServer()

	testCases := []struct {
		desc               string
		userID             string
		groupID            string
		imageData          map[string]string
		setBodyLimiter     bool
		returnVal          bool
		expectedStatusCode int
		expectedResponse   interface{}
	}{
		{
			desc:               "UpdateProfilePictureInvalidUserID",
			userID:             "1c4dccaf-a341-4920-9003-f4e0412f8e0",
			groupID:            "61fbd273-b941-471c-983a-0a3cd2c74747",
			imageData:          map[string]string{"Key": "avatarFile", "CType": "image/png"},
			setBodyLimiter:     false,
			returnVal:          false,
			expectedStatusCode: http.StatusBadRequest,
			expectedResponse:   gin.H{"err": "Invalid ID"},
		},
		{
			desc:               "UpdateProfilePictureNoFile",
			userID:             "1c4dccaf-a341-4920-9003-f24e0412f8e0",
			groupID:            "61fbd273-b941-471c-983a-0a3cd2c74747",
			imageData:          map[string]string{"Key": "WrongFile", "CType": "image/png"},
			setBodyLimiter:     false,
			returnVal:          false,
			expectedStatusCode: http.StatusBadRequest,
			expectedResponse:   gin.H{"err": "http: no such file"},
		},
		{
			desc:               "UpdateProfilePictureInvalidGroupID",
			userID:             "1c4dccaf-a341-4920-9003-f24e0412f8e0",
			groupID:            "61fbd273-b941-471c-983a-a3cd2c7477",
			imageData:          map[string]string{"Key": "avatarFile", "CType": "image/png"},
			setBodyLimiter:     false,
			returnVal:          false,
			expectedStatusCode: http.StatusBadRequest,
			expectedResponse:   gin.H{"err": "Invalid group ID"},
		},
		{
			desc:               "UpdateProfilePictureNoMember",
			userID:             "1fa00013-89b1-49ad-af29-a79afea1f8b8",
			groupID:            "61fbd273-b941-471c-983a-0a3cd2c74747",
			imageData:          map[string]string{"Key": "avatarFile", "CType": "image/png"},
			setBodyLimiter:     false,
			returnVal:          false,
			expectedStatusCode: http.StatusForbidden,
			expectedResponse:   gin.H{"err": "no rights to set"},
		},
		{
			desc:               "UpdateProfilePictureNoRightsToSet",
			userID:             "0ef41409-24b0-43e6-80a3-cf31a4b1a684",
			groupID:            "61fbd273-b941-471c-983a-0a3cd2c74747",
			imageData:          map[string]string{"Key": "avatarFile", "CType": "image/png"},
			setBodyLimiter:     false,
			returnVal:          false,
			expectedStatusCode: http.StatusForbidden,
			expectedResponse:   gin.H{"err": "no rights to set"},
		},
		{
			desc:               "UpdateProfilePictureWrongImageType",
			userID:             "1c4dccaf-a341-4920-9003-f24e0412f8e0",
			groupID:            "61fbd273-b941-471c-983a-0a3cd2c74747",
			imageData:          map[string]string{"Key": "avatarFile", "CType": "application/octet-stream"},
			setBodyLimiter:     false,
			returnVal:          false,
			expectedStatusCode: http.StatusBadRequest,
			expectedResponse:   gin.H{"err": "image extention not allowed"},
		},
		{
			desc:               "UpdateProfilePictureTooBig",
			userID:             "1c4dccaf-a341-4920-9003-f24e0412f8e0",
			groupID:            "61fbd273-b941-471c-983a-0a3cd2c74747",
			imageData:          map[string]string{"Key": "avatarFile", "CType": "image/png"},
			setBodyLimiter:     true,
			returnVal:          false,
			expectedStatusCode: http.StatusRequestEntityTooLarge,
			expectedResponse:   gin.H{"err": "too large"},
		},
		{
			desc:               "UpdateProfilePictureSuccess",
			userID:             "1c4dccaf-a341-4920-9003-f24e0412f8e0",
			groupID:            "61fbd273-b941-471c-983a-0a3cd2c74747",
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

			req, _ := http.NewRequest(http.MethodPut, "/api/group/"+tC.groupID+"/image", body)
			req.Header.Add("Content-Type", writer.FormDataContentType())

			w := httptest.NewRecorder()
			_, engine := gin.CreateTestContext(w)
			engine.Use(func(c *gin.Context) {
				c.Set("userID", tC.userID)
			})

			if tC.setBodyLimiter {
				engine.Use(limits.RequestSizeLimiter(10))
			}
			engine.Handle(http.MethodPut, "/api/group/:groupID/image", s.SetGroupProfilePicture)
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
