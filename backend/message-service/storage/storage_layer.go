package storage

const MAX_BUCKET_SIZE = 4900000000
const DEFAULT_REGION = "eu-central-1"

// StorageLayer defines functionality expected from Storage
type StorageLayer interface {
	GetPresignedGetRequests(string, ...GetFileInput) ([]GetFileOutput, error)
	GetPresignedPutRequests(string, ...PutFileInput) ([]PutFileOutput, error)
	DeleteFolder(folder string) error
	DeleteFile(key string) error
}

type PutFileInput struct {
	Name string `json:"name"`
	Size int64  `json:"size"`
}

type PutFileOutput struct {
	Name         string `json:"name"`
	Key          string `json:"key"`
	PresignedURL string `json:"url"`
}

type GetFileInput string

type GetFileOutput struct {
	Key          string `json:"key"`
	PresignedURL string `json:"url"`
}
