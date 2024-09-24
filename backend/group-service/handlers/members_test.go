package handlers_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Slimo300/Microservices-Videocall-and-Chat-App/backend/group-service/handlers"
	"github.com/Slimo300/Microservices-Videocall-and-Chat-App/backend/group-service/models"
	mockservice "github.com/Slimo300/Microservices-Videocall-and-Chat-App/backend/group-service/service/mock"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

type MembersTestSuite struct {
	suite.Suite
	IDs    map[string]uuid.UUID
	server *handlers.Server
}

func (s *MembersTestSuite) SetupSuite() {

	s.IDs = make(map[string]uuid.UUID)

	s.IDs["user"] = uuid.New()
	s.IDs["group"] = uuid.New()
	s.IDs["member"] = uuid.New()

	service := new(mockservice.GroupsMockService)
	service.On("GrantRights", s.IDs["user"], s.IDs["member"], mock.Anything).Return(&models.Member{ID: s.IDs["member"]}, nil)
	service.On("DeleteMember", s.IDs["user"], s.IDs["member"]).Return(&models.Member{ID: s.IDs["member"]}, nil)
	service.On("DeleteGroup", s.IDs["user"], s.IDs["group"]).Return(&models.Group{ID: s.IDs["group"]}, nil)

	s.server = handlers.NewServer(service, nil)
}

func (s *MembersTestSuite) TestGrantRights() {
	gin.SetMode(gin.TestMode)

	testCases := []struct {
		desc               string
		userID             string
		memberID           string
		data               map[string]interface{}
		expectedStatusCode int
		expectedResponse   interface{}
	}{
		{
			desc:               "UpdateRightsBadUserID",
			userID:             s.IDs["user"].String()[:2],
			memberID:           s.IDs["member"].String(),
			data:               map[string]interface{}{"adding": -1},
			expectedStatusCode: http.StatusBadRequest,
			expectedResponse:   gin.H{"err": "invalid ID"},
		},
		{
			desc:               "UpdateRightsBadMemberID",
			userID:             s.IDs["user"].String(),
			memberID:           s.IDs["member"].String()[:2],
			data:               map[string]interface{}{"adding": -1},
			expectedStatusCode: http.StatusBadRequest,
			expectedResponse:   gin.H{"err": "invalid member ID"},
		},
		{
			desc:               "UpdateRightsNoAction",
			userID:             s.IDs["user"].String(),
			memberID:           s.IDs["member"].String(),
			data:               map[string]interface{}{},
			expectedStatusCode: http.StatusBadRequest,
			expectedResponse:   gin.H{"err": "no action specified"},
		},
		{
			desc:               "UpdateRightsSuccess",
			userID:             s.IDs["user"].String(),
			memberID:           s.IDs["member"].String(),
			data:               map[string]interface{}{"adding": -1},
			expectedStatusCode: http.StatusOK,
			expectedResponse:   gin.H{"message": "member updated"},
		},
	}

	for _, tC := range testCases {
		s.Run(tC.desc, func() {

			requestBody, _ := json.Marshal(tC.data)
			req, _ := http.NewRequest(http.MethodPut, "/group/members/"+tC.memberID, bytes.NewBuffer(requestBody))

			w := httptest.NewRecorder()
			_, engine := gin.CreateTestContext(w)
			engine.Use(func(c *gin.Context) {
				c.Set("userID", tC.userID)
			})

			engine.Handle(http.MethodPut, "/group/members/:memberID", s.server.GrantRights)
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

func (s *MembersTestSuite) TestDeleteMember() {
	gin.SetMode(gin.TestMode)

	testCases := []struct {
		desc               string
		userID             string
		memberID           string
		expectedStatusCode int
		expectedResponse   interface{}
	}{
		{
			desc:               "DeleteMemberBadUserID",
			userID:             s.IDs["user"].String()[:2],
			memberID:           s.IDs["member"].String(),
			expectedStatusCode: http.StatusBadRequest,
			expectedResponse:   gin.H{"err": "invalid ID"},
		},
		{
			desc:               "DeleteMemberBadMemberID",
			userID:             s.IDs["user"].String(),
			memberID:           s.IDs["member"].String()[:2],
			expectedStatusCode: http.StatusBadRequest,
			expectedResponse:   gin.H{"err": "invalid member ID"},
		},
		{
			desc:               "DeleteMemberSuccess",
			userID:             s.IDs["user"].String(),
			memberID:           s.IDs["member"].String(),
			expectedStatusCode: http.StatusOK,
			expectedResponse:   gin.H{"message": "member deleted"},
		},
	}

	for _, tC := range testCases {
		s.Run(tC.desc, func() {
			req, _ := http.NewRequest(http.MethodDelete, "/group/members/"+tC.memberID, nil)

			w := httptest.NewRecorder()
			_, engine := gin.CreateTestContext(w)
			engine.Use(func(c *gin.Context) {
				c.Set("userID", tC.userID)
			})

			engine.Handle(http.MethodDelete, "/group/members/:memberID", s.server.DeleteUserFromGroup)
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

func (s *MembersTestSuite) TestDeleteGroup() {
	gin.SetMode(gin.TestMode)

	testCases := []struct {
		desc               string
		userID             string
		groupID            string
		expectedStatusCode int
		expectedResponse   interface{}
	}{

		{
			desc:               "DeleteGroupBadUserID",
			userID:             s.IDs["user"].String()[:2],
			groupID:            s.IDs["group"].String(),
			expectedStatusCode: http.StatusBadRequest,
			expectedResponse:   gin.H{"err": "invalid ID"},
		},
		{
			desc:               "DeleteGroupBadGroupID",
			userID:             s.IDs["user"].String(),
			groupID:            s.IDs["group"].String()[:2],
			expectedStatusCode: http.StatusBadRequest,
			expectedResponse:   gin.H{"err": "invalid group ID"},
		},
		{
			desc:               "DeleteGroupSuccess",
			userID:             s.IDs["user"].String(),
			groupID:            s.IDs["group"].String(),
			expectedStatusCode: http.StatusOK,
			expectedResponse:   gin.H{"message": "group deleted"},
		},
	}

	for _, tC := range testCases {
		s.Run(tC.desc, func() {
			req, _ := http.NewRequest(http.MethodDelete, "/api/group/"+tC.groupID, nil)

			w := httptest.NewRecorder()
			_, engine := gin.CreateTestContext(w)

			engine.Use(func(c *gin.Context) {
				c.Set("userID", tC.userID)
			})
			engine.Handle(http.MethodDelete, "/api/group/:groupID", s.server.DeleteGroup)
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

func TestMembers(t *testing.T) {
	suite.Run(t, &MembersTestSuite{})
}
