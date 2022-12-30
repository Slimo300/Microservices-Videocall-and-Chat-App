package handlers_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	dbmock "github.com/Slimo300/MicroservicesChatApp/backend/group-service/database/mock"
	"github.com/Slimo300/MicroservicesChatApp/backend/group-service/handlers"
	"github.com/Slimo300/MicroservicesChatApp/backend/group-service/models"
	"github.com/Slimo300/MicroservicesChatApp/backend/lib/apperrors"
	mockqueue "github.com/Slimo300/MicroservicesChatApp/backend/lib/msgqueue/mock"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

type InvitesTestSuite struct {
	suite.Suite
	IDs    map[string]uuid.UUID
	server *handlers.Server
}

func (s *InvitesTestSuite) SetupSuite() {

	s.IDs = make(map[string]uuid.UUID)

	s.IDs["inviteOK"] = uuid.MustParse("9248e828-8120-4f6d-a2c5-25a4689b9ba8")
	s.IDs["inviteNotFound"] = uuid.MustParse("2917d4d0-b3ed-49ff-93de-d5913d24a6c8")
	s.IDs["inviteAnswered"] = uuid.MustParse("a901767d-d908-471d-8a9a-f01945547da9")
	s.IDs["userOK"] = uuid.MustParse("f515cb74-99b2-4aa9-be0d-faf1a68c8064")
	s.IDs["userWithoutInvites"] = uuid.MustParse("1414bb70-a865-4a88-8c5d-adbe7fa1ec53")
	s.IDs["userNoRights"] = uuid.MustParse("58bb1c85-7f6a-4e2b-90a9-b974928a81c4")
	s.IDs["invitedUserOK"] = uuid.MustParse("34b80593-57fc-4953-8cd3-217ecb61b6eb")
	s.IDs["invitedUserNotFound"] = uuid.MustParse("6ebb22de-1bd6-4c23-bb0f-eec359d10462")
	s.IDs["invitedUserMember"] = uuid.MustParse("27df64da-a103-49fb-9724-151cdb2943b5")
	s.IDs["invitedUserInvited"] = uuid.MustParse("34234be4-fe92-49cb-9ddd-76ba9f410266")
	s.IDs["group"] = uuid.MustParse("b646e70f-3c8f-4782-84a3-0b34b0f9aecf")

	db := new(dbmock.MockGroupsDB)
	db.On("GetUserInvites", s.IDs["userOK"], 1, 0).Return([]models.Invite{{ID: s.IDs["inviteOK"]}}, nil)
	db.On("GetUserInvites", s.IDs["userWithoutInvites"], 1, 0).Return([]models.Invite{}, nil)

	db.On("AddInvite", s.IDs["userNoRights"], s.IDs["invitedUserOK"], s.IDs["group"]).
		Return(&models.Invite{}, apperrors.NewForbidden(fmt.Sprintf("User %v has no rights to add new members to group %v", s.IDs["userNoRights"], s.IDs["group"])))
	db.On("AddInvite", s.IDs["userOK"], s.IDs["invitedUserNotFound"], s.IDs["group"]).
		Return(&models.Invite{}, apperrors.NewNotFound("user", s.IDs["invitedUserNotFound"].String()))
	db.On("AddInvite", s.IDs["userOK"], s.IDs["invitedUserMember"], s.IDs["group"]).
		Return(&models.Invite{}, apperrors.NewForbidden(fmt.Sprintf("User %v already is already a member of group %v", s.IDs["invitedUserMember"], s.IDs["group"])))
	db.On("AddInvite", s.IDs["userOK"], s.IDs["invitedUserInvited"], s.IDs["group"]).
		Return(&models.Invite{}, apperrors.NewForbidden(fmt.Sprintf("User %v already invited to group %v", s.IDs["invitedUserInvited"], s.IDs["group"])))
	db.On("AddInvite", s.IDs["userOK"], s.IDs["invitedUserOK"], s.IDs["group"]).
		Return(&models.Invite{ID: s.IDs["inviteOK"]}, nil)

	db.On("AnswerInvite", s.IDs["userOK"], s.IDs["inviteOK"], true).Return(&models.Invite{ID: s.IDs["inviteOK"]}, &models.Group{ID: s.IDs["group"]}, nil, nil)
	db.On("AnswerInvite", s.IDs["userOK"], s.IDs["inviteOK"], false).Return(&models.Invite{ID: s.IDs["inviteOK"]}, nil, nil, nil)
	db.On("AnswerInvite", s.IDs["userOK"], s.IDs["inviteNotFound"], mock.Anything).
		Return(nil, nil, nil, apperrors.NewNotFound("invite", s.IDs["inviteNotFound"].String()))
	db.On("AnswerInvite", s.IDs["userOK"], s.IDs["inviteAnswered"], mock.Anything).
		Return(nil, nil, nil, apperrors.NewForbidden("invite already answered"))

	s.server = handlers.NewServer(db, nil, nil)

	emitter := new(mockqueue.MockEmitter)
	emitter.On("Emit", mock.Anything).Return(nil)
	s.server.Emitter = emitter

}

func (s InvitesTestSuite) TestGetUserInvites() {
	gin.SetMode(gin.TestMode)

	testCases := []struct {
		desc               string
		id                 string
		expectedStatusCode int
		expectedResponse   []models.Invite
	}{
		{
			desc:               "getinvitessuccess",
			id:                 s.IDs["userOK"].String(),
			expectedStatusCode: http.StatusOK,
			expectedResponse:   []models.Invite{{ID: s.IDs["inviteOK"]}},
		},
		{
			desc:               "getinvitesnocontent",
			id:                 s.IDs["userWithoutInvites"].String(),
			expectedStatusCode: http.StatusNoContent,
			expectedResponse:   []models.Invite{},
		},
	}

	for _, tC := range testCases {
		s.Run(tC.desc, func() {

			req, _ := http.NewRequest("GET", "/api/invites?num=1&offset=0", nil)

			w := httptest.NewRecorder()
			_, engine := gin.CreateTestContext(w)
			engine.Use(func(c *gin.Context) {
				c.Set("userID", tC.id)
			})

			engine.Handle(http.MethodGet, "/api/invites", s.server.GetUserInvites)
			engine.ServeHTTP(w, req)
			response := w.Result()
			defer response.Body.Close()

			s.Equal(tC.expectedStatusCode, response.StatusCode)

			respBody := []models.Invite{}
			json.NewDecoder(response.Body).Decode(&respBody)

			s.Equal(tC.expectedResponse, respBody)
		})
	}
}

func (s InvitesTestSuite) TestSendGroupInvite() {
	gin.SetMode(gin.TestMode)

	testCases := []struct {
		desc               string
		id                 string
		data               map[string]interface{}
		returnVal          bool
		expectedStatusCode int
		expectedResponse   interface{}
	}{
		{
			desc:               "inviteNoGroup",
			id:                 s.IDs["userOK"].String(),
			data:               map[string]interface{}{"target": s.IDs["invitedUserOK"].String()},
			returnVal:          false,
			expectedStatusCode: http.StatusBadRequest,
			expectedResponse:   gin.H{"err": "invalid group ID"},
		},
		{
			desc:               "inviteNoUser",
			id:                 s.IDs["userOK"].String(),
			data:               map[string]interface{}{"group": s.IDs["group"].String()},
			returnVal:          false,
			expectedStatusCode: http.StatusBadRequest,
			expectedResponse:   gin.H{"err": "invalid target user ID"},
		},
		{
			desc:               "inviteNoRights",
			id:                 s.IDs["userNoRights"].String(),
			data:               map[string]interface{}{"group": s.IDs["group"].String(), "target": s.IDs["invitedUserOK"]},
			returnVal:          false,
			expectedStatusCode: http.StatusForbidden,
			expectedResponse:   gin.H{"err": fmt.Sprintf("Forbidden action. Reason: User %v has no rights to add new members to group %v", s.IDs["userNoRights"], s.IDs["group"])},
		},
		{
			desc:               "inviteUserNotFound",
			id:                 s.IDs["userOK"].String(),
			data:               map[string]interface{}{"group": s.IDs["group"].String(), "target": s.IDs["invitedUserNotFound"]},
			returnVal:          false,
			expectedStatusCode: http.StatusNotFound,
			expectedResponse:   gin.H{"err": fmt.Sprintf("resource: user with value: %v not found", s.IDs["invitedUserNotFound"])},
		},
		{
			desc:               "inviteUserMember",
			id:                 s.IDs["userOK"].String(),
			data:               map[string]interface{}{"group": s.IDs["group"].String(), "target": s.IDs["invitedUserMember"]},
			returnVal:          false,
			expectedStatusCode: http.StatusForbidden,
			expectedResponse:   gin.H{"err": fmt.Sprintf("Forbidden action. Reason: User %v already is already a member of group %v", s.IDs["invitedUserMember"], s.IDs["group"])},
		},
		{
			desc:               "inviteUserInvited",
			id:                 s.IDs["userOK"].String(),
			data:               map[string]interface{}{"group": s.IDs["group"].String(), "target": s.IDs["invitedUserInvited"]},
			returnVal:          false,
			expectedStatusCode: http.StatusForbidden,
			expectedResponse:   gin.H{"err": fmt.Sprintf("Forbidden action. Reason: User %v already invited to group %v", s.IDs["invitedUserInvited"], s.IDs["group"])},
		},
		{
			desc:               "invitesuccess",
			id:                 s.IDs["userOK"].String(),
			data:               map[string]interface{}{"group": s.IDs["group"].String(), "target": s.IDs["invitedUserOK"].String()},
			returnVal:          true,
			expectedStatusCode: http.StatusCreated,
			expectedResponse:   models.Invite{ID: s.IDs["inviteOK"]},
		},
	}

	for _, tC := range testCases {
		s.Run(tC.desc, func() {

			requestBody, _ := json.Marshal(tC.data)
			req, _ := http.NewRequest("POST", "/api/invite", bytes.NewReader(requestBody))

			w := httptest.NewRecorder()
			_, engine := gin.CreateTestContext(w)
			engine.Use(func(c *gin.Context) {
				c.Set("userID", tC.id)
			})

			engine.Handle(http.MethodPost, "/api/invite", s.server.CreateInvite)
			engine.ServeHTTP(w, req)
			response := w.Result()
			defer response.Body.Close()

			s.Equal(tC.expectedStatusCode, response.StatusCode)

			var respBody interface{}

			if tC.returnVal {
				var invite models.Invite
				json.NewDecoder(response.Body).Decode(&invite)
				respBody = invite
			} else {
				var msg gin.H
				json.NewDecoder(response.Body).Decode(&msg)
				respBody = msg
			}

			s.Equal(tC.expectedResponse, respBody)
		})
	}
}

func (s InvitesTestSuite) TestRespondGroupInvite() {
	gin.SetMode(gin.TestMode)

	testCases := []struct {
		desc               string
		userID             string
		inviteID           string
		data               map[string]interface{}
		returnVal          bool
		expectedStatusCode int
		expectedResponse   interface{}
	}{
		{
			desc:               "respondInviteInvalidUserID",
			userID:             s.IDs["userOK"].String()[:2],
			inviteID:           s.IDs["inviteOK"].String(),
			data:               map[string]interface{}{"answer": true},
			returnVal:          false,
			expectedStatusCode: http.StatusBadRequest,
			expectedResponse:   gin.H{"err": "invalid ID"},
		},
		{
			desc:               "respondInviteInvalidInviteID",
			userID:             s.IDs["userOK"].String(),
			inviteID:           s.IDs["inviteOK"].String()[:2],
			data:               map[string]interface{}{"answer": true},
			returnVal:          false,
			expectedStatusCode: http.StatusBadRequest,
			expectedResponse:   gin.H{"err": "invalid invite id"},
		},
		{
			desc:               "respondInviteNoAnswer",
			userID:             s.IDs["userOK"].String(),
			inviteID:           s.IDs["inviteOK"].String(),
			data:               map[string]interface{}{},
			returnVal:          false,
			expectedStatusCode: http.StatusBadRequest,
			expectedResponse:   gin.H{"err": "answer not specified"},
		},
		{
			desc:               "respondInviteNotFound",
			userID:             s.IDs["userOK"].String(),
			inviteID:           s.IDs["inviteNotFound"].String(),
			data:               map[string]interface{}{"answer": true},
			returnVal:          false,
			expectedStatusCode: http.StatusNotFound,
			expectedResponse:   gin.H{"err": "resource: invite with value: 2917d4d0-b3ed-49ff-93de-d5913d24a6c8 not found"},
		},
		{
			desc:               "respondInviteAnswered",
			userID:             s.IDs["userOK"].String(),
			inviteID:           s.IDs["inviteAnswered"].String(),
			data:               map[string]interface{}{"answer": true},
			returnVal:          false,
			expectedStatusCode: http.StatusForbidden,
			expectedResponse:   gin.H{"err": "Forbidden action. Reason: invite already answered"},
		},
		{
			desc:               "respondInviteNo",
			userID:             s.IDs["userOK"].String(),
			inviteID:           s.IDs["inviteOK"].String(),
			data:               map[string]interface{}{"answer": false},
			returnVal:          true,
			expectedStatusCode: http.StatusOK,
			expectedResponse:   gin.H{"invite": models.Invite{ID: s.IDs["inviteOK"]}},
		},
		{
			desc:               "respondInviteYes",
			userID:             s.IDs["userOK"].String(),
			inviteID:           s.IDs["inviteOK"].String(),
			data:               map[string]interface{}{"answer": true},
			returnVal:          true,
			expectedStatusCode: http.StatusOK,
			expectedResponse:   gin.H{"invite": models.Invite{ID: s.IDs["inviteOK"]}, "group": models.Group{ID: s.IDs["group"]}},
		},
	}

	for _, tC := range testCases {
		s.Run(tC.desc, func() {

			requestBody, _ := json.Marshal(tC.data)
			req, _ := http.NewRequest("PUT", "/api/invite/"+tC.inviteID, bytes.NewReader(requestBody))

			w := httptest.NewRecorder()
			_, engine := gin.CreateTestContext(w)
			engine.Use(func(c *gin.Context) {
				c.Set("userID", tC.userID)
			})

			engine.Handle(http.MethodPut, "/api/invite/:inviteID", s.server.RespondGroupInvite)
			engine.ServeHTTP(w, req)
			response := w.Result()
			defer response.Body.Close()

			s.Equal(tC.expectedStatusCode, response.StatusCode)

			var respBody interface{}
			if tC.returnVal {
				var msg gin.H
				json.NewDecoder(response.Body).Decode(&msg)

				inviteAsInterface, ok := msg["invite"]
				if !ok {
					s.Fail("Returned message does not contain invite field")
				}
				inviteAsMap, ok := inviteAsInterface.(map[string]interface{})
				if !ok {
					s.Fail("group is not a map")
				}
				inviteID, err := uuid.Parse(inviteAsMap["ID"].(string))
				if err != nil {
					s.Fail("Parsing invite ID returned err: ", err.Error())
				}

				var groupID uuid.UUID
				groupAsInterface, ok := msg["group"]
				if ok {

					groupAsMap, ok := groupAsInterface.(map[string]interface{})
					if !ok {
						s.Fail("group is not a map")
					}
					groupID, err = uuid.Parse(groupAsMap["ID"].(string))
					if err != nil {
						s.Fail("Parsing invite ID returned err: ", err.Error())
					}
					respBody = gin.H{"invite": models.Invite{ID: inviteID}, "group": models.Group{ID: groupID}}
				} else {
					respBody = gin.H{"invite": models.Invite{ID: inviteID}}
				}

			} else {
				var msg gin.H
				json.NewDecoder(response.Body).Decode(&msg)
				respBody = msg
			}

			s.Equal(tC.expectedResponse, respBody)
		})
	}
}

func TestInvitesSuite(t *testing.T) {
	suite.Run(t, &InvitesTestSuite{})
}
