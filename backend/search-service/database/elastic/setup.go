package elastic

import (
	"fmt"
	"strings"

	elasticsearch "github.com/elastic/go-elasticsearch/v8"
)

const INDEX_NAME = "users"
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
		return fmt.Errorf("error when checking if index exists: %v", err)
	}
	defer exists.Body.Close()

	if exists.IsError() {
		if exists.StatusCode != 404 {
			return fmt.Errorf("unexpected status code: %d", exists.StatusCode)
		}
	} else {
		return nil
	}

	// creating index
	res, err := es.client.Indices.Create(
		INDEX_NAME,
		es.client.Indices.Create.WithBody(strings.NewReader(MAPPING)),
	)
	if err != nil {
		return fmt.Errorf("error creating index: %v", err)
	}

	if res.IsError() {
		return fmt.Errorf("server returned unexpected status code: %d", res.StatusCode)
	}

	return nil
}
