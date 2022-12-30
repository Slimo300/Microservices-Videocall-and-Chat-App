package handlers_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	mockdb "github.com/Slimo300/MicroservicesChatApp/backend/group-service/database/mock"
	"github.com/Slimo300/MicroservicesChatApp/backend/group-service/handlers"
	"github.com/Slimo300/MicroservicesChatApp/backend/group-service/models"
	"github.com/Slimo300/MicroservicesChatApp/backend/lib/apperrors"
	mockqueue "github.com/Slimo300/MicroservicesChatApp/backend/lib/msgqueue/mock"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

type MembersTestSuite struct {
	suite.Suite
	IDs    map[string]uuid.UUID
	server *handlers.Server
}

func (s *MembersTestSuite) SetupSuite() {

	s.IDs = make(map[string]uuid.UUID)

	s.IDs["userOK"] = uuid.MustParse("37dc93ba-f1ee-497e-aeaf-07588b9ea674")
	s.IDs["userWithoutRights"] = uuid.MustParse("7eb36855-d6d6-46b5-baf4-c806983d7466")
	s.IDs["groupOK"] = uuid.MustParse("9f4fce3d-26b8-46bb-a466-f13e925296a4")
	s.IDs["memberOK"] = uuid.MustParse("6aeada0c-b6f9-4b96-9797-c05d92259171")
	s.IDs["memberNotFound"] = uuid.MustParse("4eec521e-ed3d-4fd3-953e-986d77eda6ed")
	s.IDs["memberHighRank"] = uuid.MustParse("a59aac0f-f575-4412-818f-21a52f1da02d")

	db := new(mockdb.MockGroupsDB)

	db.On("DeleteMember", s.IDs["userOK"], s.IDs["groupOK"], s.IDs["memberOK"]).Return(&models.Member{ID: s.IDs["memberOK"]}, nil)
	db.On("DeleteMember", s.IDs["userWithoutRights"], s.IDs["groupOK"], s.IDs["memberOK"]).
		Return(nil, apperrors.NewForbidden(fmt.Sprintf("User %v has no right to delete members in group %v", s.IDs["userWithoutRights"], s.IDs["groupOK"])))
	db.On("DeleteMember", s.IDs["userOK"], s.IDs["groupOK"], s.IDs["memberNotFound"]).
		Return(nil, apperrors.NewNotFound("member", s.IDs["memberNotFound"].String()))
	db.On("DeleteMember", s.IDs["userOK"], s.IDs["groupOK"], s.IDs["memberHighRank"]).
		Return(nil, apperrors.NewForbidden(fmt.Sprintf("User %v cannot delete member %v", s.IDs["userOK"], s.IDs["memberHighRank"])))

	db.On("GrantRights", s.IDs["userOK"], s.IDs["groupOK"], s.IDs["memberOK"], mock.Anything).Return(nil, nil)
	db.On("GrantRights", s.IDs["userWithoutRights"], s.IDs["groupOK"], s.IDs["memberOK"], mock.Anything).
		Return(nil, apperrors.NewForbidden(fmt.Sprintf("User %v has no right to alter members in group %v", s.IDs["userWithoutRights"], s.IDs["groupOK"])))
	db.On("GrantRights", s.IDs["userOK"], s.IDs["groupOK"], s.IDs["memberNotFound"], mock.Anything).
		Return(nil, apperrors.NewNotFound("member", s.IDs["memberNotFound"].String()))
	db.On("GrantRights", s.IDs["userOK"], s.IDs["groupOK"], s.IDs["memberHighRank"], mock.Anything).
		Return(nil, apperrors.NewForbidden(fmt.Sprintf("User %v cannot alter member %v", s.IDs["userOK"], s.IDs["memberHighRank"])))

	db.On("DeleteGroup", s.IDs["userWithoutRights"], s.IDs["groupOK"]).
		Return(models.Group{}, apperrors.NewForbidden("User has no right to delete group"))
	db.On("DeleteGroup", s.IDs["userOK"], s.IDs["groupOK"]).
		Return(models.Group{ID: s.IDs["groupOK"]}, nil)

	emitter := new(mockqueue.MockEmitter)
	emitter.On("Emit", mock.Anything).Return(nil)

	s.server = handlers.NewServer(db, nil, nil)
	s.server.Emitter = emitter
}

func (s MembersTestSuite) TestGrantPriv() {
	gin.SetMode(gin.TestMode)

	testCases := []struct {
		desc               string
		userID             string
		groupID            string
		memberID           string
		data               map[string]interface{}
		expectedStatusCode int
		expectedResponse   interface{}
	}{
		{
			desc:               "UpdateRightsBadUserID",
			userID:             s.IDs["userOK"].String()[:2],
			groupID:            s.IDs["groupOK"].String(),
			memberID:           s.IDs["memberOK"].String(),
			data:               map[string]interface{}{"adding": -1},
			expectedStatusCode: http.StatusBadRequest,
			expectedResponse:   gin.H{"err": "invalid ID"},
		},
		{
			desc:               "UpdateRightsBadGroupID",
			userID:             s.IDs["userOK"].String(),
			groupID:            s.IDs["groupOK"].String()[:2],
			memberID:           s.IDs["memberOK"].String(),
			data:               map[string]interface{}{"adding": -1},
			expectedStatusCode: http.StatusBadRequest,
			expectedResponse:   gin.H{"err": "invalid group ID"},
		},
		{
			desc:               "UpdateRightsBadMemberID",
			userID:             s.IDs["userOK"].String(),
			groupID:            s.IDs["groupOK"].String(),
			memberID:           s.IDs["memberOK"].String()[:2],
			data:               map[string]interface{}{"adding": -1},
			expectedStatusCode: http.StatusBadRequest,
			expectedResponse:   gin.H{"err": "invalid member ID"},
		},
		{
			desc:               "UpdateRightsNoAction",
			userID:             s.IDs["userOK"].String(),
			groupID:            s.IDs["groupOK"].String(),
			memberID:           s.IDs["memberOK"].String(),
			data:               map[string]interface{}{},
			expectedStatusCode: http.StatusBadRequest,
			expectedResponse:   gin.H{"err": "no action specified"},
		},
		{
			desc:               "UpdateRightsNoRights",
			userID:             s.IDs["userWithoutRights"].String(),
			groupID:            s.IDs["groupOK"].String(),
			memberID:           s.IDs["memberOK"].String(),
			data:               map[string]interface{}{"adding": -1},
			expectedStatusCode: http.StatusForbidden,
			expectedResponse:   gin.H{"err": fmt.Sprintf("Forbidden action. Reason: User %v has no right to alter members in group %v", s.IDs["userWithoutRights"].String(), s.IDs["groupOK"].String())},
		},
		{
			desc:               "UpdateRightsNotFound",
			userID:             s.IDs["userOK"].String(),
			groupID:            s.IDs["groupOK"].String(),
			memberID:           s.IDs["memberNotFound"].String(),
			data:               map[string]interface{}{"adding": -1},
			expectedStatusCode: http.StatusNotFound,
			expectedResponse:   gin.H{"err": fmt.Sprintf("resource: member with value: %v not found", s.IDs["memberNotFound"].String())},
		},
		{
			desc:               "UpdateRightsHighRank",
			userID:             s.IDs["userOK"].String(),
			groupID:            s.IDs["groupOK"].String(),
			memberID:           s.IDs["memberHighRank"].String(),
			data:               map[string]interface{}{"adding": -1},
			expectedStatusCode: http.StatusForbidden,
			expectedResponse:   gin.H{"err": fmt.Sprintf("Forbidden action. Reason: User %v cannot alter member %v", s.IDs["userOK"].String(), s.IDs["memberHighRank"].String())},
		},
		{
			desc:               "UpdateRightsSuccess",
			userID:             s.IDs["userOK"].String(),
			groupID:            s.IDs["groupOK"].String(),
			memberID:           s.IDs["memberOK"].String(),
			data:               map[string]interface{}{"adding": -1},
			expectedStatusCode: http.StatusOK,
			expectedResponse:   gin.H{"message": "member updated"},
		},
	}

	for _, tC := range testCases {
		s.Run(tC.desc, func() {

			requestBody, _ := json.Marshal(tC.data)
			req, _ := http.NewRequest(http.MethodPut, "/group/"+tC.groupID+"/member/"+tC.memberID, bytes.NewBuffer(requestBody))

			w := httptest.NewRecorder()
			_, engine := gin.CreateTestContext(w)
			engine.Use(func(c *gin.Context) {
				c.Set("userID", tC.userID)
			})

			engine.Handle(http.MethodPut, "/group/:groupID/member/:memberID", s.server.GrantPriv)
			engine.ServeHTTP(w, req)
			response := w.Result()
			defer response.Body.Close()

			s.Equal(tC.expectedStatusCode, response.StatusCode)

			var msg gin.H
			json.NewDecoder(response.Body).Decode(&msg)

			s.Equal(tC.expectedResponse, msg)
		})
	}
}

func (s MembersTestSuite) TestDeleteMember() {
	gin.SetMode(gin.TestMode)

	testCases := []struct {
		desc               string
		userID             string
		groupID            string
		memberID           string
		expectedStatusCode int
		expectedResponse   interface{}
	}{
		{
			desc:               "DeleteMemberBadUserID",
			userID:             s.IDs["userOK"].String()[:2],
			groupID:            s.IDs["groupOK"].String(),
			memberID:           s.IDs["memberOK"].String(),
			expectedStatusCode: http.StatusBadRequest,
			expectedResponse:   gin.H{"err": "invalid ID"},
		},
		{
			desc:               "DeleteMemberBadGroupID",
			userID:             s.IDs["userOK"].String(),
			groupID:            s.IDs["groupOK"].String()[:2],
			memberID:           s.IDs["memberOK"].String(),
			expectedStatusCode: http.StatusBadRequest,
			expectedResponse:   gin.H{"err": "invalid group ID"},
		},
		{
			desc:               "DeleteMemberBadMemberID",
			userID:             s.IDs["userOK"].String(),
			groupID:            s.IDs["groupOK"].String(),
			memberID:           s.IDs["memberOK"].String()[:2],
			expectedStatusCode: http.StatusBadRequest,
			expectedResponse:   gin.H{"err": "invalid member ID"},
		},
		{
			desc:               "DeleteMemberNoRights",
			userID:             s.IDs["userWithoutRights"].String(),
			groupID:            s.IDs["groupOK"].String(),
			memberID:           s.IDs["memberOK"].String(),
			expectedStatusCode: http.StatusForbidden,
			expectedResponse:   gin.H{"err": fmt.Sprintf("Forbidden action. Reason: User %v has no right to delete members in group %v", s.IDs["userWithoutRights"].String(), s.IDs["groupOK"].String())},
		},
		{
			desc:               "DeleteMemberNotFound",
			userID:             s.IDs["userOK"].String(),
			groupID:            s.IDs["groupOK"].String(),
			memberID:           s.IDs["memberNotFound"].String(),
			expectedStatusCode: http.StatusNotFound,
			expectedResponse:   gin.H{"err": fmt.Sprintf("resource: member with value: %v not found", s.IDs["memberNotFound"].String())},
		},
		{
			desc:               "DeleteMemberHighRank",
			userID:             s.IDs["userOK"].String(),
			groupID:            s.IDs["groupOK"].String(),
			memberID:           s.IDs["memberHighRank"].String(),
			expectedStatusCode: http.StatusForbidden,
			expectedResponse:   gin.H{"err": fmt.Sprintf("Forbidden action. Reason: User %v cannot delete member %v", s.IDs["userOK"].String(), s.IDs["memberHighRank"].String())},
		},
		{
			desc:               "DeleteMemberSuccess",
			userID:             s.IDs["userOK"].String(),
			groupID:            s.IDs["groupOK"].String(),
			memberID:           s.IDs["memberOK"].String(),
			expectedStatusCode: http.StatusOK,
			expectedResponse:   gin.H{"message": "member deleted"},
		},
	}

	for _, tC := range testCases {
		s.Run(tC.desc, func() {
			req, _ := http.NewRequest(http.MethodDelete, "/group/"+tC.groupID+"/member/"+tC.memberID, nil)

			w := httptest.NewRecorder()
			_, engine := gin.CreateTestContext(w)
			engine.Use(func(c *gin.Context) {
				c.Set("userID", tC.userID)
			})

			engine.Handle(http.MethodDelete, "/group/:groupID/member/:memberID", s.server.DeleteUserFromGroup)
			engine.ServeHTTP(w, req)
			response := w.Result()
			defer response.Body.Close()

			s.Equal(tC.expectedStatusCode, response.StatusCode)

			var msg gin.H
			json.NewDecoder(response.Body).Decode(&msg)

			s.Equal(tC.expectedResponse, msg)
		})
	}
}

func (s MembersTestSuite) TestDeleteGroup() {
	gin.SetMode(gin.TestMode)

	testCases := []struct {
		desc               string
		userID             string
		groupID            string
		expectedStatusCode int
		expectedResponse   interface{}
	}{

		{
			desc:               "DeleteGroupBadUserID",
			userID:             s.IDs["userOK"].String()[:2],
			groupID:            s.IDs["groupOK"].String(),
			expectedStatusCode: http.StatusBadRequest,
			expectedResponse:   gin.H{"err": "invalid ID"},
		},
		{
			desc:               "DeleteGroupBadGroupID",
			userID:             s.IDs["userOK"].String(),
			groupID:            s.IDs["groupOK"].String()[:2],
			expectedStatusCode: http.StatusBadRequest,
			expectedResponse:   gin.H{"err": "invalid group ID"},
		},
		{
			desc:               "DeleteGroupNoRights",
			userID:             s.IDs["userWithoutRights"].String(),
			groupID:            s.IDs["groupOK"].String(),
			expectedStatusCode: http.StatusForbidden,
			expectedResponse:   gin.H{"err": "Forbidden action. Reason: User has no right to delete group"},
		},
		{
			desc:               "DeleteGroupSuccess",
			userID:             s.IDs["userOK"].String(),
			groupID:            s.IDs["groupOK"].String(),
			expectedStatusCode: http.StatusOK,
			expectedResponse:   gin.H{"message": "group deleted"},
		},
	}

	for _, tC := range testCases {
		s.Run(tC.desc, func() {
			req, _ := http.NewRequest(http.MethodDelete, "/api/group/"+tC.groupID, nil)

			w := httptest.NewRecorder()
			_, engine := gin.CreateTestContext(w)

			engine.Use(func(c *gin.Context) {
				c.Set("userID", tC.userID)
			})
			engine.Handle(http.MethodDelete, "/api/group/:groupID", s.server.DeleteGroup)
			engine.ServeHTTP(w, req)
			response := w.Result()
			defer response.Body.Close()

			s.Equal(tC.expectedStatusCode, response.StatusCode)

			var msg gin.H
			json.NewDecoder(response.Body).Decode(&msg)

			s.Equal(tC.expectedResponse, msg)
		})
	}
}

func TestMembers(t *testing.T) {
	suite.Run(t, &MembersTestSuite{})
}
