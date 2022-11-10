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

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func TestGetGroupMessages(t *testing.T) {
	gin.SetMode(gin.TestMode)

	group := uuid.MustParse("61fbd273-b941-471c-983a-0a3cd2c74747")
	user1 := uuid.MustParse("1c4dccaf-a341-4920-9003-f24e0412f8e0")
	user2 := uuid.MustParse("634240cf-1219-4be2-adfa-90ab6b47899b")

	mockDB := new(mockdb.MockMessageDB)
	mockDB.On("GetGroupMessages", user1, group, 0, 4).Return([]models.Message{{Text: "elo", Nick: "Mal"},
		{Text: "siema", Nick: "River"},
		{Text: "elo elo", Nick: "Mal"},
		{Text: "siema siema", Nick: "River"}}, nil)
	mockDB.On("GetGroupMessages", user2, group, 0, 4).Return([]models.Message{}, apperrors.NewForbidden("User cannot request from this group"))

	s := handlers.Server{
		DB: mockDB,
	}

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
			userID:             "1c4dccaf-a341-4920-9003-f24e0412f8e0",
			returnVal:          true,
			expectedStatusCode: http.StatusOK,
			groupID:            "61fbd273-b941-471c-983a-0a3cd2c74747",
			expectedResponse: []models.Message{{Text: "elo", Nick: "Mal"},
				{Text: "siema", Nick: "River"},
				{Text: "elo elo", Nick: "Mal"},
				{Text: "siema siema", Nick: "River"}},
		},
		{
			desc:               "getmessagesforbidden",
			userID:             "634240cf-1219-4be2-adfa-90ab6b47899b",
			returnVal:          false,
			groupID:            "61fbd273-b941-471c-983a-0a3cd2c74747",
			expectedStatusCode: http.StatusForbidden,
			expectedResponse:   gin.H{"err": "Forbidden action. Reason: User cannot request from this group"},
		},
		{
			desc:               "getmessagesnogroup",
			userID:             "1c4dccaf-a341-4920-9003-f24e0412f8e0",
			returnVal:          false,
			groupID:            "0",
			expectedStatusCode: http.StatusBadRequest,
			expectedResponse:   gin.H{"err": "invalid UUID length: 1"},
		},
	}

	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			req, _ := http.NewRequest("GET", "/api/group/"+tC.groupID+"/messages?num=4&offset=0", nil)

			w := httptest.NewRecorder()
			_, engine := gin.CreateTestContext(w)
			engine.Use(func(c *gin.Context) {
				c.Set("userID", tC.userID)
			})

			engine.Handle(http.MethodGet, "/api/group/:groupID/messages", s.GetGroupMessages)
			engine.ServeHTTP(w, req)
			response := w.Result()

			if response.StatusCode != tC.expectedStatusCode {
				t.Errorf("Received Status code %d does not match expected status %d", response.StatusCode, tC.expectedStatusCode)
			}
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

			if !reflect.DeepEqual(respBody, tC.expectedResponse) {
				t.Errorf("Received HTTP response body %+v does not match expected HTTP response Body %+v", respBody, tC.expectedResponse)
			}
		})
	}
}

func TestDeleteMessageForYourself(t *testing.T) {
	gin.SetMode(gin.TestMode)

	// Modelled database:
	// 2 users, 2 groups, "userInGroup" belongs to both groups, "userNotInGroup" to none,
	// "otherGroup" has no messages and is there to make sure user won't delete message if he doesn't
	// know to what group it belongs
	// 2 messages (message, deletedMessage) belong to "group" and "userInGroup"
	// "notExistingMessage" is for checking service reaction to sending unknown ID

	group := uuid.MustParse("61fbd273-b941-471c-983a-0a3cd2c74747")
	otherGroup := uuid.MustParse("b5f767d6-184c-42f3-8252-06cbd9e400c5")

	userInGroup := uuid.MustParse("1c4dccaf-a341-4920-9003-f24e0412f8e0")
	userNotInGroup := uuid.MustParse("634240cf-1219-4be2-adfa-90ab6b47899b")

	message := uuid.MustParse("2845fa37-7bcb-40bb-a447-6e0bc5b151b2")
	deletedMessage := uuid.MustParse("e3b82857-85fa-4b83-ae06-f85d63aab567")
	notExistingMessage := uuid.MustParse("2e5530f9-36cc-4186-b46c-821eb900ba4a")

	mockDB := new(mockdb.MockMessageDB)
	mockDB.On("DeleteMessageForYourself", userInGroup, message, group).Return(models.Message{Text: "valid"}, nil)
	mockDB.On("DeleteMessageForYourself", userNotInGroup, message, group).Return(models.Message{}, apperrors.NewForbidden("user not in group"))
	mockDB.On("DeleteMessageForYourself", userInGroup, message, otherGroup).Return(models.Message{}, apperrors.NewNotFound("message", message.String()))
	mockDB.On("DeleteMessageForYourself", userInGroup, notExistingMessage, group).Return(models.Message{}, apperrors.NewNotFound("message", notExistingMessage.String()))
	mockDB.On("DeleteMessageForYourself", userInGroup, deletedMessage, group).Return(models.Message{}, apperrors.NewConflict("deleted", userInGroup.String()))

	// TokenClient, EventEmitter and EventListener are set to nil because there is no need
	// to instantiate them as they will not be called in this method
	server := handlers.NewServer(mockDB, nil, nil, nil)

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
			groupID:            group.String(),
			messageID:          message.String(),
			returnVal:          false,
			expectedStatusCode: http.StatusBadRequest,
			expectedResponse:   gin.H{"err": "invalid user ID"},
		},
		{
			desc:               "delYourselfInvalidGroupID",
			userID:             userInGroup.String(),
			groupID:            "1",
			messageID:          message.String(),
			returnVal:          false,
			expectedStatusCode: http.StatusBadRequest,
			expectedResponse:   gin.H{"err": "invalid group ID"},
		},
		{
			desc:               "delYourselfInvalidMessageID",
			userID:             userInGroup.String(),
			groupID:            group.String(),
			messageID:          "1",
			returnVal:          false,
			expectedStatusCode: http.StatusBadRequest,
			expectedResponse:   gin.H{"err": "invalid message ID"},
		},
		{
			desc:               "delYourselfUserNotInGroup",
			userID:             userNotInGroup.String(),
			groupID:            group.String(),
			messageID:          message.String(),
			returnVal:          false,
			expectedStatusCode: http.StatusForbidden,
			expectedResponse:   gin.H{"err": "Forbidden action. Reason: user not in group"},
		},
		{
			desc:               "delYourselfMessageOfDifferentGroup",
			userID:             userInGroup.String(),
			groupID:            otherGroup.String(),
			messageID:          message.String(),
			returnVal:          false,
			expectedStatusCode: http.StatusNotFound,
			expectedResponse:   gin.H{"err": "resource: message with value: 2845fa37-7bcb-40bb-a447-6e0bc5b151b2 not found"},
		},
		{
			desc:               "delYourselfMessageNotFound",
			userID:             userInGroup.String(),
			groupID:            group.String(),
			messageID:          notExistingMessage.String(),
			returnVal:          false,
			expectedStatusCode: http.StatusNotFound,
			expectedResponse:   gin.H{"err": "resource: message with value: 2e5530f9-36cc-4186-b46c-821eb900ba4a not found"},
		},
		{
			desc:               "delYourselfMessageDeleted",
			userID:             userInGroup.String(),
			groupID:            group.String(),
			messageID:          deletedMessage.String(),
			returnVal:          false,
			expectedStatusCode: http.StatusConflict,
			expectedResponse:   gin.H{"err": "resource: deleted with value: 1c4dccaf-a341-4920-9003-f24e0412f8e0 already exists"},
		},
		{
			desc:               "delYourselfSuccess",
			userID:             userInGroup.String(),
			groupID:            group.String(),
			messageID:          message.String(),
			returnVal:          true,
			expectedStatusCode: http.StatusOK,
			expectedResponse:   models.Message{Text: "valid"},
		},
	}

	for _, tC := range testCases {

		t.Run(tC.desc, func(t *testing.T) {
			req, _ := http.NewRequest(http.MethodPatch, "/groups/"+tC.groupID+"/messages/"+tC.messageID, nil)

			w := httptest.NewRecorder()
			_, engine := gin.CreateTestContext(w)
			engine.Use(func(ctx *gin.Context) {
				ctx.Set("userID", tC.userID)
			})

			engine.PATCH("/groups/:groupID/messages/:messageID", server.DeleteMessageForYourself)
			engine.ServeHTTP(w, req)
			response := w.Result()

			if tC.expectedStatusCode != response.StatusCode {
				t.Errorf("Received Status code %d does not match expected status %d", response.StatusCode, tC.expectedStatusCode)
			}

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

			if !reflect.DeepEqual(respBody, tC.expectedResponse) {
				t.Errorf("Received HTTP response body %+v does not match expected HTTP response Body %+v", respBody, tC.expectedResponse)
			}
		})
	}
}

func TestDeleteMessageForEveryone(t *testing.T) {
	gin.SetMode(gin.TestMode)

	// Modelled database:
	// 2 users, 2 groups, "userInGroup" belongs to both groups, "userNotInGroup" to none,
	// "otherGroup" has no messages and is there to make sure user won't delete message if he doesn't
	// know to what group it belongs
	// 2 messages (message, deletedMessage) belong to "group" and "userInGroup"
	// "notExistingMessage" is for checking service reaction to sending unknown ID

	group := uuid.MustParse("61fbd273-b941-471c-983a-0a3cd2c74747")

	userInGroup := uuid.MustParse("1c4dccaf-a341-4920-9003-f24e0412f8e0")
	userWithNoRights := uuid.MustParse("889858f7-4b2a-4914-ad14-471a7300c26b")
	userNotInGroup := uuid.MustParse("634240cf-1219-4be2-adfa-90ab6b47899b")

	message := uuid.MustParse("2845fa37-7bcb-40bb-a447-6e0bc5b151b2")
	notExistingMessage := uuid.MustParse("2e5530f9-36cc-4186-b46c-821eb900ba4a")

	// mock DB init
	mockDB := new(mockdb.MockMessageDB)
	mockDB.On("DeleteMessageForEveryone", userNotInGroup, message, group).Return(models.Message{}, apperrors.NewForbidden("user not in group"))
	mockDB.On("DeleteMessageForEveryone", userInGroup, notExistingMessage, group).Return(models.Message{}, apperrors.NewNotFound("message", notExistingMessage.String()))
	mockDB.On("DeleteMessageForEveryone", userWithNoRights, message, group).Return(models.Message{}, apperrors.NewForbidden("User has no right to delete messages"))
	mockDB.On("DeleteMessageForEveryone", userInGroup, message, group).Return(models.Message{Text: "valid"}, nil)

	mockEmit := new(mockqueue.MockEmitter)
	mockEmit.On("Emit", mock.Anything).Return(nil)

	// TokenClient and EventListener are set to nil because there is no need
	// to instantiate them as they will not be called in this method
	server := handlers.NewServer(mockDB, nil, mockEmit, nil)

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
			groupID:            group.String(),
			messageID:          message.String(),
			returnVal:          false,
			expectedStatusCode: http.StatusBadRequest,
			expectedResponse:   gin.H{"err": "invalid user ID"},
		},
		{
			desc:               "DeleteForEveryone InvalidGroupID",
			userID:             userInGroup.String(),
			groupID:            "1",
			messageID:          message.String(),
			returnVal:          false,
			expectedStatusCode: http.StatusBadRequest,
			expectedResponse:   gin.H{"err": "invalid group ID"},
		},
		{
			desc:               "DeleteForEveryone InvalidMessageID",
			userID:             userInGroup.String(),
			groupID:            group.String(),
			messageID:          "1",
			returnVal:          false,
			expectedStatusCode: http.StatusBadRequest,
			expectedResponse:   gin.H{"err": "invalid message ID"},
		},
		{
			desc:               "DeleteForEveryone UserNotInGroup",
			userID:             userNotInGroup.String(),
			groupID:            group.String(),
			messageID:          message.String(),
			returnVal:          false,
			expectedStatusCode: http.StatusForbidden,
			expectedResponse:   gin.H{"err": "Forbidden action. Reason: user not in group"},
		},
		{
			desc:               "DeleteForEveryone MessageNotFound",
			userID:             userInGroup.String(),
			groupID:            group.String(),
			messageID:          notExistingMessage.String(),
			returnVal:          false,
			expectedStatusCode: http.StatusNotFound,
			expectedResponse:   gin.H{"err": "resource: message with value: 2e5530f9-36cc-4186-b46c-821eb900ba4a not found"},
		},
		{
			desc:               "DeleteForEveryone NoRights",
			userID:             userWithNoRights.String(),
			groupID:            group.String(),
			messageID:          message.String(),
			returnVal:          false,
			expectedStatusCode: http.StatusForbidden,
			expectedResponse:   gin.H{"err": "Forbidden action. Reason: User has no right to delete messages"},
		},
		{
			desc:               "DeleteForEveryone Success",
			userID:             userInGroup.String(),
			groupID:            group.String(),
			messageID:          message.String(),
			returnVal:          true,
			expectedStatusCode: http.StatusOK,
			expectedResponse:   models.Message{Text: "valid"},
		},
	}

	for _, tC := range testCases {

		t.Run(tC.desc, func(t *testing.T) {
			req, _ := http.NewRequest(http.MethodDelete, "/groups/"+tC.groupID+"/messages/"+tC.messageID, nil)

			w := httptest.NewRecorder()
			_, engine := gin.CreateTestContext(w)
			engine.Use(func(ctx *gin.Context) {
				ctx.Set("userID", tC.userID)
			})

			engine.DELETE("/groups/:groupID/messages/:messageID", server.DeleteMessageForEveryone)
			engine.ServeHTTP(w, req)
			response := w.Result()

			if tC.expectedStatusCode != response.StatusCode {
				t.Errorf("Received Status code %d does not match expected status %d", response.StatusCode, tC.expectedStatusCode)
			}

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

			if !reflect.DeepEqual(respBody, tC.expectedResponse) {
				t.Errorf("Received HTTP response body %+v does not match expected HTTP response Body %+v", respBody, tC.expectedResponse)
			}
		})
	}
}
