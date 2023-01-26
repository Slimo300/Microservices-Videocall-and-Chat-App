package handlers_test

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	dbmock "github.com/Slimo300/MicroservicesChatApp/backend/group-service/database/mock"
	"github.com/Slimo300/MicroservicesChatApp/backend/group-service/handlers"
	"github.com/Slimo300/MicroservicesChatApp/backend/lib/apperrors"
	"github.com/Slimo300/MicroservicesChatApp/backend/lib/storage"
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

	s.IDs["userOK"] = uuid.MustParse("e95bd1fc-ec1f-472c-b7f3-6b39aa7a90c4")
	s.IDs["groupOK"] = uuid.MustParse("4552667f-ea03-4ad3-8757-ea4645c8b4a0")
	s.IDs["userWithoutRights"] = uuid.MustParse("ee2c6112-1114-4d9f-8869-716068ff7159")
	s.IDs["groupWithoutPicture"] = uuid.MustParse("4399b92e-d68a-42be-9a01-a9e098df98d8")

	db := new(dbmock.MockGroupsDB)

	db.On("DeleteGroupProfilePicture", s.IDs["userOK"], s.IDs["groupOK"]).Return("picture_url", nil)
	db.On("DeleteGroupProfilePicture", s.IDs["userWithoutRights"], s.IDs["groupOK"]).
		Return("", apperrors.NewForbidden(fmt.Sprintf("User %v has no rights to set in group %v", s.IDs["userWithoutRights"], s.IDs["groupOK"])))
	db.On("DeleteGroupProfilePicture", s.IDs["userOK"], s.IDs["groupWithoutPicture"]).
		Return("", apperrors.NewForbidden(fmt.Sprintf("group %v has no profile picture", s.IDs["groupWithoutPicture"])))

	db.On("GetGroupProfilePictureURL", s.IDs["userOK"], s.IDs["groupOK"]).Return("picture_url", nil)
	db.On("GetGroupProfilePictureURL", s.IDs["userWithoutRights"], s.IDs["groupOK"]).
		Return("", apperrors.NewForbidden(fmt.Sprintf("User %v has no rights to set in group %v", s.IDs["userWithoutRights"], s.IDs["groupOK"])))

	storage := new(storage.MockStorage)

	storage.On("DeleteFile", mock.Anything).Return(nil)
	storage.On("UploadFile", mock.Anything, mock.Anything).Return(nil)

	s.server = handlers.NewServer(db, storage, nil)
}

func (s GroupPicturesTestSuite) TestDeleteGroupProfilePicture() {
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
			desc:               "DeleteProfilePictureNoRights",
			userID:             s.IDs["userWithoutRights"].String(),
			groupID:            s.IDs["groupOK"].String(),
			expectedStatusCode: http.StatusForbidden,
			expectedResponse:   gin.H{"err": "Forbidden action. Reason: User ee2c6112-1114-4d9f-8869-716068ff7159 has no rights to set in group 4552667f-ea03-4ad3-8757-ea4645c8b4a0"},
		},
		{
			desc:               "DeleteProfilePictureNoPicture",
			userID:             s.IDs["userOK"].String(),
			groupID:            s.IDs["groupWithoutPicture"].String(),
			expectedStatusCode: http.StatusForbidden,
			expectedResponse:   gin.H{"err": "Forbidden action. Reason: group 4399b92e-d68a-42be-9a01-a9e098df98d8 has no profile picture"},
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
			json.NewDecoder(response.Body).Decode(&msg)

			s.Equal(tC.expectedResponse, msg)
		})
	}
}

func (s GroupPicturesTestSuite) TestSetGroupProfilePicture() {
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
			desc:               "UpdateProfilePictureNoRights",
			userID:             s.IDs["userWithoutRights"].String(),
			groupID:            s.IDs["groupOK"].String(),
			imageData:          map[string]string{"Key": "avatarFile", "CType": "image/png"},
			setBodyLimiter:     false,
			expectedStatusCode: http.StatusForbidden,
			expectedResponse:   gin.H{"err": "Forbidden action. Reason: User ee2c6112-1114-4d9f-8869-716068ff7159 has no rights to set in group 4552667f-ea03-4ad3-8757-ea4645c8b4a0"},
		},
		{
			desc:               "UpdateProfilePictureWrongImageType",
			userID:             s.IDs["userOK"].String(),
			groupID:            s.IDs["groupOK"].String(),
			imageData:          map[string]string{"Key": "avatarFile", "CType": "application/octet-stream"},
			setBodyLimiter:     false,
			expectedStatusCode: http.StatusBadRequest,
			expectedResponse:   gin.H{"err": "image extention not allowed"},
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

			var msg gin.H
			json.NewDecoder(response.Body).Decode(&msg)

			if !tC.setBodyLimiter {
				s.Equal(tC.expectedResponse, msg)
			}
		})
	}
}

func TestGroupPicturesSuite(t *testing.T) {
	suite.Run(t, &GroupPicturesTestSuite{})
}
