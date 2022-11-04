package mock

import (
	"encoding/json"

	"github.com/Slimo300/MicroservicesChatApp/backend/group-service/models"
)

type MockDB struct {
	Users   []models.User
	Groups  []models.Group
	Members []models.Member
	Invites []models.Invite
}

func NewMockDB() *MockDB {

	USERS := `[
		{
			"ID": "1c4dccaf-a341-4920-9003-f24e0412f8e0",
			"signup": "2018-08-14T07:52:54Z",
			"active": "2019-01-13T22:00:45Z",
			"username": "Mal",
			"email": "mal.zein@email.com",
			"pictureUrl": "16fc5e9d-da47-4923-8475-9f444177990d",
			"password": "$2a$10$6BSuuiaPdRJJF2AygYAfnOGkrKLY2o0wDWbEpebn.9Rk0O95D3hW.",
			"logged": true
		},
		{
			"ID": "0ef41409-24b0-43e6-80a3-cf31a4b1a684",
			"signup": "2018-08-14T07:52:55Z",
			"active": "2019-01-12T22:39:01Z",
			"username": "River",
			"email": "river.sam@email.com",
			"password": "$2a$10$BvQjoN.PH8FkVPV3ZMBK1O.3LGhrF/RhZ2kM/h9M3jPA1d2lZkL.C",
			"logged": false
		},
		{
			"ID": "634240cf-1219-4be2-adfa-90ab6b47899b",
			"username": "John",
			"signup": "2019-01-13T08:43:44Z",
			"active": "2019-01-13T15:12:25Z",
			"email": "john.doe@bla.com",
			"password": "$2a$10$T4c8rmpbgKrUA0sIqtHCaO0g2XGWWxFY4IGWkkpVQOD/iuBrwKrZu",
			"logged": false
		},
		{
			"ID": "1fa00013-89b1-49ad-af29-a79afea1f8b8",
			"username": "Kal",
			"signup": "2019-01-13T08:53:44Z",
			"active": "2019-01-13T15:52:25Z",
			"email": "kal.doe@bla.com",
			"password": "$2a$10$T4c8rmpbgKrUA0sIqtHCaO0g2XGWWxFY4IGWkkpVQOD/iuBrwKrZu",
			"logged": false
		},
		{
			"ID": "fe98176b-cf09-4da5-94ae-81207519a75f",
			"username": "Kel",
			"signup": "2019-01-12T08:53:44Z",
			"active": "2019-01-12T15:52:25Z",
			"email": "kel.doa@bla.com",
			"password": "$2a$10$T4c8rmpbgKrUA0sIqtHCaO0g2XGWWxFY4IGWkkpVQOD/iuBrwKrZu",
			"logged": false
		}
	]`

	GROUPS := `[
		{
			"ID": "61fbd273-b941-471c-983a-0a3cd2c74747",
			"name": "New Group",
			"pictureUrl": "16fc5e9d-da47-4923-8475-9f444177990d",
			"desc": "totally new group",
			"created": "2019-01-13T08:47:44Z"
		},
		{
			"ID": "87a0c639-e590-422e-9664-6aedd5ef85ba",
			"name": "New Group2",
			"desc": "totally new group2",
			"created": "2019-01-13T08:47:45Z"
		}
	]`

	MEMBERS := `[
		{
			"ID": "3208e6cc-858e-4ca0-b03b-e500cd335290",
			"group_id": "87a0c639-e590-422e-9664-6aedd5ef85ba",
			"user_id": "1c4dccaf-a341-4920-9003-f24e0412f8e0",
			"nick": "Mal",
			"adding": true,
			"deleting": true,
			"setting": true,
			"creator": true,
			"deleted": false
		},
		{
			"ID": "e4372b71-30ca-42e1-8c1e-7df6d033fd3f",
			"group_id": "61fbd273-b941-471c-983a-0a3cd2c74747",
			"user_id": "1c4dccaf-a341-4920-9003-f24e0412f8e0",
			"nick": "Mal",
			"adding": true,
			"deleting": true,
			"setting": true,
			"creator": true,
			"deleted": false
		},
		{
			"ID": "b38aaff8-6733-4a1d-8eaf-fc10e656d02b",
			"group_id": "61fbd273-b941-471c-983a-0a3cd2c74747",
			"user_id": "0ef41409-24b0-43e6-80a3-cf31a4b1a684",
			"nick": "River",
			"adding": false,
			"deleting": false,
			"setting": false,
			"creator": false,
			"deleted": false
		},
		{
			"ID": "324c0237-e196-4c96-a3be-fee55e745b89",
			"group_id": "61fbd273-b941-471c-983a-0a3cd2c74747",
			"user_id": "1fa00013-89b1-49ad-af29-a79afea1f8b8",
			"nick": "Kal",
			"adding": false,
			"deleting": false,
			"setting": false,
			"creator": false,
			"deleted": true
		}
	]`

	INVITES := `[
		{
			"ID": "0916b355-323c-45fd-b79e-4160eaac1320",
			"issID": "1c4dccaf-a341-4920-9003-f24e0412f8e0",
			"targetID": "634240cf-1219-4be2-adfa-90ab6b47899b",
			"groupID": "61fbd273-b941-471c-983a-0a3cd2c74747",
			"status": 1,
			"created": "2019-03-17T22:04:45Z",
			"modified": "2019-03-17T22:04:45Z"
		}
	]`

	var users []models.User
	json.Unmarshal([]byte(USERS), &users)

	var groups []models.Group
	json.Unmarshal([]byte(GROUPS), &groups)

	var members []models.Member
	json.Unmarshal([]byte(MEMBERS), &members)

	var invites []models.Invite
	json.Unmarshal([]byte(INVITES), &invites)

	// add data
	return &MockDB{Users: users, Groups: groups, Members: members, Invites: invites}
}
