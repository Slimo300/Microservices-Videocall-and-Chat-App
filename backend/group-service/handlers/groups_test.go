package handlers_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
	"time"

	"github.com/Slimo300/MicrosevicesChatApp/backend/group-service/models"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func TestGetUserGroups(t *testing.T) {
	gin.SetMode(gin.TestMode)
	s := setupTestServer()

	date1, _ := time.Parse("2006-01-02T15:04:05Z", "2019-01-13T08:47:44Z")
	groupID1, _ := uuid.Parse("61fbd273-b941-471c-983a-0a3cd2c74747")

	date2, _ := time.Parse("2006-01-02T15:04:05Z", "2019-01-13T08:47:45Z")
	groupID2, _ := uuid.Parse("87a0c639-e590-422e-9664-6aedd5ef85ba")

	testCases := []struct {
		desc               string
		id                 string
		returnVal          bool
		expectedStatusCode int
		expectedResponse   interface{}
	}{
		{
			desc:               "getgroupssuccess",
			id:                 "1c4dccaf-a341-4920-9003-f24e0412f8e0",
			returnVal:          true,
			expectedStatusCode: http.StatusOK,
			expectedResponse: []models.Group{
				{ID: groupID2, Name: "New Group2", Desc: "totally new group2", Created: date2},
				{ID: groupID1, Name: "New Group", Desc: "totally new group", Picture: "16fc5e9d-da47-4923-8475-9f444177990d", Created: date1},
			},
		},
		{
			desc:               "getgroupsnone",
			id:                 "634240cf-1219-4be2-adfa-90ab6b47899b",
			returnVal:          false,
			expectedStatusCode: http.StatusNoContent,
			expectedResponse:   nil,
		},
	}

	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {

			req, _ := http.NewRequest("GET", "/api/group/get", nil)

			w := httptest.NewRecorder()
			_, engine := gin.CreateTestContext(w)

			engine.Use(func(c *gin.Context) {
				c.Set("userID", tC.id)
			})
			engine.Handle(http.MethodGet, "/api/group/get", s.GetUserGroups)
			engine.ServeHTTP(w, req)
			response := w.Result()

			if response.StatusCode != tC.expectedStatusCode {
				t.Errorf("Received Status code %d does not match expected status %d", response.StatusCode, tC.expectedStatusCode)
			}
			if tC.returnVal {
				respBody := []models.Group{}
				json.NewDecoder(response.Body).Decode(&respBody)

				if !reflect.DeepEqual(respBody, tC.expectedResponse) {
					t.Errorf("Received HTTP response body %+v does not match expected HTTP response Body %+v", respBody, tC.expectedResponse)
				}
			}
		})
	}
}

func TestDeleteGroup(t *testing.T) {
	gin.SetMode(gin.TestMode)
	s := setupTestServerWithHub()

	testCases := []struct {
		desc               string
		userID             string
		groupID            string
		expectedStatusCode int
		expectedResponse   interface{}
	}{
		// user is not a creator of the group so he can't delete it
		{
			desc:               "deletegroupnosuccess",
			userID:             "634240cf-1219-4be2-adfa-90ab6b47899b",
			groupID:            "61fbd273-b941-471c-983a-0a3cd2c74747",
			expectedStatusCode: http.StatusForbidden,
			expectedResponse:   gin.H{"err": "couldn't delete group"},
		},
		// user hasn't specified a group in a query
		{
			desc:               "deletegroupnotspecified",
			userID:             "1c4dccaf-a341-4920-9003-f24e0412f8e0",
			groupID:            "sa",
			expectedStatusCode: http.StatusBadRequest,
			expectedResponse:   gin.H{"err": "invalid group ID"},
		},
		// creator deletes a group
		{
			desc:               "deletegroupsuccess",
			userID:             "1c4dccaf-a341-4920-9003-f24e0412f8e0",
			groupID:            "61fbd273-b941-471c-983a-0a3cd2c74747",
			expectedStatusCode: http.StatusOK,
			expectedResponse:   gin.H{"message": "ok"},
		},
	}

	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			req, _ := http.NewRequest(http.MethodDelete, "/api/group/"+tC.groupID, nil)

			w := httptest.NewRecorder()
			_, engine := gin.CreateTestContext(w)

			engine.Use(func(c *gin.Context) {
				c.Set("userID", tC.userID)
			})
			engine.Handle(http.MethodDelete, "/api/group/:groupID", s.DeleteGroup)
			engine.ServeHTTP(w, req)
			response := w.Result()

			if response.StatusCode != tC.expectedStatusCode {
				t.Errorf("Received Status code %d does not match expected status %d", response.StatusCode, tC.expectedStatusCode)
			}

			var respBody gin.H
			json.NewDecoder(response.Body).Decode(&respBody)

			if !reflect.DeepEqual(respBody, tC.expectedResponse) {
				t.Errorf("Received HTTP response body %+v does not match expected HTTP response Body %+v", respBody, tC.expectedResponse)
			}
		})
	}
}

func TestCreateGroup(t *testing.T) {
	gin.SetMode(gin.TestMode)
	s := setupTestServerWithHub()

	testCases := []struct {
		desc               string
		ID                 string
		data               map[string]interface{}
		returnVal          bool
		expectedStatusCode int
		expectedResponse   interface{}
	}{
		// no name provided in request body
		{
			desc:               "creategroupnoname",
			ID:                 "634240cf-1219-4be2-adfa-90ab6b47899b",
			data:               map[string]interface{}{"name": "", "desc": "ng1"},
			returnVal:          false,
			expectedStatusCode: http.StatusBadRequest,
			expectedResponse:   gin.H{"err": "bad name"},
		},
		// no description provided in request body
		{
			desc:               "creategroupnodesc",
			ID:                 "634240cf-1219-4be2-adfa-90ab6b47899b",
			data:               map[string]interface{}{"name": "ng1", "desc": ""},
			returnVal:          false,
			expectedStatusCode: http.StatusBadRequest,
			expectedResponse:   gin.H{"err": "bad description"},
		},
		// creator deletes a group
		{
			desc:               "creategroupsuccess",
			ID:                 "634240cf-1219-4be2-adfa-90ab6b47899b",
			data:               map[string]interface{}{"name": "ng1", "desc": "ng1"},
			returnVal:          true,
			expectedStatusCode: http.StatusCreated,
			expectedResponse:   models.Group{Name: "ng1", Desc: "ng1"},
		},
	}

	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {

			requestBody, _ := json.Marshal(tC.data)
			req, _ := http.NewRequest("POST", "/api/group/create", bytes.NewBuffer(requestBody))

			w := httptest.NewRecorder()
			_, engine := gin.CreateTestContext(w)

			engine.Use(func(c *gin.Context) {
				c.Set("userID", tC.ID)
			})

			engine.Handle(http.MethodPost, "/api/group/create", s.CreateGroup)
			engine.ServeHTTP(w, req)
			response := w.Result()

			if response.StatusCode != tC.expectedStatusCode {
				t.Errorf("Received Status code %d does not match expected status %d", response.StatusCode, tC.expectedStatusCode)
			}

			var respBody interface{}
			if tC.returnVal {
				group := models.Group{}
				json.NewDecoder(response.Body).Decode(&group)
				group.Created = time.Time{}
				group.ID = uuid.UUID{}
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
