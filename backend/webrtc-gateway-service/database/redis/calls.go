package redis

import "strconv"

const (
	CALL_COUNTER_COL_NAME      = "counter"
	CALL_DOMAIN_COL_NAME       = "domain_name"
	INSTANCE_COUNTER_HASH_NAME = "instance_counter"
)

func (db *DB) AddConnection(callID, domainName string) error {

	newCount, err := db.HIncrBy(callID, CALL_COUNTER_COL_NAME, 1).Result()
	if err != nil {
		return err
	}

	if newCount == 1 {
		if err := db.HSet(callID, CALL_DOMAIN_COL_NAME, domainName).Err(); err != nil {
			return err
		}
	}

	return db.HIncrBy(INSTANCE_COUNTER_HASH_NAME, domainName, 1).Err()
}

func (db *DB) DeleteConnection(callID, domainName string) error {
	newCount, err := db.HIncrBy(callID, CALL_COUNTER_COL_NAME, -1).Result()
	if err != nil {
		return err
	}

	if newCount == 0 {
		if err := db.Del(callID).Err(); err != nil {
			return err
		}
	}

	return db.HIncrBy(INSTANCE_COUNTER_HASH_NAME, domainName, -1).Err()
}

func (db *DB) NewInstance(domainName string) error {
	return db.HSet(INSTANCE_COUNTER_HASH_NAME, domainName, 0).Err()
}

func (db *DB) GetCallInstanceDomainName(callID string) (string, error) {
	return db.HGet(callID, CALL_DOMAIN_COL_NAME).Result()
}

func (db *DB) GetLeastUsedInstanceDomainName() (string, error) {
	instances, err := db.HGetAll(INSTANCE_COUNTER_HASH_NAME).Result()
	if err != nil {
		return "", err
	}

	var domainName string
	min := int(^uint(0) >> 1)

	for instance, counter := range instances {
		intCounter, err := strconv.Atoi(counter)
		if err != nil {
			return "", err
		}
		if intCounter < min {
			min = intCounter
			domainName = instance
		}
	}

	return domainName, nil
}
