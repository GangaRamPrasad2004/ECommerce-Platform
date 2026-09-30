package providers

import (
	"mime/multipart"
	"os"
	"path/filepath"
)

// LocalUploadProvider stores uploaded files on the local filesystem.
type LocalUploadProvider struct {
	basePath string
}

// NewLocalUploadProvider creates a local upload provider rooted at basePath.
func NewLocalUploadProvider(basePath string) *LocalUploadProvider {
	return &LocalUploadProvider{basePath: basePath}
}

// UploadFile stores a multipart file at the requested relative path.
func (p *LocalUploadProvider) UploadFile(file *multipart.FileHeader, path string) (string, error) {

	fullPath := filepath.Join(p.basePath, path)

	if err := os.MkdirAll(filepath.Dir(fullPath), 0o755); err != nil {
		return "", err
	}

	// Open source
	src, err := file.Open()
	if err != nil {
		return "", err
	}
	defer func() { _ = src.Close() }()

	// create destination
	dst, err := os.Create(fullPath)
	if err != nil {
		return "", err
	}
	defer func() { _ = dst.Close() }()

	// read from source to destination
	if _, err := dst.ReadFrom(src); err != nil {
		return "", err
	}

	return path, nil

}

// DeleteFile removes a previously uploaded local file.
func (p *LocalUploadProvider) DeleteFile(path string) error {
	fullPath := filepath.Join(p.basePath, path)
	return os.Remove(fullPath)
}
