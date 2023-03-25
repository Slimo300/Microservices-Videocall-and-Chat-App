package handlers_test

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Slimo300/MicroservicesChatApp/backend/lib/storage"
	mockdb "github.com/Slimo300/chat-messageservice/internal/database/mock"
	"github.com/Slimo300/chat-messageservice/internal/handlers"
	"github.com/Slimo300/chat-messageservice/internal/models"
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
	mockStorage.On("GetPresignedPutRequest", mock.AnythingOfType("string")).Return("someUrl", nil)

	s.server = *handlers.NewServer(
		mockDB,
		nil,
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
		num                string
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
			desc:               "invalidNumQuery",
			userID:             s.uuids["userInGroup"].String(),
			groupID:            s.uuids["groupID"].String(),
			num:                "a",
			returnVal:          false,
			expectedStatusCode: http.StatusBadRequest,
			expectedResponse:   gin.H{"err": "files query value is not a valid integer"},
		},
		{
			desc:               "userNotInGroup",
			userID:             s.uuids["userNotInGroup"].String(),
			groupID:            s.uuids["groupID"].String(),
			num:                "",
			returnVal:          false,
			expectedStatusCode: http.StatusForbidden,
			expectedResponse:   gin.H{"err": "user cannot send messages to this group"},
		},
		{
			desc:               "success",
			userID:             s.uuids["userInGroup"].String(),
			groupID:            s.uuids["groupID"].String(),
			num:                "2",
			returnVal:          true,
			expectedStatusCode: http.StatusOK,
			expectedResponse:   []string{"someUrl", "someUrl"},
		},
	}

	for _, tC := range testCases {
		s.Run(tC.desc, func() {

			r, _ := http.NewRequest(http.MethodGet, "/group/"+tC.groupID+"/uploads?files="+tC.num, nil)
			w := httptest.NewRecorder()
			_, engine := gin.CreateTestContext(w)

			engine.Use(func(ctx *gin.Context) {
				ctx.Set("userID", tC.userID)
			})

			engine.GET("/group/:groupID/uploads", s.server.GetPresignedPutRequest)
			engine.ServeHTTP(w, r)

			response := w.Result()
			defer response.Body.Close()

			s.Equal(tC.expectedStatusCode, response.StatusCode)

			var respBody interface{}
			if tC.returnVal {
				var msg gin.H
				if err := json.NewDecoder(response.Body).Decode(&msg); err != nil {
					s.Fail(err.Error())
				}

				var urls []string
				for _, url := range msg["requests"].([]interface{}) {
					urls = append(urls, url.(map[string]interface{})["url"].(string))
				}

				respBody = urls
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
