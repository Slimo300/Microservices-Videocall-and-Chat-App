package handlers_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
	"time"

	"github.com/Slimo300/MicroservicesChatApp/backend/message-service/database"
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
	message1 := uuid.MustParse("f6008d3c-965e-46e6-8ec9-ad7d4da02e93")
	message2 := uuid.MustParse("c108dbe0-b6d2-4cb1-8a61-f5b03875d41e")
	message3 := uuid.MustParse("54bc782e-6140-4558-82a5-c2d2460d6325")
	message4 := uuid.MustParse("34787baa-91f6-4a83-b4fc-e56d0d37e4f6")
	member1 := uuid.MustParse("e4372b71-30ca-42e1-8c1e-7df6d033fd3f")
	member2 := uuid.MustParse("b38aaff8-6733-4a1d-8eaf-fc10e656d02b")
	posted1, _ := time.Parse("2006-02-01 15:04:05", "2019-13-01 22:00:45")
	posted2, _ := time.Parse("2006-02-01 15:04:05", "2019-15-01 22:00:45")
	posted3, _ := time.Parse("2006-02-01 15:04:05", "2019-16-01 22:00:45")
	posted4, _ := time.Parse("2006-02-01 15:04:05", "2019-17-01 22:00:45")

	mockDB := new(database.MockMessageDB)
	mockDB.On("IsUserInGroup", user1, group).Return(true)
	mockDB.On("IsUserInGroup", user2, group).Return(false)
	mockDB.On("GetGroupMessages", group, 0, 4).Return([]models.Message{{ID: message1, GroupID: group, MemberID: member1, Text: "elo", Nick: "Mal", Posted: posted1},
		{ID: message2, GroupID: group, MemberID: member2, Text: "siema", Nick: "River", Posted: posted2},
		{ID: message3, GroupID: group, MemberID: member1, Text: "elo elo", Nick: "Mal", Posted: posted3},
		{ID: message4, GroupID: group, MemberID: member2, Text: "siema siema", Nick: "River", Posted: posted4}}, nil)

	s := handlers.Server{
		DB: mockDB,
	}

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
			userID:             user1.String(),
			returnVal:          true,
			expectedStatusCode: http.StatusOK,
			group:              group.String(),
			expectedResponse: []models.Message{{ID: message1, GroupID: group, MemberID: member1, Text: "elo", Nick: "Mal", Posted: posted1},
				{ID: message2, GroupID: group, MemberID: member2, Text: "siema", Nick: "River", Posted: posted2},
				{ID: message3, GroupID: group, MemberID: member1, Text: "elo elo", Nick: "Mal", Posted: posted3},
				{ID: message4, GroupID: group, MemberID: member2, Text: "siema siema", Nick: "River", Posted: posted4}},
		},
		{
			desc:               "getmessagesforbidden",
			userID:             user2.String(),
			returnVal:          false,
			group:              group.String(),
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
