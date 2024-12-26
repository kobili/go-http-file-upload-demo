package storage_backends

import (
	"mime/multipart"
)

type StorageBackend interface {
	SaveFile(file multipart.File, path string, fileName string) (string, error)
	RetrieveFile(filePath string) ([]byte, error)
	DeleteFile(filePath string) error
}
