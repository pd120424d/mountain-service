package storage

//go:generate mockgen -source=azure_blob_service.go -destination=azure_blob_service_gomock.go -package=storage -imports=gomock=go.uber.org/mock/gomock

import (
	"context"
	"fmt"
	"mime/multipart"
	"path/filepath"
	"strings"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob"
	"github.com/google/uuid"
	"github.com/pd120424d/mountain-service/api/shared/utils"
)

const (
	maxSize              = int64(5 * 1024 * 1024) // 5MB
	defaultContainerName = "employee-profiles"
)

type AzureBlobService interface {
	UploadProfilePicture(ctx context.Context, file multipart.File, header *multipart.FileHeader, employeeID uint) (*UploadResult, error)
	DeleteProfilePicture(ctx context.Context, blobName string) error
}

type azureBlobService struct {
	client AzureBlobClientWrapper
	log    utils.Logger
}

type UploadResult struct {
	BlobURL  string `json:"blobUrl"`
	BlobName string `json:"blobName"`
	Size     int64  `json:"size"`
}

func NewAzureBlobService(log utils.Logger, client AzureBlobClientWrapper) (AzureBlobService, error) {

	service := &azureBlobService{
		log:    log.WithName("AzureBlobService"),
		client: client,
	}

	if err := service.ensureContainer(context.Background()); err != nil {
		return nil, fmt.Errorf("failed to ensure container exists: %w", err)
	}

	return service, nil
}

func (s *azureBlobService) UploadProfilePicture(ctx context.Context, file multipart.File, header *multipart.FileHeader, employeeID uint) (*UploadResult, error) {
	if err := s.validateImageFile(header); err != nil {
		return nil, err
	}

	blobName := s.generateBlobName(employeeID, header.Filename)

	fileSize := header.Size

	contentType := header.Header.Get("Content-Type")
	if contentType == "" {
		contentType = "application/octet-stream"
	}

	_, err := s.client.UploadStream(ctx, blobName, file, &azblob.UploadStreamOptions{
		BlockSize:   int64(1024 * 1024), // 1MB blocks
		Concurrency: 3,
		Metadata: map[string]*string{
			"employeeId":  stringPtr(fmt.Sprintf("%d", employeeID)),
			"uploadedAt":  stringPtr(time.Now().UTC().Format(time.RFC3339)),
			"contentType": stringPtr(contentType),
		},
	})

	if err != nil {
		s.log.Errorf("Failed to upload blob %s: %v", blobName, err)
		return nil, fmt.Errorf("failed to upload file to Azure Blob Storage: %w", err)
	}

	blobURL := s.client.GetBlobURL(blobName)

	result := &UploadResult{
		BlobURL:  blobURL,
		BlobName: blobName,
		Size:     fileSize,
	}

	s.log.Infof("Successfully uploaded profile picture for employee %d: %s", employeeID, blobURL)
	return result, nil
}

func (s *azureBlobService) DeleteProfilePicture(ctx context.Context, blobName string) error {
	if blobName == "" {
		return nil // Nothing to delete
	}

	_, err := s.client.DeleteBlob(ctx, blobName, nil)
	if err != nil {
		s.log.Errorf("Failed to delete blob %s: %v", blobName, err)
		return fmt.Errorf("failed to delete blob: %w", err)
	}

	s.log.Infof("Successfully deleted blob: %s", blobName)
	return nil
}

func (s *azureBlobService) ensureContainer(ctx context.Context) error {
	_, err := s.client.CreateContainer(ctx, &azblob.CreateContainerOptions{Access: nil})
	if err != nil {
		if !strings.Contains(err.Error(), "ContainerAlreadyExists") {
			return err
		}
	}
	return nil
}

func (s *azureBlobService) validateImageFile(header *multipart.FileHeader) error {
	if header.Size > maxSize {
		return fmt.Errorf("file size exceeds maximum allowed size of 5MB")
	}

	ext := strings.ToLower(filepath.Ext(header.Filename))
	allowedExts := map[string]bool{
		".jpg":  true,
		".jpeg": true,
		".png":  true,
		".gif":  true,
		".webp": true,
	}

	if !allowedExts[ext] {
		return fmt.Errorf("unsupported file type: %s. Allowed types: jpg, jpeg, png, gif, webp", ext)
	}

	contentType := header.Header.Get("Content-Type")
	allowedMimes := map[string]bool{
		"image/jpeg": true,
		"image/jpg":  true,
		"image/png":  true,
		"image/gif":  true,
		"image/webp": true,
	}

	if !allowedMimes[contentType] {
		return fmt.Errorf("unsupported content type: %s", contentType)
	}

	return nil
}

func (s *azureBlobService) generateBlobName(employeeID uint, originalFilename string) string {
	ext := filepath.Ext(originalFilename)
	timestamp := time.Now().UTC().Format("20060102-150405")
	uniqueID := uuid.New().String()[:8]

	return fmt.Sprintf("employee-%d/%s-%s%s", employeeID, timestamp, uniqueID, ext)
}

func stringPtr(s string) *string {
	return &s
}
