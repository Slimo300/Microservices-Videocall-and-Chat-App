package redis

// NewAccessCode creates a new access code for specific user with TTL specified by db object and saves it to redis
func (db *DB) NewAccessCode(accessCode, memberID string) error {
	return db.Set(accessCode, memberID, db.AccessCodeTTL).Err()
}

// CheckAccessCode checks if access code is present in storage and returns userID saved with it and information about error
func (db *DB) CheckAccessCode(accessCode string) (string, error) {
	return db.Get(accessCode).Result()
}
