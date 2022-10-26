package handlers_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/Slimo300/MicroservicesChatApp/backend/lib/apperrors"
	"github.com/Slimo300/MicroservicesChatApp/backend/message-service/database/mock"
	"github.com/Slimo300/MicroservicesChatApp/backend/message-service/handlers"
	"github.com/Slimo300/MicroservicesChatApp/backend/message-service/models"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func TestGetGroupMessages(t *testing.T) {
	gin.SetMode(gin.TestMode)

	group := uuid.MustParse("61fbd273-b941-471c-983a-0a3cd2c74747")
	user1 := uuid.MustParse("1c4dccaf-a341-4920-9003-f24e0412f8e0")
	user2 := uuid.MustParse("634240cf-1219-4be2-adfa-90ab6b47899b")

	mockDB := new(mock.MockMessageDB)
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

			json.NewDecoder(response.Body).Decode(&respBody)
			if !reflect.DeepEqual(respBody, tC.expectedResponse) {
				t.Errorf("Received HTTP response body %+v does not match expected HTTP response Body %+v", respBody, tC.expectedResponse)
			}
		})
	}
}
