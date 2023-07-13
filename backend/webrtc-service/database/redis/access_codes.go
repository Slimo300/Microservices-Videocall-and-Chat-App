package redis

import (
	"fmt"
	"strings"
)

// NewAccessCode creates a new access code for specific user with TTL specified by db object and saves it to redis
func (db *DB) NewAccessCode(groupID, userID, accessCode string) error {
	return db.Set(accessCode, fmt.Sprintf("%s:%s", userID, groupID), db.AccessCodeTTL).Err()
}

// CheckAccessCode checks if access code is present in storage and returns userID saved with it and information about error
func (db *DB) CheckAccessCode(accessCode string) (string, string, error) {
	res, err := db.Get(accessCode).Result()
	if err != nil {
		return "", "", err
	}
	userID := strings.Split(res, "")[0]
	groupID := strings.Split(res, "")[1]

	return userID, groupID, nil
}
