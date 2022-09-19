package storage

import "mime/multipart"

type MockStorage struct{}

func (m MockStorage) UpdateProfilePicture(img multipart.File, key string) error {
	return nil
}

func (m MockStorage) DeleteProfilePicture(key string) error {
	return nil
}
