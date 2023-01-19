package elastic

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/Slimo300/MicroservicesChatApp/backend/lib/events"
	"github.com/Slimo300/MicroservicesChatApp/backend/search-service/models"
	"github.com/elastic/go-elasticsearch/v8/esapi"
	"github.com/google/uuid"
)

func (es *elasticSearchDB) UpdateProfilePicture(ev events.UserPictureModifiedEvent) error {
	data := map[string]map[string]string{
		"doc": {
			"picture": ev.PictureURL,
		},
	}
	dataJSON, err := json.Marshal(data)
	if err != nil {
		return err
	}

	req := esapi.UpdateRequest{
		Index:      INDEX_NAME,
		DocumentID: ev.ID.String(),
		Body:       bytes.NewReader(dataJSON),
		Refresh:    "true",
	}

	res, err := req.Do(context.TODO(), es.client)
	if err != nil {
		return err
	}

	defer res.Body.Close()

	if res.IsError() {
		var respBody map[string]interface{}
		if err := json.NewDecoder(res.Body).Decode(&respBody); err != nil {
			return fmt.Errorf("Error decoding error response body: %v", err)
		}
		return errors.New(respBody["error"].(map[string]interface{})["root_cause"].(map[string]interface{})["reason"].(string))
	}

	return nil
}

func (es *elasticSearchDB) AddUser(user events.UserRegisteredEvent) error {

	data := map[string]string{
		"username": user.Username,
		"picture":  user.PictureURL,
	}
	dataJSON, err := json.Marshal(data)
	if err != nil {
		return err
	}

	req := esapi.IndexRequest{
		Index:      INDEX_NAME,
		DocumentID: user.ID.String(),
		Body:       bytes.NewReader(dataJSON),
		Refresh:    "true",
	}

	res, err := req.Do(context.TODO(), es.client)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	if res.IsError() {
		var respBody map[string]interface{}
		if err := json.NewDecoder(res.Body).Decode(&respBody); err != nil {
			return fmt.Errorf("Error decoding error response body: %v", err)
		}
		return errors.New(respBody["error"].(map[string]interface{})["root_cause"].(map[string]interface{})["reason"].(string))
	}

	return nil
}

func (es *elasticSearchDB) GetUsers(query string, num int) ([]models.User, error) {

	var buffer bytes.Buffer

	reqBody := map[string]interface{}{
		"from": 0,
		"size": num,
		"query": map[string]interface{}{
			"match": map[string]interface{}{
				"username": map[string]interface{}{
					"query":                query,
					"minimum_should_match": MINIMUM_SHOULD_MATCH,
				},
			},
		},
	}

	if err := json.NewEncoder(&buffer).Encode(reqBody); err != nil {
		return nil, fmt.Errorf("Error when encoding query: %v", err)
	}

	res, err := es.client.Search(
		es.client.Search.WithContext(context.Background()),
		es.client.Search.WithIndex(INDEX_NAME),
		es.client.Search.WithBody(&buffer),
		es.client.Search.WithPretty(),
	)
	if err != nil {
		return nil, fmt.Errorf("Error when sending search request: %v", err)
	}
	defer res.Body.Close()

	var respBody map[string]interface{}
	if res.IsError() {
		if err := json.NewDecoder(res.Body).Decode(&respBody); err != nil {
			return nil, fmt.Errorf("Error when parsing error response body: %v", err)
		}
		return nil, errors.New(respBody["error"].(map[string]interface{})["reason"].(string))
	}

	if err := json.NewDecoder(res.Body).Decode(&respBody); err != nil {
		return nil, fmt.Errorf("Error when parsing response body: %v", err)
	}

	var users []models.User

	for _, hit := range respBody["hits"].(map[string]interface{})["hits"].([]interface{}) {
		users = append(users, models.User{
			ID:         uuid.MustParse(hit.(map[string]interface{})["_id"].(string)),
			Username:   hit.(map[string]interface{})["_source"].(map[string]interface{})["username"].(string),
			PictureURL: hit.(map[string]interface{})["_source"].(map[string]interface{})["picture"].(string),
		})
	}
	return users, nil
}
