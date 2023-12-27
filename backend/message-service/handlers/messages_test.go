package handlers_test

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/Slimo300/Microservices-Videocall-and-Chat-App/backend/lib/apperrors"
	mockqueue "github.com/Slimo300/Microservices-Videocall-and-Chat-App/backend/lib/msgqueue/mock"
	mockdb "github.com/Slimo300/Microservices-Videocall-and-Chat-App/backend/message-service/database/mock"
	"github.com/Slimo300/Microservices-Videocall-and-Chat-App/backend/message-service/handlers"
	"github.com/Slimo300/Microservices-Videocall-and-Chat-App/backend/message-service/models"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type MessageTestSuite struct {
	suite.Suite
	uuids  map[string]uuid.UUID
	server handlers.Server
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

	s.uuids["group"] = uuid.MustParse("1c4dccaf-a341-4920-9023-f24e0412f8e0")

	s.uuids["userInGroup"] = uuid.MustParse("1c4dccaf-a341-4920-9003-f24e0412f8e0")
	s.uuids["userNotInGroup"] = uuid.MustParse("634240cf-1219-4be2-adfa-90ab6b47899b")
	s.uuids["userWithNoRights"] = uuid.MustParse("889858f7-4b2a-4914-ad14-471a7300c26b")

	s.uuids["message"] = uuid.MustParse("2845fa37-7bcb-40bb-a447-6e0bc5b151b2")
	s.uuids["deletedMessage"] = uuid.MustParse("e3b82857-85fa-4b83-ae06-f85d63aab567")
	s.uuids["notExistingMessage"] = uuid.MustParse("2e5530f9-36cc-4186-b46c-821eb900ba4a")

	// Setting up our mocks
	mockDB := new(mockdb.MockMessageDB)

	mockDB.On("GetGroupMessages", s.uuids["userInGroup"], s.uuids["group"], 0, 4).Return([]models.Message{{Text: "elo"},
		{Text: "siema"},
		{Text: "elo elo"},
		{Text: "siema siema"}}, nil)
	mockDB.On("GetGroupMessages", s.uuids["userNotInGroup"], s.uuids["group"], 0, 4).Return([]models.Message{}, apperrors.NewForbidden("User cannot request from this group"))

	mockDB.On("DeleteMessageForYourself", s.uuids["userInGroup"], s.uuids["message"]).Return(&models.Message{Text: "valid"}, nil)
	mockDB.On("DeleteMessageForYourself", s.uuids["userNotInGroup"], s.uuids["message"]).Return(nil, apperrors.NewNotFound(fmt.Sprintf("message with id %s not found", s.uuids["message"])))
	mockDB.On("DeleteMessageForYourself", s.uuids["userInGroup"], s.uuids["notExistingMessage"]).Return(nil, apperrors.NewNotFound(fmt.Sprintf("Message with id %v not found", s.uuids["notExistingMessage"].String())))
	mockDB.On("DeleteMessageForYourself", s.uuids["userInGroup"], s.uuids["deletedMessage"]).Return(nil, apperrors.NewConflict(fmt.Sprintf("Message %v already deleted", s.uuids["deletedMessage"].String())))

	mockDB.On("DeleteMessageForEveryone", s.uuids["userNotInGroup"], s.uuids["message"]).Return(nil, apperrors.NewNotFound(fmt.Sprintf("message with id %s not found", s.uuids["message"])))
	mockDB.On("DeleteMessageForEveryone", s.uuids["userInGroup"], s.uuids["notExistingMessage"]).Return(nil, apperrors.NewNotFound(fmt.Sprintf("Message with id %v not found", s.uuids["notExistingMessage"].String())))
	mockDB.On("DeleteMessageForEveryone", s.uuids["userWithNoRights"], s.uuids["message"]).Return(nil, apperrors.NewForbidden("User has no right to delete messages"))
	mockDB.On("DeleteMessageForEveryone", s.uuids["userInGroup"], s.uuids["message"]).Return(&models.Message{Text: "valid"}, nil)

	mockEmit := new(mockqueue.MockEmitter)
	mockEmit.On("Emit", mock.Anything).Return(nil)

	s.server = *handlers.NewServer(
		mockDB,
		nil,
		mockEmit,
		nil,
	)
}

func (s *MessageTestSuite) TestGetGroupMessages() {
	gin.SetMode(gin.TestMode)

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
			expectedResponse: []models.Message{{Text: "elo"},
				{Text: "siema"},
				{Text: "elo elo"},
				{Text: "siema siema"}},
		},
		{
			desc:               "getmessagesforbidden",
			userID:             s.uuids["userNotInGroup"].String(),
			returnVal:          false,
			groupID:            s.uuids["group"].String(),
			expectedStatusCode: http.StatusForbidden,
			expectedResponse:   gin.H{"err": "User cannot request from this group"},
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

			engine.Handle(http.MethodGet, "/api/group/:groupID/messages", s.server.GetGroupMessages)
			engine.ServeHTTP(w, req)
			response := w.Result()
			defer response.Body.Close()

			s.Equal(tC.expectedStatusCode, response.StatusCode)

			var respBody interface{}
			if tC.returnVal {
				groups := []models.Message{}
				_ = json.NewDecoder(response.Body).Decode(&groups)
				respBody = groups
			} else {
				var msg gin.H
				_ = json.NewDecoder(response.Body).Decode(&msg)
				respBody = msg
			}

			s.True(reflect.DeepEqual(respBody, tC.expectedResponse))
		})
	}
}

func (s *MessageTestSuite) TestDeleteMessageForYourself() {
	gin.SetMode(gin.TestMode)

	testCases := []struct {
		desc               string
		userID             string
		messageID          string
		returnVal          bool
		expectedStatusCode int
		expectedResponse   interface{}
	}{
		{
			desc:               "delYourselfInvalidUserID",
			userID:             "1",
			messageID:          s.uuids["message"].String(),
			returnVal:          false,
			expectedStatusCode: http.StatusBadRequest,
			expectedResponse:   gin.H{"err": "invalid user ID"},
		},
		{
			desc:               "delYourselfInvalidMessageID",
			userID:             s.uuids["userInGroup"].String(),
			messageID:          "1",
			returnVal:          false,
			expectedStatusCode: http.StatusBadRequest,
			expectedResponse:   gin.H{"err": "invalid message ID"},
		},
		{
			desc:               "delYourselfUserNotInGroup",
			userID:             s.uuids["userNotInGroup"].String(),
			messageID:          s.uuids["message"].String(),
			returnVal:          false,
			expectedStatusCode: http.StatusNotFound,
			expectedResponse:   gin.H{"err": "message with id 2845fa37-7bcb-40bb-a447-6e0bc5b151b2 not found"},
		},
		{
			desc:               "delYourselfMessageNotFound",
			userID:             s.uuids["userInGroup"].String(),
			messageID:          s.uuids["notExistingMessage"].String(),
			returnVal:          false,
			expectedStatusCode: http.StatusNotFound,
			expectedResponse:   gin.H{"err": "Message with id 2e5530f9-36cc-4186-b46c-821eb900ba4a not found"},
		},
		{
			desc:               "delYourselfMessageDeleted",
			userID:             s.uuids["userInGroup"].String(),
			messageID:          s.uuids["deletedMessage"].String(),
			returnVal:          false,
			expectedStatusCode: http.StatusConflict,
			expectedResponse:   gin.H{"err": "Message e3b82857-85fa-4b83-ae06-f85d63aab567 already deleted"},
		},
		{
			desc:               "delYourselfSuccess",
			userID:             s.uuids["userInGroup"].String(),
			messageID:          s.uuids["message"].String(),
			returnVal:          true,
			expectedStatusCode: http.StatusOK,
			expectedResponse:   &models.Message{Text: "valid"},
		},
	}

	for _, tC := range testCases {

		s.Run(tC.desc, func() {
			req, _ := http.NewRequest(http.MethodPatch, "/messages/"+tC.messageID+"/hide", nil)

			w := httptest.NewRecorder()
			_, engine := gin.CreateTestContext(w)
			engine.Use(func(ctx *gin.Context) {
				ctx.Set("userID", tC.userID)
			})

			engine.PATCH("/messages/:messageID/hide", s.server.DeleteMessageForYourself)
			engine.ServeHTTP(w, req)
			response := w.Result()
			defer response.Body.Close()

			s.Equal(tC.expectedStatusCode, response.StatusCode)

			var respBody interface{}
			if tC.returnVal {
				var msg *models.Message
				_ = json.NewDecoder(response.Body).Decode(&msg)
				respBody = msg
			} else {
				var msg gin.H
				_ = json.NewDecoder(response.Body).Decode(&msg)
				respBody = msg
			}

			s.Equal(respBody, tC.expectedResponse)
		})
	}
}

func (s *MessageTestSuite) TestDeleteMessageForEveryone() {
	gin.SetMode(gin.TestMode)

	testCases := []struct {
		desc               string
		userID             string
		messageID          string
		returnVal          bool
		expectedStatusCode int
		expectedResponse   interface{}
	}{
		{
			desc:               "DeleteForEveryone InvalidUserID",
			userID:             "1",
			messageID:          s.uuids["message"].String(),
			returnVal:          false,
			expectedStatusCode: http.StatusBadRequest,
			expectedResponse:   gin.H{"err": "invalid user ID"},
		},
		{
			desc:               "DeleteForEveryone InvalidMessageID",
			userID:             s.uuids["userInGroup"].String(),
			messageID:          "1",
			returnVal:          false,
			expectedStatusCode: http.StatusBadRequest,
			expectedResponse:   gin.H{"err": "invalid message ID"},
		},
		{
			desc:               "DeleteForEveryone UserNotInGroup",
			userID:             s.uuids["userNotInGroup"].String(),
			messageID:          s.uuids["message"].String(),
			returnVal:          false,
			expectedStatusCode: http.StatusNotFound,
			expectedResponse:   gin.H{"err": "message with id 2845fa37-7bcb-40bb-a447-6e0bc5b151b2 not found"},
		},
		{
			desc:               "DeleteForEveryone MessageNotFound",
			userID:             s.uuids["userInGroup"].String(),
			messageID:          s.uuids["notExistingMessage"].String(),
			returnVal:          false,
			expectedStatusCode: http.StatusNotFound,
			expectedResponse:   gin.H{"err": "Message with id 2e5530f9-36cc-4186-b46c-821eb900ba4a not found"},
		},
		{
			desc:               "DeleteForEveryone NoRights",
			userID:             s.uuids["userWithNoRights"].String(),
			messageID:          s.uuids["message"].String(),
			returnVal:          false,
			expectedStatusCode: http.StatusForbidden,
			expectedResponse:   gin.H{"err": "User has no right to delete messages"},
		},
		{
			desc:               "DeleteForEveryone Success",
			userID:             s.uuids["userInGroup"].String(),
			messageID:          s.uuids["message"].String(),
			returnVal:          true,
			expectedStatusCode: http.StatusOK,
			expectedResponse:   &models.Message{Text: "valid"},
		},
	}

	for _, tC := range testCases {

		s.Run(tC.desc, func() {
			req, _ := http.NewRequest(http.MethodDelete, "/messages/"+tC.messageID, nil)

			w := httptest.NewRecorder()
			_, engine := gin.CreateTestContext(w)
			engine.Use(func(ctx *gin.Context) {
				ctx.Set("userID", tC.userID)
			})

			engine.DELETE("/messages/:messageID", s.server.DeleteMessageForEveryone)
			engine.ServeHTTP(w, req)
			response := w.Result()
			defer response.Body.Close()

			s.Equal(tC.expectedStatusCode, response.StatusCode)

			var respBody interface{}
			if tC.returnVal {
				var msg *models.Message
				_ = json.NewDecoder(response.Body).Decode(&msg)
				respBody = msg
			} else {
				var msg gin.H
				_ = json.NewDecoder(response.Body).Decode(&msg)
				respBody = msg
			}

			s.Equal(respBody, tC.expectedResponse)
		})
	}
}

func TestMessageSuite(t *testing.T) {
	suite.Run(t, &MessageTestSuite{})
}
