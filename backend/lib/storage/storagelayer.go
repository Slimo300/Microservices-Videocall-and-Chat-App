package storage

import "mime/multipart"

type StorageLayer interface {
	UpdateProfilePicture(img multipart.File, key string) error
	DeleteProfilePicture(key string) error
}
