package elastic

import (
	"encoding/json"
	"fmt"
	"strings"

	elasticsearch "github.com/elastic/go-elasticsearch/v8"
)

const INDEX_NAME = "users"
const INDEX_EXISTS_EXCEPTION = "resource_already_exists_exception"
const MINIMUM_SHOULD_MATCH = "60%"

type elasticSearchDB struct {
	client *elasticsearch.Client
}

func NewElasticSearchDB(addresses []string, username, password string) (*elasticSearchDB, error) {
	client, err := elasticsearch.NewClient(elasticsearch.Config{
		Addresses: addresses,
		Username:  username,
		Password:  password,
	})
	if err != nil {
		return nil, err
	}

	es := &elasticSearchDB{client: client}
	if err := es.setupIndex(); err != nil {
		return nil, err
	}

	return es, nil
}

func (es *elasticSearchDB) setupIndex() error {
	// checking if index 'INDEX_NAME' exists
	exists, err := es.client.Indices.Exists([]string{INDEX_NAME})
	if err != nil {
		return fmt.Errorf("Error when checking if index exists: %v", err)
	}
	defer exists.Body.Close()

	var respBody map[string]interface{}
	if exists.IsError() {
		if err := json.NewDecoder(exists.Body).Decode(&respBody); err != nil {
			return fmt.Errorf("Error decoding error body message: %v", err)
		}
		cause := respBody["error"].(map[string]interface{})["root_cause"].(map[string]interface{})["type"]
		if cause != INDEX_EXISTS_EXCEPTION {
			return fmt.Errorf("Unexpected error from request: %v", err)
		}
		// index exists return nil
		return nil
	}

	// creating index
	res, err := es.client.Indices.Create(
		INDEX_NAME,
		es.client.Indices.Create.WithBody(strings.NewReader(MAPPING)),
	)
	if err != nil {
		return fmt.Errorf("Error creating index: %v", err)
	}

	if res.IsError() {
		return fmt.Errorf("Error when creating database: %v", err)
	}

	return nil
}
