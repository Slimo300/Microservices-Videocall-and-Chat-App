package redis

import "github.com/google/uuid"

// NewAccessCode creates a new access code for specific user with TTL specified by db object and saves it to redis
func (db *DB) NewAccessCode(userID uuid.UUID, accessCode string) error {
	return db.Set(accessCode, userID.String(), db.AccessCodeTTL).Err()
}

// CheckAccessCode checks if access code is present in storage and returns userID saved with it and information about error
func (db *DB) CheckAccessCode(accessCode string) (uuid.UUID, error) {
	user, err := db.Get(accessCode).Result()
	if err != nil {
		return uuid.Nil, err
	}

	if err := db.Del(accessCode).Err(); err != nil {
		return uuid.Nil, err
	}

	userUID, err := uuid.Parse(user)
	if err != nil {
		return uuid.Nil, err
	}

	return userUID, nil
}
