package interfaces

import "mime/multipart"

// UploadProvider stores and deletes uploaded files.
type UploadProvider interface {
	UploadFile(file *multipart.FileHeader, path string) (string, error)
	DeleteFile(path string) error
}
