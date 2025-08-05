package storage

import (
	"bytes"
	"context"
	"errors"
	"mime/multipart"
	"net/textproto"
	"strings"
	"testing"

	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob"
	"github.com/pd120424d/mountain-service/api/shared/utils"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestAzureBlobService_NewAzureBlobService(t *testing.T) {
	t.Parallel()

	t.Run("it fails when account name is missing", func(t *testing.T) {
		log := utils.NewTestLogger()
		service, err := NewAzureBlobService(AzureBlobConfig{}, log)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "azure storage account name and key are required")
		assert.Nil(t, service)
	})

	t.Run("it fails when account key is missing", func(t *testing.T) {
		log := utils.NewTestLogger()
		service, err := NewAzureBlobService(AzureBlobConfig{AccountName: "testaccount"}, log)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "azure storage account name and key are required")
		assert.Nil(t, service)
	})
}

func TestLoadConfigFromEnv(t *testing.T) {
	tests := []struct {
		name     string
		envVars  map[string]string
		expected AzureBlobConfig
	}{
		{
			name: "it succeeds when all environment variables are set",
			envVars: map[string]string{
				"AZURE_STORAGE_ACCOUNT_NAME":   "testaccount",
				"AZURE_STORAGE_ACCOUNT_KEY":    "testkey",
				"AZURE_STORAGE_CONTAINER_NAME": "testcontainer",
			},
			expected: AzureBlobConfig{
				AccountName:   "testaccount",
				AccountKey:    "testkey",
				ContainerName: "testcontainer",
			},
		},
		{
			name: "it succeeds when using default container name",
			envVars: map[string]string{
				"AZURE_STORAGE_ACCOUNT_NAME": "testaccount",
				"AZURE_STORAGE_ACCOUNT_KEY":  "testkey",
			},
			expected: AzureBlobConfig{
				AccountName:   "testaccount",
				AccountKey:    "testkey",
				ContainerName: "employee-profiles",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Set environment variables
			for key, value := range tt.envVars {
				t.Setenv(key, value)
			}

			config := LoadConfigFromEnv()
			assert.Equal(t, tt.expected, config)
		})
	}
}

func TestValidateImageFile(t *testing.T) {
	t.Parallel()
	log := utils.NewTestLogger()

	// Create a service instance for testing validation
	service := &azureBlobService{
		log: log.WithName("AzureBlobService"),
	}

	tests := []struct {
		name        string
		header      *multipart.FileHeader
		expectError bool
		errorMsg    string
	}{
		{
			name: "it succeeds when file is valid JPEG",
			header: &multipart.FileHeader{
				Filename: "test.jpg",
				Size:     1024 * 1024, // 1MB
				Header: textproto.MIMEHeader{
					"Content-Type": []string{"image/jpeg"},
				},
			},
			expectError: false,
		},
		{
			name: "it succeeds when file is valid PNG",
			header: &multipart.FileHeader{
				Filename: "test.png",
				Size:     2 * 1024 * 1024, // 2MB
				Header: textproto.MIMEHeader{
					"Content-Type": []string{"image/png"},
				},
			},
			expectError: false,
		},
		{
			name: "it fails when file size exceeds limit",
			header: &multipart.FileHeader{
				Filename: "large.jpg",
				Size:     6 * 1024 * 1024, // 6MB
				Header: textproto.MIMEHeader{
					"Content-Type": []string{"image/jpeg"},
				},
			},
			expectError: true,
			errorMsg:    "file size exceeds maximum allowed size of 5MB",
		},
		{
			name: "it fails when file extension is invalid",
			header: &multipart.FileHeader{
				Filename: "test.txt",
				Size:     1024,
				Header: textproto.MIMEHeader{
					"Content-Type": []string{"text/plain"},
				},
			},
			expectError: true,
			errorMsg:    "unsupported file type: .txt",
		},
		{
			name: "it fails when MIME type is invalid",
			header: &multipart.FileHeader{
				Filename: "test.jpg",
				Size:     1024,
				Header: textproto.MIMEHeader{
					"Content-Type": []string{"text/plain"},
				},
			},
			expectError: true,
			errorMsg:    "unsupported content type: text/plain",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := service.validateImageFile(tt.header)

			if tt.expectError {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.errorMsg)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestGenerateBlobName(t *testing.T) {
	t.Parallel()
	log := utils.NewTestLogger()

	service := &azureBlobService{
		log: log.WithName("AzureBlobService"),
	}

	tests := []struct {
		name             string
		employeeID       uint
		originalFilename string
		expectedPrefix   string
		expectedSuffix   string
	}{
		{
			name:             "it succeeds when generating blob name for JPEG",
			employeeID:       123,
			originalFilename: "profile.jpg",
			expectedPrefix:   "employee-123/",
			expectedSuffix:   ".jpg",
		},
		{
			name:             "it succeeds when generating blob name for PNG",
			employeeID:       456,
			originalFilename: "avatar.png",
			expectedPrefix:   "employee-456/",
			expectedSuffix:   ".png",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			blobName := service.generateBlobName(tt.employeeID, tt.originalFilename)

			assert.True(t, len(blobName) > 0)
			assert.Contains(t, blobName, tt.expectedPrefix)
			assert.True(t, len(blobName) > len(tt.expectedPrefix)+len(tt.expectedSuffix))
			assert.True(t, blobName[len(blobName)-len(tt.expectedSuffix):] == tt.expectedSuffix)
		})
	}
}

func TestNewAzureBlobService(t *testing.T) {
	t.Parallel()
	log := utils.NewTestLogger()

	tests := []struct {
		name        string
		config      AzureBlobConfig
		expectError bool
		errorMsg    string
	}{
		{
			name: "it fails when account name is missing",
			config: AzureBlobConfig{
				AccountName: "",
				AccountKey:  "testkey",
			},
			expectError: true,
			errorMsg:    "azure storage account name and key are required",
		},
		{
			name: "it fails when account key is missing",
			config: AzureBlobConfig{
				AccountName: "testaccount",
				AccountKey:  "",
			},
			expectError: true,
			errorMsg:    "azure storage account name and key are required",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			service, err := NewAzureBlobService(tt.config, log)

			if tt.expectError {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.errorMsg)
				assert.Nil(t, service)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, service)
			}
		})
	}
}

// createTestFile creates a test multipart file for testing
func createTestFile(filename, content string) (multipart.File, *multipart.FileHeader, error) {
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	part, err := writer.CreateFormFile("file", filename)
	if err != nil {
		return nil, nil, err
	}

	_, err = part.Write([]byte(content))
	if err != nil {
		return nil, nil, err
	}

	err = writer.Close()
	if err != nil {
		return nil, nil, err
	}

	// Parse the multipart form to get the file
	reader := multipart.NewReader(body, writer.Boundary())
	form, err := reader.ReadForm(10 << 20) // 10MB max
	if err != nil {
		return nil, nil, err
	}

	files := form.File["file"]
	if len(files) == 0 {
		return nil, nil, errors.New("no file found")
	}

	fileHeader := files[0]
	file, err := fileHeader.Open()
	if err != nil {
		return nil, nil, err
	}

	return file, fileHeader, nil
}

func TestAzureBlobService_UploadProfilePicture(t *testing.T) {
	t.Parallel()

	t.Run("it fails when file validation fails - file too large", func(t *testing.T) {
		log := utils.NewTestLogger()
		service := &azureBlobService{
			log: log.WithName("AzureBlobService"),
		}

		// Create a file that's too large (over 5MB)
		largeContent := strings.Repeat("a", 6*1024*1024) // 6MB
		file, header, err := createTestFile("large.jpg", largeContent)
		assert.NoError(t, err)
		defer file.Close()

		ctx := context.Background()
		result, err := service.UploadProfilePicture(ctx, file, header, 123)

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Contains(t, err.Error(), "file size exceeds maximum allowed size")
	})

	t.Run("it fails when file has invalid extension", func(t *testing.T) {
		log := utils.NewTestLogger()
		service := &azureBlobService{
			log: log.WithName("AzureBlobService"),
		}

		file, header, err := createTestFile("test.txt", "test content")
		assert.NoError(t, err)
		defer file.Close()

		ctx := context.Background()
		result, err := service.UploadProfilePicture(ctx, file, header, 123)

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Contains(t, err.Error(), "unsupported file type")
	})

	t.Run("it fails when MIME type is invalid", func(t *testing.T) {
		log := utils.NewTestLogger()
		service := &azureBlobService{
			log: log.WithName("AzureBlobService"),
		}

		// Create a file with .jpg extension but wrong MIME type
		file, header, err := createTestFile("test.jpg", "test content")
		assert.NoError(t, err)
		defer file.Close()

		// Manually set wrong MIME type
		header.Header = textproto.MIMEHeader{
			"Content-Type": []string{"text/plain"},
		}

		ctx := context.Background()
		result, err := service.UploadProfilePicture(ctx, file, header, 123)

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Contains(t, err.Error(), "unsupported content type")
	})

	// Note: We can't test successful upload without a real Azure client
	// The method would require a properly initialized azblob.Client
	// These tests focus on validation logic that can be tested in isolation
}

func TestAzureBlobService_DeleteProfilePicture(t *testing.T) {
	t.Parallel()

	t.Run("it succeeds when blob name is empty", func(t *testing.T) {
		log := utils.NewTestLogger()
		service := &azureBlobService{
			log: log.WithName("AzureBlobService"),
		}

		ctx := context.Background()
		err := service.DeleteProfilePicture(ctx, "")

		// TODO: Test this properly when it is implemented correctly
		assert.NoError(t, err)
	})

}

// TestAzureBlobService_UploadProfilePicture_WithMocks tests the upload functionality with mocked Azure client
func TestAzureBlobService_UploadProfilePicture_WithMocks(t *testing.T) {
	t.Parallel()

	t.Run("it succeeds when uploading valid image with mocked client", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockClient := NewMockAzureBlobClient(ctrl)
		log := utils.NewTestLogger()

		service := &azureBlobService{
			client:        mockClient,
			containerName: "test-container",
			log:           log.WithName("AzureBlobService"),
		}

		// Create a valid test file
		file, header, err := createTestFile("test.jpg", "fake image content")
		assert.NoError(t, err)
		defer file.Close()

		// Set proper MIME type
		header.Header = textproto.MIMEHeader{
			"Content-Type": []string{"image/jpeg"},
		}

		// Mock the UploadStream call
		mockClient.EXPECT().
			UploadStream(gomock.Any(), "test-container", gomock.Any(), file, gomock.Any()).
			Return(azblob.UploadStreamResponse{}, nil)

		ctx := context.Background()
		result, err := service.UploadProfilePicture(ctx, file, header, 123)

		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Contains(t, result.BlobURL, "test-container")
		assert.NotEmpty(t, result.BlobName)
		assert.Equal(t, int64(18), result.Size) // "fake image content" length
	})

	t.Run("it fails when Azure client returns error", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockClient := NewMockAzureBlobClient(ctrl)
		log := utils.NewTestLogger()

		service := &azureBlobService{
			client:        mockClient,
			containerName: "test-container",
			log:           log.WithName("AzureBlobService"),
		}

		// Create a valid test file
		file, header, err := createTestFile("test.jpg", "fake image content")
		assert.NoError(t, err)
		defer file.Close()

		// Set proper MIME type
		header.Header = textproto.MIMEHeader{
			"Content-Type": []string{"image/jpeg"},
		}

		// Mock the UploadStream call to return an error
		mockClient.EXPECT().
			UploadStream(gomock.Any(), "test-container", gomock.Any(), file, gomock.Any()).
			Return(azblob.UploadStreamResponse{}, errors.New("Azure upload failed"))

		ctx := context.Background()
		result, err := service.UploadProfilePicture(ctx, file, header, 123)

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Contains(t, err.Error(), "failed to upload file to Azure Blob Storage")
	})
}

// TestAzureBlobService_DeleteProfilePicture_WithMocks tests the delete functionality with mocked Azure client
func TestAzureBlobService_DeleteProfilePicture_WithMocks(t *testing.T) {
	t.Parallel()

	t.Run("it succeeds when deleting existing blob with mocked client", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockClient := NewMockAzureBlobClient(ctrl)
		log := utils.NewTestLogger()

		service := &azureBlobService{
			client:        mockClient,
			containerName: "test-container",
			log:           log.WithName("AzureBlobService"),
		}

		// Mock the DeleteBlob call
		mockClient.EXPECT().
			DeleteBlob(gomock.Any(), "test-container", "test-blob.jpg", gomock.Any()).
			Return(azblob.DeleteBlobResponse{}, nil)

		ctx := context.Background()
		err := service.DeleteProfilePicture(ctx, "test-blob.jpg")

		assert.NoError(t, err)
	})

	t.Run("it fails when Azure client returns error", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockClient := NewMockAzureBlobClient(ctrl)
		log := utils.NewTestLogger()

		service := &azureBlobService{
			client:        mockClient,
			containerName: "test-container",
			log:           log.WithName("AzureBlobService"),
		}

		// Mock the DeleteBlob call to return an error
		mockClient.EXPECT().
			DeleteBlob(gomock.Any(), "test-container", "test-blob.jpg", gomock.Any()).
			Return(azblob.DeleteBlobResponse{}, errors.New("Azure delete failed"))

		ctx := context.Background()
		err := service.DeleteProfilePicture(ctx, "test-blob.jpg")

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "failed to delete blob")
	})

	t.Run("it succeeds when blob name is empty (nothing to delete)", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockClient := NewMockAzureBlobClient(ctrl)
		log := utils.NewTestLogger()

		service := &azureBlobService{
			client:        mockClient,
			containerName: "test-container",
			log:           log.WithName("AzureBlobService"),
		}

		// No mock expectations since the method should return early

		ctx := context.Background()
		err := service.DeleteProfilePicture(ctx, "")

		assert.NoError(t, err)
	})
}

// TestAzureBlobService_NewAzureBlobService_WithMocks tests the constructor with mocked Azure client
func TestAzureBlobService_NewAzureBlobService_WithMocks(t *testing.T) {
	t.Parallel()

	t.Run("it creates service and ensures container exists", func(t *testing.T) {
		// Note: This test demonstrates how we could test the constructor if we refactored it
		// to accept an AzureBlobClient interface instead of creating the client internally
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockClient := NewMockAzureBlobClient(ctrl)
		log := utils.NewTestLogger()

		// Mock the CreateContainer call (this would be called by ensureContainer)
		mockClient.EXPECT().
			CreateContainer(gomock.Any(), "test-container", gomock.Any()).
			Return(azblob.CreateContainerResponse{}, nil)

		service := &azureBlobService{
			client:        mockClient,
			containerName: "test-container",
			log:           log.WithName("AzureBlobService"),
		}

		// Test the ensureContainer method directly
		ctx := context.Background()
		err := service.ensureContainer(ctx)

		assert.NoError(t, err)
	})

	t.Run("it handles container already exists error gracefully", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockClient := NewMockAzureBlobClient(ctrl)
		log := utils.NewTestLogger()

		// Mock the CreateContainer call to return a "ContainerAlreadyExists" error
		mockClient.EXPECT().
			CreateContainer(gomock.Any(), "test-container", gomock.Any()).
			Return(azblob.CreateContainerResponse{}, errors.New("ContainerAlreadyExists: The specified container already exists"))

		service := &azureBlobService{
			client:        mockClient,
			containerName: "test-container",
			log:           log.WithName("AzureBlobService"),
		}

		// Test the ensureContainer method directly
		ctx := context.Background()
		err := service.ensureContainer(ctx)

		// The method should handle the "ContainerAlreadyExists" error gracefully
		assert.NoError(t, err)
	})

	t.Run("it fails when container creation returns other errors", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockClient := NewMockAzureBlobClient(ctrl)
		log := utils.NewTestLogger()

		// Mock the CreateContainer call to return a different error
		mockClient.EXPECT().
			CreateContainer(gomock.Any(), "test-container", gomock.Any()).
			Return(azblob.CreateContainerResponse{}, errors.New("access denied"))

		service := &azureBlobService{
			client:        mockClient,
			containerName: "test-container",
			log:           log.WithName("AzureBlobService"),
		}

		// Test the ensureContainer method directly
		ctx := context.Background()
		err := service.ensureContainer(ctx)

		// The method should return the error for non-"ContainerAlreadyExists" errors
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "access denied")
	})
}
