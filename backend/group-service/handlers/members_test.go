package handlers_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/gin-gonic/gin"
)

func TestGrantPriv(t *testing.T) {
	gin.SetMode(gin.TestMode)
	s := setupTestServer()

	testCases := []struct {
		desc               string
		userID             string
		memberID           string
		data               map[string]interface{}
		expectedStatusCode int
		expectedResponse   interface{}
	}{
		{
			desc:               "grantprivsuccess",
			userID:             "1c4dccaf-a341-4920-9003-f24e0412f8e0",
			memberID:           "b38aaff8-6733-4a1d-8eaf-fc10e656d02b",
			data:               map[string]interface{}{"adding": true, "deleting": true, "setting": false},
			expectedStatusCode: http.StatusOK,
			expectedResponse:   gin.H{"message": "ok"},
		},
		{
			desc:               "grantprivbadrequest",
			userID:             "1c4dccaf-a341-4920-9003-f24e0412f8e0",
			memberID:           "b38aaff8-6733-4a1d-8eaf-fc10e656d02b",
			data:               map[string]interface{}{"adding": true, "deleting": true},
			expectedStatusCode: http.StatusBadRequest,
			expectedResponse:   gin.H{"err": "bad request, all 3 fields must be present"},
		},
		// no member provided in request body
		{
			desc:               "grantprivnomember",
			userID:             "1c4dccaf-a341-4920-9003-f24e0412f8e0",
			memberID:           "1c4dccaf-a341-4920-9003-f24e0412f8e0",
			data:               map[string]interface{}{"adding": true, "deleting": true, "setting": false},
			expectedStatusCode: http.StatusNotFound,
			expectedResponse:   gin.H{"err": "resource not found"},
		},
		// issuer has no right to add
		{
			desc:               "grantprivmemberdeleted",
			userID:             "1c4dccaf-a341-4920-9003-f24e0412f8e0",
			memberID:           "0ef41409-24b0-43e6-80a3-cf31a4b1a684",
			data:               map[string]interface{}{"adding": true, "deleting": true, "setting": false},
			expectedStatusCode: http.StatusNotFound,
			expectedResponse:   gin.H{"err": "resource not found"},
		},
		{
			desc:               "grantprivnorights",
			userID:             "0ef41409-24b0-43e6-80a3-cf31a4b1a684",
			memberID:           "b38aaff8-6733-4a1d-8eaf-fc10e656d02b",
			data:               map[string]interface{}{"adding": true, "deleting": true, "setting": true},
			expectedStatusCode: http.StatusForbidden,
			expectedResponse:   gin.H{"err": "no rights to put"},
		},
		{
			desc:               "grantprivcreator",
			userID:             "1c4dccaf-a341-4920-9003-f24e0412f8e0",
			memberID:           "e4372b71-30ca-42e1-8c1e-7df6d033fd3f",
			data:               map[string]interface{}{"adding": true, "deleting": true, "setting": false},
			expectedStatusCode: http.StatusForbidden,
			expectedResponse:   gin.H{"err": "creator can't be modified"},
		},
	}

	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {

			requestBody, _ := json.Marshal(tC.data)
			req, _ := http.NewRequest(http.MethodPut, "/api/member/"+tC.memberID, bytes.NewBuffer(requestBody))

			w := httptest.NewRecorder()
			_, engine := gin.CreateTestContext(w)
			engine.Use(func(c *gin.Context) {
				c.Set("userID", tC.userID)
			})

			engine.Handle(http.MethodPut, "/api/member/:memberID", s.GrantPriv)
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
func TestDeleteMember(t *testing.T) {
	gin.SetMode(gin.TestMode)
	s := setupTestServerWithHub()

	testCases := []struct {
		desc               string
		userID             string
		data               map[string]interface{}
		memberID           string
		expectedStatusCode int
		expectedResponse   interface{}
	}{
		{
			desc:               "deleteusersuccess",
			userID:             "1c4dccaf-a341-4920-9003-f24e0412f8e0",
			memberID:           "b38aaff8-6733-4a1d-8eaf-fc10e656d02b",
			expectedStatusCode: http.StatusOK,
			expectedResponse:   gin.H{"message": "ok"},
		},
		{
			desc:               "deletebadurl",
			userID:             "1c4dccaf-a341-4920-9003-f24e0412f8e0",
			memberID:           "0ef41409-24bx-43e6-80a3-cf31a4b1a684",
			expectedStatusCode: http.StatusBadRequest,
			expectedResponse:   gin.H{"err": "invalid member ID"},
		},
		{
			desc:               "deletenopriv",
			userID:             "0ef41409-24b0-43e6-80a3-cf31a4b1a684",
			memberID:           "e4372b71-30ca-42e1-8c1e-7df6d033fd3f",
			expectedStatusCode: http.StatusForbidden,
			expectedResponse:   gin.H{"err": "no rights to delete"},
		},
	}

	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			req, _ := http.NewRequest(http.MethodDelete, "/api/member/"+tC.memberID, nil)

			w := httptest.NewRecorder()
			_, engine := gin.CreateTestContext(w)
			engine.Use(func(c *gin.Context) {
				c.Set("userID", tC.userID)
			})

			engine.Handle(http.MethodDelete, "/api/member/:memberID", s.DeleteUserFromGroup)
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
