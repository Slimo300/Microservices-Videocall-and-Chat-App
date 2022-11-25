package handlers_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/Slimo300/MicroservicesChatApp/backend/lib/apperrors"
	mockqueue "github.com/Slimo300/MicroservicesChatApp/backend/lib/msgqueue/mock"
	mockdb "github.com/Slimo300/MicroservicesChatApp/backend/message-service/database/mock"
	"github.com/Slimo300/MicroservicesChatApp/backend/message-service/handlers"
	"github.com/Slimo300/MicroservicesChatApp/backend/message-service/models"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type MessageTestSuite struct {
	suite.Suite
	uuids map[string]uuid.UUID
	db    *mockdb.MockMessageDB
}

func (s *MessageTestSuite) SetupSuite() {
	// Modelled database:
	// 3 users, 2 groups, 3 messages
	// "userInGroup" belongs to both groups, "userNotInGroup" to none,
	// "userWithNoRights belongs to "group" but can't delete messages
	// "otherGroup" has no messages and is there to make sure user won't delete message if he doesn't
	// know to what group it belongs
	// 2 messages (message, deletedMessage) belong to "group" and "userInGroup"
	// "notExistingMessage" is for checking service reaction to sending unknown ID

	s.uuids = make(map[string]uuid.UUID)

	s.uuids["userInGroup"] = uuid.MustParse("1c4dccaf-a341-4920-9003-f24e0412f8e0")
	s.uuids["userNotInGroup"] = uuid.MustParse("634240cf-1219-4be2-adfa-90ab6b47899b")
	s.uuids["userWithNoRights"] = uuid.MustParse("889858f7-4b2a-4914-ad14-471a7300c26b")

	s.uuids["group"] = uuid.MustParse("61fbd273-b941-471c-983a-0a3cd2c74747")
	s.uuids["otherGroup"] = uuid.MustParse("b5f767d6-184c-42f3-8252-06cbd9e400c5")

	s.uuids["message"] = uuid.MustParse("2845fa37-7bcb-40bb-a447-6e0bc5b151b2")
	s.uuids["deletedMessage"] = uuid.MustParse("e3b82857-85fa-4b83-ae06-f85d63aab567")
	s.uuids["notExistingMessage"] = uuid.MustParse("2e5530f9-36cc-4186-b46c-821eb900ba4a")

	// Setting up our mocks
	mockDB := new(mockdb.MockMessageDB)

	mockDB.On("GetGroupMessages", s.uuids["userInGroup"], s.uuids["group"], 0, 4).Return([]models.Message{{Text: "elo", Nick: "Mal"},
		{Text: "siema", Nick: "River"},
		{Text: "elo elo", Nick: "Mal"},
		{Text: "siema siema", Nick: "River"}}, nil)
	mockDB.On("GetGroupMessages", s.uuids["userNotInGroup"], s.uuids["group"], 0, 4).Return([]models.Message{}, apperrors.NewForbidden("User cannot request from this group"))

	mockDB.On("DeleteMessageForYourself", s.uuids["userInGroup"], s.uuids["message"], s.uuids["group"]).Return(models.Message{Text: "valid"}, nil)
	mockDB.On("DeleteMessageForYourself", s.uuids["userNotInGroup"], s.uuids["message"], s.uuids["group"]).Return(models.Message{}, apperrors.NewForbidden("user not in group"))
	mockDB.On("DeleteMessageForYourself", s.uuids["userInGroup"], s.uuids["message"], s.uuids["otherGroup"]).Return(models.Message{}, apperrors.NewNotFound("message", s.uuids["message"].String()))
	mockDB.On("DeleteMessageForYourself", s.uuids["userInGroup"], s.uuids["notExistingMessage"], s.uuids["group"]).Return(models.Message{}, apperrors.NewNotFound("message", s.uuids["notExistingMessage"].String()))
	mockDB.On("DeleteMessageForYourself", s.uuids["userInGroup"], s.uuids["deletedMessage"], s.uuids["group"]).Return(models.Message{}, apperrors.NewConflict("deleted", s.uuids["userInGroup"].String()))

	mockDB.On("DeleteMessageForEveryone", s.uuids["userNotInGroup"], s.uuids["message"], s.uuids["group"]).Return(models.Message{}, apperrors.NewForbidden("user not in group"))
	mockDB.On("DeleteMessageForEveryone", s.uuids["userInGroup"], s.uuids["notExistingMessage"], s.uuids["group"]).Return(models.Message{}, apperrors.NewNotFound("message", s.uuids["notExistingMessage"].String()))
	mockDB.On("DeleteMessageForEveryone", s.uuids["userWithNoRights"], s.uuids["message"], s.uuids["group"]).Return(models.Message{}, apperrors.NewForbidden("User has no right to delete messages"))
	mockDB.On("DeleteMessageForEveryone", s.uuids["userInGroup"], s.uuids["message"], s.uuids["group"]).Return(models.Message{Text: "valid"}, nil)

	s.db = mockDB
}

func (s MessageTestSuite) TestGetGroupMessages() {
	gin.SetMode(gin.TestMode)

	server := handlers.NewServer(s.db, nil, nil, nil)

	testCases := []struct {
		desc               string
		userID             string
		returnVal          bool
		groupID            string
		expectedStatusCode int
		expectedResponse   interface{}
	}{
		{
			desc:               "getmessagessuccess",
			userID:             s.uuids["userInGroup"].String(),
			returnVal:          true,
			expectedStatusCode: http.StatusOK,
			groupID:            s.uuids["group"].String(),
			expectedResponse: []models.Message{{Text: "elo", Nick: "Mal"},
				{Text: "siema", Nick: "River"},
				{Text: "elo elo", Nick: "Mal"},
				{Text: "siema siema", Nick: "River"}},
		},
		{
			desc:               "getmessagesforbidden",
			userID:             s.uuids["userNotInGroup"].String(),
			returnVal:          false,
			groupID:            "61fbd273-b941-471c-983a-0a3cd2c74747",
			expectedStatusCode: http.StatusForbidden,
			expectedResponse:   gin.H{"err": "Forbidden action. Reason: User cannot request from this group"},
		},
		{
			desc:               "getmessagesnogroup",
			userID:             s.uuids["userInGroup"].String(),
			returnVal:          false,
			groupID:            "0",
			expectedStatusCode: http.StatusBadRequest,
			expectedResponse:   gin.H{"err": "invalid UUID length: 1"},
		},
	}

	for _, tC := range testCases {
		s.Run(tC.desc, func() {
			req, _ := http.NewRequest("GET", "/api/group/"+tC.groupID+"/messages?num=4&offset=0", nil)

			w := httptest.NewRecorder()
			_, engine := gin.CreateTestContext(w)
			engine.Use(func(c *gin.Context) {
				c.Set("userID", tC.userID)
			})

			engine.Handle(http.MethodGet, "/api/group/:groupID/messages", server.GetGroupMessages)
			engine.ServeHTTP(w, req)
			response := w.Result()

			s.Equal(tC.expectedStatusCode, response.StatusCode)

			var respBody interface{}
			if tC.returnVal {
				groups := []models.Message{}
				json.NewDecoder(response.Body).Decode(&groups)
				respBody = groups
			} else {
				var msg gin.H
				json.NewDecoder(response.Body).Decode(&msg)
				respBody = msg
			}

			s.True(reflect.DeepEqual(respBody, tC.expectedResponse))
		})
	}
}

func (s *MessageTestSuite) TestDeleteMessageForYourself() {
	gin.SetMode(gin.TestMode)

	// TokenClient, EventEmitter and EventListener are set to nil because there is no need
	// to instantiate them as they will not be called in this method
	server := handlers.NewServer(s.db, nil, nil, nil)

	testCases := []struct {
		desc               string
		userID             string
		groupID            string
		messageID          string
		returnVal          bool
		expectedStatusCode int
		expectedResponse   interface{}
	}{
		{
			desc:               "delYourselfInvalidUserID",
			userID:             "1",
			groupID:            s.uuids["group"].String(),
			messageID:          s.uuids["message"].String(),
			returnVal:          false,
			expectedStatusCode: http.StatusBadRequest,
			expectedResponse:   gin.H{"err": "invalid user ID"},
		},
		{
			desc:               "delYourselfInvalidGroupID",
			userID:             s.uuids["userInGroup"].String(),
			groupID:            "1",
			messageID:          s.uuids["message"].String(),
			returnVal:          false,
			expectedStatusCode: http.StatusBadRequest,
			expectedResponse:   gin.H{"err": "invalid group ID"},
		},
		{
			desc:               "delYourselfInvalidMessageID",
			userID:             s.uuids["userInGroup"].String(),
			groupID:            s.uuids["group"].String(),
			messageID:          "1",
			returnVal:          false,
			expectedStatusCode: http.StatusBadRequest,
			expectedResponse:   gin.H{"err": "invalid message ID"},
		},
		{
			desc:               "delYourselfUserNotInGroup",
			userID:             s.uuids["userNotInGroup"].String(),
			groupID:            s.uuids["group"].String(),
			messageID:          s.uuids["message"].String(),
			returnVal:          false,
			expectedStatusCode: http.StatusForbidden,
			expectedResponse:   gin.H{"err": "Forbidden action. Reason: user not in group"},
		},
		{
			desc:               "delYourselfMessageOfDifferentGroup",
			userID:             s.uuids["userInGroup"].String(),
			groupID:            s.uuids["otherGroup"].String(),
			messageID:          s.uuids["message"].String(),
			returnVal:          false,
			expectedStatusCode: http.StatusNotFound,
			expectedResponse:   gin.H{"err": "resource: message with value: 2845fa37-7bcb-40bb-a447-6e0bc5b151b2 not found"},
		},
		{
			desc:               "delYourselfMessageNotFound",
			userID:             s.uuids["userInGroup"].String(),
			groupID:            s.uuids["group"].String(),
			messageID:          s.uuids["notExistingMessage"].String(),
			returnVal:          false,
			expectedStatusCode: http.StatusNotFound,
			expectedResponse:   gin.H{"err": "resource: message with value: 2e5530f9-36cc-4186-b46c-821eb900ba4a not found"},
		},
		{
			desc:               "delYourselfMessageDeleted",
			userID:             s.uuids["userInGroup"].String(),
			groupID:            s.uuids["group"].String(),
			messageID:          s.uuids["deletedMessage"].String(),
			returnVal:          false,
			expectedStatusCode: http.StatusConflict,
			expectedResponse:   gin.H{"err": "resource: deleted with value: 1c4dccaf-a341-4920-9003-f24e0412f8e0 already exists"},
		},
		{
			desc:               "delYourselfSuccess",
			userID:             s.uuids["userInGroup"].String(),
			groupID:            s.uuids["group"].String(),
			messageID:          s.uuids["message"].String(),
			returnVal:          true,
			expectedStatusCode: http.StatusOK,
			expectedResponse:   models.Message{Text: "valid"},
		},
	}

	for _, tC := range testCases {

		s.Run(tC.desc, func() {
			req, _ := http.NewRequest(http.MethodPatch, "/groups/"+tC.groupID+"/messages/"+tC.messageID, nil)

			w := httptest.NewRecorder()
			_, engine := gin.CreateTestContext(w)
			engine.Use(func(ctx *gin.Context) {
				ctx.Set("userID", tC.userID)
			})

			engine.PATCH("/groups/:groupID/messages/:messageID", server.DeleteMessageForYourself)
			engine.ServeHTTP(w, req)
			response := w.Result()

			s.Equal(tC.expectedStatusCode, response.StatusCode)

			var respBody interface{}
			if tC.returnVal {
				var msg models.Message
				_ = json.NewDecoder(response.Body).Decode(&msg)
				respBody = msg
			} else {
				var msg gin.H
				_ = json.NewDecoder(response.Body).Decode(&msg)
				respBody = msg
			}

			s.True(reflect.DeepEqual(respBody, tC.expectedResponse))
		})
	}
}

func (s *MessageTestSuite) TestDeleteMessageForEveryone() {
	gin.SetMode(gin.TestMode)

	mockEmit := new(mockqueue.MockEmitter)
	mockEmit.On("Emit", mock.Anything).Return(nil)

	// TokenClient and EventListener are set to nil because there is no need
	// to instantiate them as they will not be called in this method
	server := handlers.NewServer(s.db, nil, mockEmit, nil)

	testCases := []struct {
		desc               string
		userID             string
		groupID            string
		messageID          string
		returnVal          bool
		expectedStatusCode int
		expectedResponse   interface{}
	}{
		{
			desc:               "DeleteForEveryone InvalidUserID",
			userID:             "1",
			groupID:            s.uuids["group"].String(),
			messageID:          s.uuids["message"].String(),
			returnVal:          false,
			expectedStatusCode: http.StatusBadRequest,
			expectedResponse:   gin.H{"err": "invalid user ID"},
		},
		{
			desc:               "DeleteForEveryone InvalidGroupID",
			userID:             s.uuids["userInGroup"].String(),
			groupID:            "1",
			messageID:          s.uuids["message"].String(),
			returnVal:          false,
			expectedStatusCode: http.StatusBadRequest,
			expectedResponse:   gin.H{"err": "invalid group ID"},
		},
		{
			desc:               "DeleteForEveryone InvalidMessageID",
			userID:             s.uuids["userInGroup"].String(),
			groupID:            s.uuids["group"].String(),
			messageID:          "1",
			returnVal:          false,
			expectedStatusCode: http.StatusBadRequest,
			expectedResponse:   gin.H{"err": "invalid message ID"},
		},
		{
			desc:               "DeleteForEveryone UserNotInGroup",
			userID:             s.uuids["userNotInGroup"].String(),
			groupID:            s.uuids["group"].String(),
			messageID:          s.uuids["message"].String(),
			returnVal:          false,
			expectedStatusCode: http.StatusForbidden,
			expectedResponse:   gin.H{"err": "Forbidden action. Reason: user not in group"},
		},
		{
			desc:               "DeleteForEveryone MessageNotFound",
			userID:             s.uuids["userInGroup"].String(),
			groupID:            s.uuids["group"].String(),
			messageID:          s.uuids["notExistingMessage"].String(),
			returnVal:          false,
			expectedStatusCode: http.StatusNotFound,
			expectedResponse:   gin.H{"err": "resource: message with value: 2e5530f9-36cc-4186-b46c-821eb900ba4a not found"},
		},
		{
			desc:               "DeleteForEveryone NoRights",
			userID:             s.uuids["userWithNoRights"].String(),
			groupID:            s.uuids["group"].String(),
			messageID:          s.uuids["message"].String(),
			returnVal:          false,
			expectedStatusCode: http.StatusForbidden,
			expectedResponse:   gin.H{"err": "Forbidden action. Reason: User has no right to delete messages"},
		},
		{
			desc:               "DeleteForEveryone Success",
			userID:             s.uuids["userInGroup"].String(),
			groupID:            s.uuids["group"].String(),
			messageID:          s.uuids["message"].String(),
			returnVal:          true,
			expectedStatusCode: http.StatusOK,
			expectedResponse:   models.Message{Text: "valid"},
		},
	}

	for _, tC := range testCases {

		s.Run(tC.desc, func() {
			req, _ := http.NewRequest(http.MethodDelete, "/groups/"+tC.groupID+"/messages/"+tC.messageID, nil)

			w := httptest.NewRecorder()
			_, engine := gin.CreateTestContext(w)
			engine.Use(func(ctx *gin.Context) {
				ctx.Set("userID", tC.userID)
			})

			engine.DELETE("/groups/:groupID/messages/:messageID", server.DeleteMessageForEveryone)
			engine.ServeHTTP(w, req)
			response := w.Result()

			s.Equal(tC.expectedStatusCode, response.StatusCode)

			var respBody interface{}
			if tC.returnVal {
				var msg models.Message
				_ = json.NewDecoder(response.Body).Decode(&msg)
				respBody = msg
			} else {
				var msg gin.H
				_ = json.NewDecoder(response.Body).Decode(&msg)
				respBody = msg
			}

			s.True(reflect.DeepEqual(respBody, tC.expectedResponse))
		})
	}
}

func TestMessageSuite(t *testing.T) {
	suite.Run(t, &MessageTestSuite{})
}
