package handlers_test

// import (
// 	"bytes"
// 	"encoding/json"
// 	"net/http"
// 	"net/http/httptest"
// 	"testing"

// 	"github.com/Slimo300/Microservices-Videocall-and-Chat-App/backend/group-service/handlers"
// 	"github.com/Slimo300/Microservices-Videocall-and-Chat-App/backend/group-service/models"
// 	mockservice "github.com/Slimo300/Microservices-Videocall-and-Chat-App/backend/group-service/service/mock"
// 	"github.com/gin-gonic/gin"
// 	"github.com/google/uuid"
// 	"github.com/stretchr/testify/mock"
// 	"github.com/stretchr/testify/suite"
// )

// type GroupTestSuite struct {
// 	suite.Suite
// 	IDs    map[string]uuid.UUID
// 	server *handlers.Server
// }

// func (s *GroupTestSuite) SetupSuite() {

// 	s.IDs = make(map[string]uuid.UUID)

// 	s.IDs["user1"] = uuid.MustParse("d562da33-062b-4f18-afa0-f9fd1b7aadd3")
// 	s.IDs["user2"] = uuid.MustParse("b86dcd49-2fb3-48b1-bee7-ee46d4682244")
// 	s.IDs["group1"] = uuid.MustParse("e7564f8c-8917-4527-a020-1db4b901d4b9")
// 	s.IDs["group2"] = uuid.MustParse("be0accb1-50db-4698-b048-fb0128e35684")

// 	s.IDs["member"] = uuid.MustParse("6c564875-cd55-4e20-a035-44f1750d25b9")

// 	service := new(mockservice.GroupsMockService)
// 	service.On("GetUserGroups", mock.Anything, s.IDs["user1"]).Return([]*models.Group{
// 		{ID: s.IDs["group1"]},
// 		{ID: s.IDs["group2"]},
// 	}, nil)
// 	service.On("GetUserGroups", mock.Anything, s.IDs["user2"]).Return([]*models.Group{}, nil)

// 	service.On("CreateGroup", mock.Anything, mock.Anything, mock.Anything).Return(&models.Group{}, nil)

// 	service.On("DeleteGroup", mock.Anything, mock.Anything).Return(&models.Group{}, nil)

// 	s.server = handlers.NewServer(service, nil)
// }

// func (s *GroupTestSuite) TestGetUserGroups() {
// 	gin.SetMode(gin.TestMode)

// 	testCases := []struct {
// 		desc               string
// 		userID             string
// 		returnVal          bool
// 		expectedStatusCode int
// 		expectedResponse   interface{}
// 	}{
// 		{
// 			desc:               "GetGroupsSuccess",
// 			userID:             s.IDs["user1"].String(),
// 			returnVal:          true,
// 			expectedStatusCode: http.StatusOK,
// 			expectedResponse: []models.Group{
// 				{ID: s.IDs["group1"]},
// 				{ID: s.IDs["group2"]},
// 			},
// 		},
// 		{
// 			desc:               "GetGroupsNone",
// 			userID:             s.IDs["user2"].String(),
// 			returnVal:          false,
// 			expectedStatusCode: http.StatusNoContent,
// 			expectedResponse:   nil,
// 		},
// 	}

// 	for _, tC := range testCases {
// 		s.Run(tC.desc, func() {

// 			req, _ := http.NewRequest("GET", "/api/group/get", nil)

// 			w := httptest.NewRecorder()
// 			_, engine := gin.CreateTestContext(w)

// 			engine.Use(func(c *gin.Context) {
// 				c.Set("userID", tC.userID)
// 			})
// 			engine.Handle(http.MethodGet, "/api/group/get", s.server.GetUserGroups)
// 			engine.ServeHTTP(w, req)
// 			response := w.Result()
// 			defer response.Body.Close()

// 			s.Equal(response.StatusCode, tC.expectedStatusCode)

// 			if tC.returnVal {
// 				respBody := []models.Group{}
// 				if err := json.NewDecoder(response.Body).Decode(&respBody); err != nil {
// 					s.Fail(err.Error())
// 				}

// 				s.Equal(respBody, tC.expectedResponse)
// 			}
// 		})
// 	}
// }

// func (s *GroupTestSuite) TestCreateGroup() {
// 	gin.SetMode(gin.TestMode)

// 	testCases := []struct {
// 		desc               string
// 		userID             string
// 		data               map[string]interface{}
// 		returnVal          bool
// 		expectedStatusCode int
// 		expectedResponse   interface{}
// 	}{
// 		{
// 			desc:               "CreateGroupBadUserID",
// 			userID:             s.IDs["user1"].String()[:2],
// 			data:               map[string]interface{}{"name": "New Group"},
// 			returnVal:          false,
// 			expectedStatusCode: http.StatusBadRequest,
// 			expectedResponse:   gin.H{"err": "invalid ID"},
// 		},
// 		{
// 			desc:               "CreateGroupNoName",
// 			userID:             s.IDs["user1"].String(),
// 			data:               map[string]interface{}{"name": ""},
// 			returnVal:          false,
// 			expectedStatusCode: http.StatusBadRequest,
// 			expectedResponse:   gin.H{"err": "name not specified"},
// 		},
// 		{
// 			desc:               "CreateGroupSuccess",
// 			userID:             s.IDs["user1"].String(),
// 			data:               map[string]interface{}{"name": "New Group"},
// 			returnVal:          true,
// 			expectedStatusCode: http.StatusCreated,
// 			expectedResponse:   models.Group{},
// 		},
// 	}

// 	for _, tC := range testCases {
// 		s.Run(tC.desc, func() {

// 			requestBody, _ := json.Marshal(tC.data)
// 			req, _ := http.NewRequest("POST", "/api/group/create", bytes.NewBuffer(requestBody))

// 			w := httptest.NewRecorder()
// 			_, engine := gin.CreateTestContext(w)

// 			engine.Use(func(c *gin.Context) {
// 				c.Set("userID", tC.userID)
// 			})

// 			engine.Handle(http.MethodPost, "/api/group/create", s.server.CreateGroup)
// 			engine.ServeHTTP(w, req)
// 			response := w.Result()
// 			defer response.Body.Close()

// 			s.Equal(tC.expectedStatusCode, response.StatusCode)

// 			var respBody interface{}
// 			if tC.returnVal {
// 				group := models.Group{}
// 				if err := json.NewDecoder(response.Body).Decode(&group); err != nil {
// 					s.Fail(err.Error())
// 				}
// 				respBody = group
// 			} else {
// 				var msg gin.H
// 				if err := json.NewDecoder(response.Body).Decode(&msg); err != nil {
// 					s.Fail(err.Error())
// 				}
// 				respBody = msg
// 			}

// 			s.Equal(tC.expectedResponse, respBody)
// 		})
// 	}
// }

// func TestGroupSuite(t *testing.T) {
// 	suite.Run(t, &GroupTestSuite{})
// }
