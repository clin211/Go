package storage

import "mime/multipart"

type Storage interface {
	Save(fileHeader *multipart.FileHeader, dstPath string) error
	SaveFromBytes(data []byte, dstPath string) error
	GetURL(objectKey string) (string, error)
	GetURLWithFilename(objectKey string, filename string) (string, error)
	Delete(objectKey string) error
}
