package handlers_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
	"time"

	"github.com/Slimo300/MicroservicesChatApp/backend/group-service/models"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func TestSendGroupInvite(t *testing.T) {
	gin.SetMode(gin.TestMode)
	s := setupTestServer()

	testCases := []struct {
		desc               string
		id                 string
		data               map[string]interface{}
		expectedStatusCode int
		expectedResponse   interface{}
	}{
		{
			desc:               "invitesuccess",
			id:                 "1c4dccaf-a341-4920-9003-f24e0412f8e0",
			data:               map[string]interface{}{"group": "61fbd273-b941-471c-983a-0a3cd2c74747", "target": "Kel"},
			expectedStatusCode: http.StatusCreated,
			expectedResponse:   gin.H{"message": "invite sent"},
		},
		{
			desc:               "invitenosuchuser",
			id:                 "1c4dccaf-a341-4920-9003-f24e0412f8e0",
			data:               map[string]interface{}{"group": "61fbd273-b941-471c-983a-0a3cd2c74747", "target": "Raul"},
			expectedStatusCode: http.StatusNotFound,
			expectedResponse:   gin.H{"err": "no user with name: Raul"},
		},
		{
			desc:               "invitenorights",
			id:                 "0ef41409-24b0-43e6-80a3-cf31a4b1a684",
			data:               map[string]interface{}{"group": "61fbd273-b941-471c-983a-0a3cd2c74747", "target": "Kel"},
			expectedStatusCode: http.StatusForbidden,
			expectedResponse:   gin.H{"err": "no rights to add"},
		},
		{
			desc:               "inviteuserismember",
			id:                 "1c4dccaf-a341-4920-9003-f24e0412f8e0",
			data:               map[string]interface{}{"group": "61fbd273-b941-471c-983a-0a3cd2c74747", "target": "River"},
			expectedStatusCode: http.StatusConflict,
			expectedResponse:   gin.H{"err": "user is already a member of group"},
		},
		{
			desc:               "invitealreadyindatabase",
			id:                 "1c4dccaf-a341-4920-9003-f24e0412f8e0",
			data:               map[string]interface{}{"group": "61fbd273-b941-471c-983a-0a3cd2c74747", "target": "John"},
			expectedStatusCode: http.StatusConflict,
			expectedResponse:   gin.H{"err": "user already invited"},
		},
		{
			desc:               "invitenogroup",
			id:                 "1c4dccaf-a341-4920-9003-f24e0412f8e0",
			data:               map[string]interface{}{"target": "Kel"},
			expectedStatusCode: http.StatusBadRequest,
			expectedResponse:   gin.H{"err": "invalid group ID"},
		},
		{
			desc:               "invitenouser",
			id:                 "1c4dccaf-a341-4920-9003-f24e0412f8e0",
			data:               map[string]interface{}{"group": "61fbd273-b941-471c-983a-0a3cd2c74747"},
			expectedStatusCode: http.StatusBadRequest,
			expectedResponse:   gin.H{"err": "user not specified"},
		},
	}

	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {

			requestBody, _ := json.Marshal(tC.data)
			req, _ := http.NewRequest("POST", "/api/invite", bytes.NewReader(requestBody))

			w := httptest.NewRecorder()
			_, engine := gin.CreateTestContext(w)
			engine.Use(func(c *gin.Context) {
				c.Set("userID", tC.id)
			})

			engine.Handle(http.MethodPost, "/api/invite", s.SendGroupInvite)
			engine.ServeHTTP(w, req)
			response := w.Result()

			if response.StatusCode != tC.expectedStatusCode {
				t.Errorf("Received Status code %d does not match expected status %d", response.StatusCode, tC.expectedStatusCode)
			}

			var msg gin.H
			json.NewDecoder(response.Body).Decode(&msg)

			if !reflect.DeepEqual(msg, tC.expectedResponse) {
				t.Errorf("Received HTTP response body %+v does not match expected HTTP response Body %+v", msg, tC.expectedResponse)
			}
		})
	}
}

func TestGetUserInvites(t *testing.T) {
	gin.SetMode(gin.TestMode)
	s := setupTestServer()

	dateCreated, _ := time.Parse("2006-01-02T15:04:05Z", "2019-03-17T22:04:45Z")
	dateModified, _ := time.Parse("2006-01-02T15:04:05Z", "2019-03-17T22:04:45Z")
	issId, _ := uuid.Parse("1c4dccaf-a341-4920-9003-f24e0412f8e0")
	groupId, _ := uuid.Parse("61fbd273-b941-471c-983a-0a3cd2c74747")
	targetId, _ := uuid.Parse("634240cf-1219-4be2-adfa-90ab6b47899b")
	inviteId, _ := uuid.Parse("0916b355-323c-45fd-b79e-4160eaac1320")

	testCases := []struct {
		desc               string
		id                 string
		expectedStatusCode int
		expectedResponse   []models.Invite
	}{
		{
			desc:               "getinvitessuccess",
			id:                 "634240cf-1219-4be2-adfa-90ab6b47899b",
			expectedStatusCode: http.StatusOK,
			expectedResponse:   []models.Invite{{ID: inviteId, IssId: issId, TargetID: targetId, GroupID: groupId, Status: 1, Created: dateCreated, Modified: dateModified}},
		},
		{
			desc:               "getinvitesnocontent",
			id:                 "1c4dccaf-a341-4920-9003-f24e0412f8e0",
			expectedStatusCode: http.StatusNoContent,
			expectedResponse:   []models.Invite{},
		},
	}

	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {

			req, _ := http.NewRequest("GET", "/api/invites", nil)

			w := httptest.NewRecorder()
			_, engine := gin.CreateTestContext(w)
			engine.Use(func(c *gin.Context) {
				c.Set("userID", tC.id)
			})

			engine.Handle(http.MethodGet, "/api/invites", s.GetUserInvites)
			engine.ServeHTTP(w, req)
			response := w.Result()

			if response.StatusCode != tC.expectedStatusCode {
				t.Errorf("Received Status code %d does not match expected status %d", response.StatusCode, tC.expectedStatusCode)
			}

			respBody := []models.Invite{}
			json.NewDecoder(response.Body).Decode(&respBody)

			if !reflect.DeepEqual(respBody, tC.expectedResponse) {
				t.Errorf("Received HTTP response body %+v does not match expected HTTP response Body %+v", respBody, tC.expectedResponse)
			}
		})
	}
}

func TestRespondGroupInvite(t *testing.T) {
	gin.SetMode(gin.TestMode)
	s := setupTestServer()

	dateGroupCreated, _ := time.Parse("2006-01-02T15:04:05Z", "2019-01-13T08:47:44Z")
	groupId, _ := uuid.Parse("61fbd273-b941-471c-983a-0a3cd2c74747")

	testCases := []struct {
		desc               string
		userID             string
		answer             bool
		data               map[string]interface{}
		inviteID           string
		returnVal          bool
		expectedStatusCode int
		expectedResponse   interface{}
	}{
		{
			desc:               "respondInviteYes",
			userID:             "634240cf-1219-4be2-adfa-90ab6b47899b",
			data:               map[string]interface{}{"answer": true},
			inviteID:           "0916b355-323c-45fd-b79e-4160eaac1320",
			returnVal:          true,
			expectedStatusCode: http.StatusOK,
			expectedResponse:   models.Group{ID: groupId, Name: "New Group", Desc: "totally new group", Picture: "16fc5e9d-da47-4923-8475-9f444177990d", Created: dateGroupCreated},
		},
		{
			desc:               "respondInviteNo",
			userID:             "634240cf-1219-4be2-adfa-90ab6b47899b",
			data:               map[string]interface{}{"answer": false},
			inviteID:           "0916b355-323c-45fd-b79e-4160eaac1320",
			returnVal:          false,
			expectedStatusCode: http.StatusOK,
			expectedResponse:   gin.H{"message": "invite declined"},
		},
		{
			desc:               "respondInviteNotInDatabase",
			userID:             "634240cf-1219-4be2-adfa-90ab6b47899b",
			data:               map[string]interface{}{"answer": true},
			inviteID:           "0916b355-323c-45fd-b79e-4161eaac1320",
			returnVal:          false,
			expectedStatusCode: http.StatusNotFound,
			expectedResponse:   gin.H{"err": "resource not found"},
		},
		{
			desc:               "respondInviteWrongUser",
			userID:             "1c4dccaf-a341-4920-9003-f24e0412f8e0",
			data:               map[string]interface{}{"answer": true},
			inviteID:           "0916b355-323c-45fd-b79e-4160eaac1320",
			returnVal:          false,
			expectedStatusCode: http.StatusForbidden,
			expectedResponse:   gin.H{"err": "no rights to respond"},
		},
		{
			desc:               "respondInviteNoAnswer",
			userID:             "1c4dccaf-a341-4920-9003-f24e0412f8e0",
			inviteID:           "0916b355-323c-45fd-b79e-4160eaac1320",
			data:               map[string]interface{}{},
			returnVal:          false,
			expectedStatusCode: http.StatusBadRequest,
			expectedResponse:   gin.H{"err": "answer not specified"},
		},
	}

	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {

			requestBody, _ := json.Marshal(tC.data)
			req, _ := http.NewRequest("PUT", "/api/invite/"+tC.inviteID, bytes.NewReader(requestBody))

			w := httptest.NewRecorder()
			_, engine := gin.CreateTestContext(w)
			engine.Use(func(c *gin.Context) {
				c.Set("userID", tC.userID)
			})

			engine.Handle(http.MethodPut, "/api/invite/:inviteID", s.RespondGroupInvite)
			engine.ServeHTTP(w, req)
			response := w.Result()

			if response.StatusCode != tC.expectedStatusCode {
				t.Errorf("Received Status code %d does not match expected status %d", response.StatusCode, tC.expectedStatusCode)
			}

			var respBody interface{}
			if tC.returnVal {
				group := models.Group{}
				json.NewDecoder(response.Body).Decode(&group)
				respBody = group
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
