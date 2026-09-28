package service

import (
	"fmt"
	"mime/multipart"
	"path/filepath"
	"strings"

	"github.com/google/uuid"

	"github.com/GangaRamPrasad2004/ECommerce-Platform/internal/interfaces"
)

// UploadService validates and stores product image uploads.
type UploadService struct {
	provider interfaces.UploadProvider
}

// NewUploadService creates an upload service using provider.
func NewUploadService(provider interfaces.UploadProvider) *UploadService {
	return &UploadService{provider: provider}
}

// UploadProductImage validates and stores an image for the specified product.
func (s *UploadService) UploadProductImage(productID uint, file *multipart.FileHeader) (string, error) {
	if file == nil {
		return "", fmt.Errorf("file is required")
	}

	ext := strings.ToLower(filepath.Ext(file.Filename))
	if !isValidImageExt(ext) {
		return "", fmt.Errorf("invalid file type: %s", ext)
	}
	newFileName := uuid.New().String()
	path := fmt.Sprintf("products/%d/%s%s", productID, newFileName, ext)

	return s.provider.UploadFile(file, path)
}

func isValidImageExt(ext string) bool {
	validExts := []string{".jpg", ".jpeg", ".png", ".gif", ".webp"}
	for _, validExt := range validExts {
		if ext == validExt {
			return true
		}
	}

	return false
}
