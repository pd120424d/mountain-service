package storage

//go:generate mockgen -source=azure_storage_client_wrapper.go -destination=azure_storage_client_wrapper_gomock.go -package=storage -imports=gomock=go.uber.org/mock/gomock

import (
	"context"
	"fmt"
	"io"
	"os"

	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob"
	"github.com/pd120424d/mountain-service/api/shared/utils"
)

type AzureBlobConfig struct {
	AccountName   string
	AccountKey    string
	ContainerName string
}

// AzureBlobClientWrapper interface wraps the Azure blob client operations for testing
type AzureBlobClientWrapper interface {
	UploadStream(ctx context.Context, blobName string, body io.Reader, options *azblob.UploadStreamOptions) (azblob.UploadStreamResponse, error)
	DeleteBlob(ctx context.Context, blobName string, options *azblob.DeleteBlobOptions) (azblob.DeleteBlobResponse, error)
	CreateContainer(ctx context.Context, options *azblob.CreateContainerOptions) (azblob.CreateContainerResponse, error)
	GetBlobURL(blobName string) string
}

// azureBlobClientWrapper wraps the real Azure client
type azureBlobClientWrapper struct {
	log           utils.Logger
	client        *azblob.Client
	containerName string
	accountName   string
}

func NewAzureBlobClientWrapper(log utils.Logger, config AzureBlobConfig) (AzureBlobClientWrapper, error) {
	if config.AccountName == "" || config.AccountKey == "" {
		return nil, fmt.Errorf("azure storage account name and key are required")
	}

	if config.ContainerName == "" {
		config.ContainerName = defaultContainerName
	}

	credential, err := azblob.NewSharedKeyCredential(config.AccountName, config.AccountKey)
	if err != nil {
		return nil, fmt.Errorf("failed to create Azure credentials: %w", err)
	}

	serviceURL := fmt.Sprintf("https://%s.blob.core.windows.net/", config.AccountName)

	client, err := azblob.NewClientWithSharedKeyCredential(serviceURL, credential, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create Azure blob client: %w", err)
	}
	return &azureBlobClientWrapper{
		log:           log,
		client:        client,
		containerName: config.ContainerName,
		accountName:   config.AccountName,
	}, nil
}

func (w *azureBlobClientWrapper) UploadStream(ctx context.Context, blobName string, body io.Reader, options *azblob.UploadStreamOptions) (azblob.UploadStreamResponse, error) {
	return w.client.UploadStream(ctx, w.containerName, blobName, body, options)
}

func (w *azureBlobClientWrapper) DeleteBlob(ctx context.Context, blobName string, options *azblob.DeleteBlobOptions) (azblob.DeleteBlobResponse, error) {
	return w.client.DeleteBlob(ctx, w.containerName, blobName, options)
}

func (w *azureBlobClientWrapper) CreateContainer(ctx context.Context, options *azblob.CreateContainerOptions) (azblob.CreateContainerResponse, error) {
	return w.client.CreateContainer(ctx, w.containerName, options)
}

func (w *azureBlobClientWrapper) GetBlobURL(blobName string) string {
	return fmt.Sprintf("https://%s.blob.core.windows.net/%s/%s", w.accountName, w.containerName, blobName)
}

func loadConfigFromEnv() AzureBlobConfig {
	return AzureBlobConfig{
		AccountName:   os.Getenv("AZURE_STORAGE_ACCOUNT_NAME"),
		AccountKey:    os.Getenv("AZURE_STORAGE_ACCOUNT_KEY"),
		ContainerName: getEnvOrDefault("AZURE_STORAGE_CONTAINER_NAME", "employee-profiles"),
	}
}

func getEnvOrDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
