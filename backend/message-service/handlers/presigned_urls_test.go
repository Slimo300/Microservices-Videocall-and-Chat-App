package handlers_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"

	mockdb "github.com/Slimo300/Microservices-Videocall-and-Chat-App/backend/message-service/database/mock"
	"github.com/Slimo300/Microservices-Videocall-and-Chat-App/backend/message-service/handlers"
	"github.com/Slimo300/Microservices-Videocall-and-Chat-App/backend/message-service/models"
	"github.com/Slimo300/Microservices-Videocall-and-Chat-App/backend/message-service/storage"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

type PresignedUrlSuite struct {
	suite.Suite
	uuids  map[string]uuid.UUID
	server handlers.Server
}

func (s *PresignedUrlSuite) SetupSuite() {
	s.uuids = make(map[string]uuid.UUID)
	s.uuids["userInGroup"] = uuid.MustParse("1c4dccaf-a341-4920-9003-f24e0412f8e0")
	s.uuids["userNotInGroup"] = uuid.MustParse("634240cf-1219-4be2-adfa-90ab6b47899b")
	s.uuids["groupID"] = uuid.MustParse("61fbd273-b941-471c-983a-0a3cd2c74747")
	s.uuids["memberID"] = uuid.MustParse("cf003fcf-47c4-497b-bc8d-b2f5df481979")

	mockDB := new(mockdb.MockMessageDB)
	mockDB.On("GetGroupMembership", s.uuids["userNotInGroup"], s.uuids["groupID"]).Return(models.Membership{}, errors.New("no membership"))
	mockDB.On("GetGroupMembership", s.uuids["userInGroup"], s.uuids["groupID"]).Return(models.Membership{MembershipID: s.uuids["memberID"]}, nil)

	mockStorage := new(storage.MockStorage)
	mockStorage.On("GetPresignedPutRequests", mock.AnythingOfType("string"), mock.Anything, mock.Anything).Return([]storage.FileOutput{{
		Name:         "duckie.jpeg",
		PresignedURL: "someUrl",
	}, {
		Name:         "kitty.jpeg",
		PresignedURL: "someUrl",
	}}, nil)

	s.server = *handlers.NewServer(
		mockDB,
		nil,
		nil,
		mockStorage,
	)
	log.Println(s.server)

}

func (s *PresignedUrlSuite) TestGetPresignedUrl() {
	gin.SetMode(gin.TestMode)

	testCases := []struct {
		desc               string
		userID             string
		groupID            string
		data               map[string]interface{}
		returnVal          bool
		expectedStatusCode int
		expectedResponse   interface{}
	}{
		{
			desc:               "invalidUserID",
			userID:             "1",
			returnVal:          false,
			expectedStatusCode: http.StatusBadRequest,
			expectedResponse:   gin.H{"err": "invalid user ID"},
		},
		{
			desc:               "invalidGroupID",
			userID:             s.uuids["userInGroup"].String(),
			groupID:            "1",
			returnVal:          false,
			expectedStatusCode: http.StatusBadRequest,
			expectedResponse:   gin.H{"err": "invalid group ID"},
		},
		{
			desc:               "invalidRequestBody",
			userID:             s.uuids["userInGroup"].String(),
			groupID:            s.uuids["groupID"].String(),
			data:               map[string]interface{}{},
			returnVal:          false,
			expectedStatusCode: http.StatusBadRequest,
			expectedResponse:   gin.H{"err": "invalid request body"},
		},
		{
			desc:    "userNotInGroup",
			userID:  s.uuids["userNotInGroup"].String(),
			groupID: s.uuids["groupID"].String(),
			data: map[string]interface{}{
				"files": []map[string]interface{}{
					{
						"name": "duckie.jpeg",
						"size": 1024,
					},
				},
			},
			returnVal:          false,
			expectedStatusCode: http.StatusForbidden,
			expectedResponse:   gin.H{"err": "user cannot send messages to this group"},
		},
		{
			desc:    "success",
			userID:  s.uuids["userInGroup"].String(),
			groupID: s.uuids["groupID"].String(),
			data: map[string]interface{}{
				"files": []storage.FileInput{
					{
						Name: "duckie.jpeg",
						Size: 1024,
					},
					{
						Name: "kitty.jpeg",
						Size: 1024,
					},
				},
			},
			returnVal:          true,
			expectedStatusCode: http.StatusOK,
			expectedResponse: []storage.FileOutput{{
				Name:         "duckie.jpeg",
				PresignedURL: "someUrl",
			}, {
				Name:         "kitty.jpeg",
				PresignedURL: "someUrl",
			}},
		},
	}

	for _, tC := range testCases {
		s.Run(tC.desc, func() {

			requestBody, _ := json.Marshal(tC.data)

			r, _ := http.NewRequest(http.MethodPost, "/group/"+tC.groupID+"/uploads", bytes.NewBuffer(requestBody))
			w := httptest.NewRecorder()
			_, engine := gin.CreateTestContext(w)

			engine.Use(func(ctx *gin.Context) {
				ctx.Set("userID", tC.userID)
			})

			engine.POST("/group/:groupID/uploads", s.server.GetPresignedPutRequest)
			engine.ServeHTTP(w, r)

			response := w.Result()
			defer response.Body.Close()

			s.Equal(tC.expectedStatusCode, response.StatusCode)

			var respBody interface{}
			if tC.returnVal {
				var msg []storage.FileOutput
				if err := json.NewDecoder(response.Body).Decode(&msg); err != nil {
					s.Fail(err.Error())
				}

				respBody = msg
			} else {
				var msg gin.H
				if err := json.NewDecoder(response.Body).Decode(&msg); err != nil {
					s.Fail(err.Error())
				}
				respBody = msg
			}

			s.Equal(tC.expectedResponse, respBody)
		})
	}
}

func TestPresignedUrlSuite(t *testing.T) {
	suite.Run(t, &PresignedUrlSuite{})
}
