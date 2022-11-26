package handlers_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Slimo300/MicroservicesChatApp/backend/group-service/database"
	mockdb "github.com/Slimo300/MicroservicesChatApp/backend/group-service/database/mock"
	"github.com/Slimo300/MicroservicesChatApp/backend/group-service/handlers"
	"github.com/Slimo300/MicroservicesChatApp/backend/group-service/models"
	mockqueue "github.com/Slimo300/MicroservicesChatApp/backend/lib/msgqueue/mock"
	"github.com/Slimo300/MicroservicesChatApp/backend/lib/storage"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/stretchr/testify/suite"
)

type GroupTestSuite struct {
	suite.Suite
	uuids   map[string]uuid.UUID
	db      database.DBlayer
	storage *storage.MockStorage
	emitter *mockqueue.MockEmitter
}

func (s *GroupTestSuite) SetupSuite() {

	s.uuids = make(map[string]uuid.UUID)

	s.uuids["user1"] = uuid.MustParse("d562da33-062b-4f18-afa0-f9fd1b7aadd3")
	s.uuids["user2"] = uuid.MustParse("b86dcd49-2fb3-48b1-bee7-ee46d4682244")
	s.uuids["group1"] = uuid.MustParse("e7564f8c-8917-4527-a020-1db4b901d4b9")
	s.uuids["group2"] = uuid.MustParse("be0accb1-50db-4698-b048-fb0128e35684")

	db := new(mockdb.MockGroupsDB)
	db.On("GetUserGroups", s.uuids["user1"]).Return([]models.Group{
		{ID: s.uuids["group1"]},
		{ID: s.uuids["group2"]},
	}, nil)
	db.On("GetUserGroups", s.uuids["user2"]).Return([]models.Group{}, nil)
	s.db = db

	// Handlers don't handle emitter errors so there is no need to mock one
	// s.emitter = new(mockqueue.MockEmitter)
	// s.emitter.On("Emit", mock.Anything).Return(nil)

	s.storage = new(storage.MockStorage)
}

func (s GroupTestSuite) TestGetUserGroups() {
	gin.SetMode(gin.TestMode)

	server := handlers.NewServer(s.db, s.storage, nil)

	testCases := []struct {
		desc               string
		userID             string
		returnVal          bool
		expectedStatusCode int
		expectedResponse   interface{}
	}{
		{
			desc:               "getgroupssuccess",
			userID:             s.uuids["user1"].String(),
			returnVal:          true,
			expectedStatusCode: http.StatusOK,
			expectedResponse: []models.Group{
				{ID: s.uuids["group1"]},
				{ID: s.uuids["group2"]},
			},
		},
		{
			desc:               "getgroupsnone",
			userID:             s.uuids["user2"].String(),
			returnVal:          false,
			expectedStatusCode: http.StatusNoContent,
			expectedResponse:   nil,
		},
	}

	for _, tC := range testCases {
		s.Run(tC.desc, func() {

			req, _ := http.NewRequest("GET", "/api/group/get", nil)

			w := httptest.NewRecorder()
			_, engine := gin.CreateTestContext(w)

			engine.Use(func(c *gin.Context) {
				c.Set("userID", tC.userID)
			})
			engine.Handle(http.MethodGet, "/api/group/get", server.GetUserGroups)
			engine.ServeHTTP(w, req)
			response := w.Result()

			s.Equal(response.StatusCode, tC.expectedStatusCode)

			if tC.returnVal {
				respBody := []models.Group{}
				json.NewDecoder(response.Body).Decode(&respBody)

				s.Equal(respBody, tC.expectedResponse)
			}
		})
	}
}

// func (s GroupTestSuite) TestCreateGroup()

func TestGroupSuite(t *testing.T) {
	suite.Run(t, &GroupTestSuite{})
}
