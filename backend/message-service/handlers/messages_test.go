package handlers_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/Slimo300/MicroservicesChatApp/backend/lib/communication"
	"github.com/Slimo300/MicroservicesChatApp/backend/message-service/database"
	"github.com/Slimo300/MicroservicesChatApp/backend/message-service/handlers"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func TestGetGroupMessages(t *testing.T) {
	gin.SetMode(gin.TestMode)
	s := handlers.Server{
		DB: new(database.MockMessageDB),
	}

	groupId, _ := uuid.Parse("61fbd273-b941-471c-983a-0a3cd2c74747")
	member1, _ := uuid.Parse("e4372b71-30ca-42e1-8c1e-7df6d033fd3f")
	member2, _ := uuid.Parse("b38aaff8-6733-4a1d-8eaf-fc10e656d02b")

	testCases := []struct {
		desc               string
		userID             string
		returnVal          bool
		group              string
		expectedStatusCode int
		expectedResponse   interface{}
	}{
		{
			desc:               "getmessagessuccess",
			userID:             "1c4dccaf-a341-4920-9003-f24e0412f8e0",
			returnVal:          true,
			expectedStatusCode: http.StatusOK,
			group:              "61fbd273-b941-471c-983a-0a3cd2c74747",
			expectedResponse: []communication.Message{{Group: groupId, Member: member1, Message: "elo", Nick: "Mal", When: "2019-13-01 22:00:45"},
				{Group: groupId, Member: member2, Message: "siema", Nick: "River", When: "2019-15-01 22:00:45"},
				{Group: groupId, Member: member1, Message: "elo elo", Nick: "Mal", When: "2019-16-01 22:00:45"},
				{Group: groupId, Member: member2, Message: "siema siema", Nick: "River", When: "2019-17-01 22:00:45"}},
		},
		{
			desc:               "getmessagesforbidden",
			userID:             "634240cf-1219-4be2-adfa-90ab6b47899b",
			returnVal:          false,
			group:              "61fbd273-b941-471c-983a-0a3cd2c74747",
			expectedStatusCode: http.StatusForbidden,
			expectedResponse:   gin.H{"err": "User cannot request from this group"},
		},
		{
			desc:               "getmessagesnogroup",
			userID:             "1c4dccaf-a341-4920-9003-f24e0412f8e0",
			returnVal:          false,
			group:              "0",
			expectedStatusCode: http.StatusBadRequest,
			expectedResponse:   gin.H{"err": "invalid group ID"},
		},
	}

	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			req, _ := http.NewRequest("GET", "/api/group/"+tC.group+"/messages?num=4&offset=0", nil)

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
				groups := []communication.Message{}
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
