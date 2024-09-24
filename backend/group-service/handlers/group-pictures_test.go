package handlers_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Slimo300/Microservices-Videocall-and-Chat-App/backend/group-service/handlers"
	mockservice "github.com/Slimo300/Microservices-Videocall-and-Chat-App/backend/group-service/service/mock"
	limits "github.com/gin-contrib/size"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

type GroupPicturesTestSuite struct {
	suite.Suite
	IDs    map[string]uuid.UUID
	server *handlers.Server
}

func (s *GroupPicturesTestSuite) SetupSuite() {

	s.IDs = make(map[string]uuid.UUID)

	s.IDs["user"] = uuid.New()
	s.IDs["group"] = uuid.New()

	service := new(mockservice.GroupsMockService)
	service.On("SetGroupPicture", mock.Anything, s.IDs["userOK"], s.IDs["groupOK"], mock.Anything).Return("picture_url", nil)
	service.On("DeleteGroupPicture", mock.Anything, s.IDs["userOK"], s.IDs["groupOK"], mock.Anything).Return(nil)

	s.server = handlers.NewServer(service, nil)
}

func (s *GroupPicturesTestSuite) TestDeleteGroupProfilePicture() {
	gin.SetMode(gin.TestMode)

	testCases := []struct {
		desc               string
		userID             string
		groupID            string
		expectedStatusCode int
		expectedResponse   interface{}
	}{
		{
			desc:               "DeleteProfilePictureInvalidUserID",
			userID:             s.IDs["userOK"].String()[:2],
			groupID:            s.IDs["groupOK"].String(),
			expectedStatusCode: http.StatusBadRequest,
			expectedResponse:   gin.H{"err": "invalid ID"},
		},
		{
			desc:               "DeleteProfilePictureInvalidGroupID",
			userID:             s.IDs["userOK"].String(),
			groupID:            s.IDs["groupOK"].String()[:2],
			expectedStatusCode: http.StatusBadRequest,
			expectedResponse:   gin.H{"err": "invalid group ID"},
		},
		{
			desc:               "DeleteProfilePictureSuccess",
			userID:             s.IDs["userOK"].String(),
			groupID:            s.IDs["groupOK"].String(),
			expectedStatusCode: http.StatusOK,
			expectedResponse:   gin.H{"message": "success"},
		},
	}

	for _, tC := range testCases {
		s.Run(tC.desc, func() {

			req, _ := http.NewRequest(http.MethodDelete, "/api/group/"+tC.groupID+"/image", nil)
			w := httptest.NewRecorder()
			_, engine := gin.CreateTestContext(w)
			engine.Use(func(c *gin.Context) {
				c.Set("userID", tC.userID)
			})

			engine.Handle(http.MethodDelete, "/api/group/:groupID/image", s.server.DeleteGroupProfilePicture)
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

func (s *GroupPicturesTestSuite) TestSetGroupProfilePicture() {
	gin.SetMode(gin.TestMode)

	testCases := []struct {
		desc               string
		userID             string
		groupID            string
		imageData          map[string]string
		setBodyLimiter     bool
		expectedStatusCode int
		expectedResponse   interface{}
	}{
		{
			desc:               "UpdateProfilePictureInvalidUserID",
			userID:             s.IDs["userOK"].String()[:2],
			groupID:            s.IDs["groupOK"].String(),
			imageData:          map[string]string{"Key": "avatarFile", "CType": "image/png"},
			setBodyLimiter:     false,
			expectedStatusCode: http.StatusBadRequest,
			expectedResponse:   gin.H{"err": "Invalid ID"},
		},
		{
			desc:               "UpdateProfilePictureInvalidGroupID",
			userID:             s.IDs["userOK"].String(),
			groupID:            s.IDs["groupOK"].String()[:2],
			imageData:          map[string]string{"Key": "avatarFile", "CType": "image/png"},
			setBodyLimiter:     false,
			expectedStatusCode: http.StatusBadRequest,
			expectedResponse:   gin.H{"err": "Invalid group ID"},
		},
		{
			desc:               "UpdateProfilePictureNoFile",
			userID:             s.IDs["userOK"].String(),
			groupID:            s.IDs["groupOK"].String(),
			imageData:          map[string]string{"Key": "WrongFile", "CType": "image/png"},
			setBodyLimiter:     false,
			expectedStatusCode: http.StatusBadRequest,
			expectedResponse:   gin.H{"err": "http: no such file"},
		},
		{
			desc:               "UpdateProfilePictureWrongImageType",
			userID:             s.IDs["userOK"].String(),
			groupID:            s.IDs["groupOK"].String(),
			imageData:          map[string]string{"Key": "avatarFile", "CType": "application/octet-stream"},
			setBodyLimiter:     false,
			expectedStatusCode: http.StatusBadRequest,
			expectedResponse:   gin.H{"err": "image extension not allowed"},
		},
		{
			desc:               "UpdateProfilePictureTooBig",
			userID:             s.IDs["userOK"].String(),
			groupID:            s.IDs["groupOK"].String(),
			imageData:          map[string]string{"Key": "avatarFile", "CType": "image/png"},
			setBodyLimiter:     true,
			expectedStatusCode: http.StatusRequestEntityTooLarge,
		},
		{
			desc:               "UpdateProfilePictureSuccess",
			userID:             s.IDs["userOK"].String(),
			groupID:            s.IDs["groupOK"].String(),
			imageData:          map[string]string{"Key": "avatarFile", "CType": "image/png"},
			setBodyLimiter:     false,
			expectedStatusCode: http.StatusOK,
			expectedResponse:   gin.H{"newUrl": "picture_url"},
		},
	}

	for _, tC := range testCases {
		s.Run(tC.desc, func() {

			body, writer, err := createTestFormFile(tC.imageData["Key"], tC.imageData["CType"])
			if err != nil {
				s.Fail("error when creating form file: %v", err)
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
			engine.Handle(http.MethodPut, "/api/group/:groupID/image", s.server.SetGroupProfilePicture)
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

			if !tC.setBodyLimiter {
				s.Equal(tC.expectedResponse, respBody)
			}
		})
	}
}

func TestGroupPicturesSuite(t *testing.T) {
	suite.Run(t, &GroupPicturesTestSuite{})
}
