package storage

import (
	"testing"

	"github.com/pd120424d/mountain-service/api/shared/utils"
	"github.com/stretchr/testify/assert"
)

func TestLoadConfigFromEnv(t *testing.T) {
	// Cannot use t.Parallel() with t.Setenv()

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

			config := loadConfigFromEnv()
			if config != tt.expected {
				t.Errorf("Expected config: %v, got: %v", tt.expected, config)
			}
		})
	}
}

func TestNewAzureBlobClientWrapper(t *testing.T) {
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
				AccountName:   "",
				AccountKey:    "testkey",
				ContainerName: "test-container",
			},
			expectError: true,
			errorMsg:    "azure storage account name and key are required",
		},
		{
			name: "it fails when account key is missing",
			config: AzureBlobConfig{
				AccountName:   "testaccount",
				AccountKey:    "",
				ContainerName: "test-container",
			},
			expectError: true,
			errorMsg:    "azure storage account name and key are required",
		},
		{
			name: "it succeeds when container name is missing (uses default)",
			config: AzureBlobConfig{
				AccountName:   "testaccount",
				AccountKey:    "dGVzdGtleQ==", // base64 encoded "testkey"
				ContainerName: "",
			},
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			wrapper, err := NewAzureBlobClientWrapper(log, tt.config)

			if tt.expectError {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.errorMsg)
				assert.Nil(t, wrapper)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, wrapper)

				// For the test that uses default container name, verify it was set
				if tt.config.ContainerName == "" {
					// We can't directly access the container name from the interface,
					// but we can test the GetBlobURL method to verify the default was used
					url := wrapper.GetBlobURL("test.jpg")
					assert.Contains(t, url, "employee-profiles") // default container name
				}
			}
		})
	}
}

func TestAzureBlobClientWrapper_GetBlobURL(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name        string
		accountName string
		container   string
		blobName    string
		expected    string
	}{
		{
			name:        "it returns correct URL for simple blob name",
			accountName: "testaccount",
			container:   "test-container",
			blobName:    "test.jpg",
			expected:    "https://testaccount.blob.core.windows.net/test-container/test.jpg",
		},
		{
			name:        "it returns correct URL for blob with path",
			accountName: "myaccount",
			container:   "images",
			blobName:    "profiles/user123/avatar.png",
			expected:    "https://myaccount.blob.core.windows.net/images/profiles/user123/avatar.png",
		},
		{
			name:        "it returns correct URL for blob with special characters",
			accountName: "storage123",
			container:   "documents",
			blobName:    "file%20with%20spaces.pdf",
			expected:    "https://storage123.blob.core.windows.net/documents/file%20with%20spaces.pdf",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			wrapper := &azureBlobClientWrapper{
				accountName:   tt.accountName,
				containerName: tt.container,
			}

			result := wrapper.GetBlobURL(tt.blobName)
			assert.Equal(t, tt.expected, result)
		})
	}
}

// TestAzureBlobClientWrapper_Structure tests the wrapper structure
func TestAzureBlobClientWrapper_Structure(t *testing.T) {
	t.Parallel()

	t.Run("it has correct structure when created", func(t *testing.T) {
		log := utils.NewTestLogger()

		wrapper := &azureBlobClientWrapper{
			log:           log,
			client:        nil, // Would be a real client in integration tests
			containerName: "test-container",
			accountName:   "testaccount",
		}

		// Test that the wrapper has the correct structure
		assert.Equal(t, "test-container", wrapper.containerName)
		assert.Equal(t, "testaccount", wrapper.accountName)
		assert.NotNil(t, wrapper.log)
	})
}

// TestAzureBlobClientWrapper_Integration demonstrates how the wrapper would be used
func TestAzureBlobClientWrapper_Integration(t *testing.T) {
	t.Parallel()

	t.Run("it demonstrates the complete workflow", func(t *testing.T) {
		log := utils.NewTestLogger()

		// This test demonstrates how the wrapper would be used in practice
		// In a real scenario, you would:
		// 1. Load config from environment
		// 2. Create the wrapper
		// 3. Create the blob service
		// 4. Use the service for operations

		config := AzureBlobConfig{
			AccountName:   "testaccount",
			AccountKey:    "dGVzdGtleQ==", // base64 encoded "testkey"
			ContainerName: "test-container",
		}

		// Step 1: Create wrapper (this would normally connect to Azure)
		// In this test, we just verify the structure
		wrapper, err := NewAzureBlobClientWrapper(log, config)
		assert.NoError(t, err)
		assert.NotNil(t, wrapper)

		// Step 2: Verify URL generation works correctly
		url := wrapper.GetBlobURL("test-file.jpg")
		expectedURL := "https://testaccount.blob.core.windows.net/test-container/test-file.jpg"
		assert.Equal(t, expectedURL, url)

		// Step 3: In a real scenario, you would create the blob service
		// service, err := NewAzureBlobService(log, wrapper)
		// assert.NoError(t, err)
		// assert.NotNil(t, service)

		// Note: Actual Azure operations (UploadStream, DeleteBlob, CreateContainer)
		// would require integration testing with real Azure credentials
	})
}
