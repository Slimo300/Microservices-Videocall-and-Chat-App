package handlers_test

// import (
// 	"bytes"
// 	"encoding/json"
// 	"io"
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

// type InvitesTestSuite struct {
// 	suite.Suite
// 	IDs    map[string]uuid.UUID
// 	server *handlers.Server
// }

// func (s *InvitesTestSuite) SetupSuite() {

// 	s.IDs = make(map[string]uuid.UUID)

// 	s.IDs["invite"] = uuid.MustParse("9248e828-8120-4f6d-a2c5-25a4689b9ba8")
// 	s.IDs["user"] = uuid.MustParse("f515cb74-99b2-4aa9-be0d-faf1a68c8064")
// 	s.IDs["target"] = uuid.New()
// 	s.IDs["userWithoutInvites"] = uuid.MustParse("1414bb70-a865-4a88-8c5d-adbe7fa1ec53")
// 	s.IDs["group"] = uuid.MustParse("b646e70f-3c8f-4782-84a3-0b34b0f9aecf")

// 	service := new(mockservice.GroupsMockService)
// 	service.On("GetUserInvites", mock.Anything, s.IDs["user"], 2, 0).Return([]*models.Invite{{ID: s.IDs["invite"]}, {ID: s.IDs["inviteAnswered"]}}, nil)
// 	service.On("GetUserInvites", mock.Anything, s.IDs["userWithoutInvites"], 2, 0).Return([]*models.Invite{}, nil)

// 	service.On("AddInvite", mock.Anything, s.IDs["user"], s.IDs["target"], s.IDs["group"]).Return(&models.Invite{ID: s.IDs["invite"]}, nil)

// 	service.On("RespondInvite", mock.Anything, s.IDs["user"], s.IDs["invite"], false).Return(&models.Invite{ID: s.IDs["invite"]}, nil, nil)
// 	service.On("RespondInvite", mock.Anything, s.IDs["user"], s.IDs["invite"], true).Return(&models.Invite{ID: s.IDs["invite"]}, &models.Group{ID: s.IDs["group"]}, nil)

// 	s.server = handlers.NewServer(service, nil)
// }

// func (s *InvitesTestSuite) TestGetUserInvites() {
// 	gin.SetMode(gin.TestMode)

// 	testCases := []struct {
// 		desc               string
// 		id                 string
// 		expectedStatusCode int
// 		expectedResponse   []models.Invite
// 	}{
// 		{
// 			desc:               "getinvitessuccess",
// 			id:                 s.IDs["user"].String(),
// 			expectedStatusCode: http.StatusOK,
// 			expectedResponse:   []models.Invite{{ID: s.IDs["invite"]}, {ID: s.IDs["inviteAnswered"]}},
// 		},
// 		{
// 			desc:               "getinvitesnocontent",
// 			id:                 s.IDs["userWithoutInvites"].String(),
// 			expectedStatusCode: http.StatusNoContent,
// 			expectedResponse:   []models.Invite{},
// 		},
// 	}

// 	for _, tC := range testCases {
// 		s.Run(tC.desc, func() {

// 			req, _ := http.NewRequest("GET", "/api/invites?num=2&offset=0", nil)

// 			w := httptest.NewRecorder()
// 			_, engine := gin.CreateTestContext(w)
// 			engine.Use(func(c *gin.Context) {
// 				c.Set("userID", tC.id)
// 			})

// 			engine.Handle(http.MethodGet, "/api/invites", s.server.GetUserInvites)
// 			engine.ServeHTTP(w, req)
// 			response := w.Result()
// 			defer response.Body.Close()

// 			s.Equal(tC.expectedStatusCode, response.StatusCode)

// 			respBody := []models.Invite{}
// 			if err := json.NewDecoder(response.Body).Decode(&respBody); err != nil && err != io.EOF {
// 				s.Fail(err.Error())
// 			}

// 			s.Equal(tC.expectedResponse, respBody)
// 		})
// 	}
// }

// func (s *InvitesTestSuite) TestSendGroupInvite() {
// 	gin.SetMode(gin.TestMode)

// 	testCases := []struct {
// 		desc               string
// 		id                 string
// 		data               map[string]interface{}
// 		returnVal          bool
// 		expectedStatusCode int
// 		expectedResponse   interface{}
// 	}{
// 		{
// 			desc:               "inviteNoGroup",
// 			id:                 s.IDs["user"].String(),
// 			data:               map[string]interface{}{"target": s.IDs["invitedUser"].String()},
// 			returnVal:          false,
// 			expectedStatusCode: http.StatusBadRequest,
// 			expectedResponse:   gin.H{"err": "invalid group ID"},
// 		},
// 		{
// 			desc:               "inviteNoUser",
// 			id:                 s.IDs["user"].String(),
// 			data:               map[string]interface{}{"group": s.IDs["group"].String()},
// 			returnVal:          false,
// 			expectedStatusCode: http.StatusBadRequest,
// 			expectedResponse:   gin.H{"err": "invalid target user ID"},
// 		},
// 		{
// 			desc:               "invitesuccess",
// 			id:                 s.IDs["user"].String(),
// 			data:               map[string]interface{}{"group": s.IDs["group"].String(), "target": s.IDs["target"].String()},
// 			returnVal:          true,
// 			expectedStatusCode: http.StatusCreated,
// 			expectedResponse:   models.Invite{ID: s.IDs["invite"]},
// 		},
// 	}

// 	for _, tC := range testCases {
// 		s.Run(tC.desc, func() {

// 			requestBody, _ := json.Marshal(tC.data)
// 			req, _ := http.NewRequest("POST", "/api/invite", bytes.NewReader(requestBody))

// 			w := httptest.NewRecorder()
// 			_, engine := gin.CreateTestContext(w)
// 			engine.Use(func(c *gin.Context) {
// 				c.Set("userID", tC.id)
// 			})

// 			engine.Handle(http.MethodPost, "/api/invite", s.server.CreateInvite)
// 			engine.ServeHTTP(w, req)
// 			response := w.Result()
// 			defer response.Body.Close()

// 			s.Equal(tC.expectedStatusCode, response.StatusCode)

// 			var respBody interface{}

// 			if tC.returnVal {
// 				var invite models.Invite
// 				if err := json.NewDecoder(response.Body).Decode(&invite); err != nil {
// 					s.Fail(err.Error())
// 				}
// 				respBody = invite
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

// func (s *InvitesTestSuite) TestRespondGroupInvite() {
// 	gin.SetMode(gin.TestMode)

// 	testCases := []struct {
// 		desc               string
// 		userID             string
// 		inviteID           string
// 		data               map[string]interface{}
// 		returnVal          bool
// 		expectedStatusCode int
// 		expectedResponse   interface{}
// 	}{
// 		{
// 			desc:               "respondInviteInvalidUserID",
// 			userID:             s.IDs["user"].String()[:2],
// 			inviteID:           s.IDs["invite"].String(),
// 			data:               map[string]interface{}{"answer": true},
// 			returnVal:          false,
// 			expectedStatusCode: http.StatusBadRequest,
// 			expectedResponse:   gin.H{"err": "invalid ID"},
// 		},
// 		{
// 			desc:               "respondInviteInvalidInviteID",
// 			userID:             s.IDs["user"].String(),
// 			inviteID:           s.IDs["invite"].String()[:2],
// 			data:               map[string]interface{}{"answer": true},
// 			returnVal:          false,
// 			expectedStatusCode: http.StatusBadRequest,
// 			expectedResponse:   gin.H{"err": "invalid invite id"},
// 		},
// 		{
// 			desc:               "respondInviteNoAnswer",
// 			userID:             s.IDs["user"].String(),
// 			inviteID:           s.IDs["invite"].String(),
// 			data:               map[string]interface{}{},
// 			returnVal:          false,
// 			expectedStatusCode: http.StatusBadRequest,
// 			expectedResponse:   gin.H{"err": "answer not specified"},
// 		},
// 		{
// 			desc:               "respondInviteNo",
// 			userID:             s.IDs["user"].String(),
// 			inviteID:           s.IDs["invite"].String(),
// 			data:               map[string]interface{}{"answer": false},
// 			returnVal:          true,
// 			expectedStatusCode: http.StatusOK,
// 			expectedResponse:   gin.H{"invite": models.Invite{ID: s.IDs["invite"]}},
// 		},
// 		{
// 			desc:               "respondInviteYes",
// 			userID:             s.IDs["user"].String(),
// 			inviteID:           s.IDs["invite"].String(),
// 			data:               map[string]interface{}{"answer": true},
// 			returnVal:          true,
// 			expectedStatusCode: http.StatusOK,
// 			expectedResponse:   gin.H{"invite": models.Invite{ID: s.IDs["invite"]}, "group": models.Group{ID: s.IDs["group"]}},
// 		},
// 	}

// 	for _, tC := range testCases {
// 		s.Run(tC.desc, func() {

// 			requestBody, _ := json.Marshal(tC.data)
// 			req, _ := http.NewRequest("PUT", "/api/invite/"+tC.inviteID, bytes.NewReader(requestBody))

// 			w := httptest.NewRecorder()
// 			_, engine := gin.CreateTestContext(w)
// 			engine.Use(func(c *gin.Context) {
// 				c.Set("userID", tC.userID)
// 			})

// 			engine.Handle(http.MethodPut, "/api/invite/:inviteID", s.server.RespondGroupInvite)
// 			engine.ServeHTTP(w, req)
// 			response := w.Result()
// 			defer response.Body.Close()

// 			s.Equal(tC.expectedStatusCode, response.StatusCode)

// 			var respBody interface{}
// 			if tC.returnVal {
// 				var msg gin.H
// 				if err := json.NewDecoder(response.Body).Decode(&msg); err != nil {
// 					s.Fail(err.Error())
// 				}

// 				inviteAsInterface, ok := msg["invite"]
// 				if !ok {
// 					s.Fail("Returned message does not contain invite field")
// 				}
// 				inviteAsMap, ok := inviteAsInterface.(map[string]interface{})
// 				if !ok {
// 					s.Fail("group is not a map")
// 				}
// 				inviteID, err := uuid.Parse(inviteAsMap["ID"].(string))
// 				if err != nil {
// 					s.Fail("Parsing invite ID returned err: ", err.Error())
// 				}

// 				var groupID uuid.UUID
// 				groupAsInterface, ok := msg["group"]
// 				if ok {

// 					groupAsMap, ok := groupAsInterface.(map[string]interface{})
// 					if !ok {
// 						s.Fail("group is not a map")
// 					}
// 					groupID, err = uuid.Parse(groupAsMap["ID"].(string))
// 					if err != nil {
// 						s.Fail("Parsing invite ID returned err: ", err.Error())
// 					}
// 					respBody = gin.H{"invite": models.Invite{ID: inviteID}, "group": models.Group{ID: groupID}}
// 				} else {
// 					respBody = gin.H{"invite": models.Invite{ID: inviteID}}
// 				}

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

// func TestInvitesSuite(t *testing.T) {
// 	suite.Run(t, &InvitesTestSuite{})
// }
