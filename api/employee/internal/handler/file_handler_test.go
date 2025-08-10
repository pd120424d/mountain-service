package handler

import (
	"bytes"
	"errors"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"

	employeeV1 "github.com/pd120424d/mountain-service/api/contracts/employee/v1"
	"github.com/pd120424d/mountain-service/api/employee/internal/service"
	"github.com/pd120424d/mountain-service/api/shared/storage"
	"github.com/pd120424d/mountain-service/api/shared/utils"
)

func TestFileHandler_UploadProfilePicture(t *testing.T) {
	t.Parallel()

	t.Run("it returns an error when employee ID is invalid", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockBlobService := storage.NewMockAzureBlobService(ctrl)
		mockEmployeeService := service.NewMockEmployeeService(ctrl)
		log := utils.NewTestLogger()
		handler := NewFileHandler(log, mockBlobService, mockEmployeeService)

		gin.SetMode(gin.TestMode)
		w := httptest.NewRecorder()
		ctx, _ := gin.CreateTestContext(w)
		ctx.Request = httptest.NewRequest("POST", "/employees/invalid/profile-picture", nil)
		ctx.Params = gin.Params{{Key: "id", Value: "invalid"}}

		handler.UploadProfilePicture(ctx)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		assert.Contains(t, w.Body.String(), "Invalid employee ID")
	})

	t.Run("it returns an error when no file is provided", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockBlobService := storage.NewMockAzureBlobService(ctrl)
		mockEmployeeService := service.NewMockEmployeeService(ctrl)
		log := utils.NewTestLogger()
		handler := NewFileHandler(log, mockBlobService, mockEmployeeService)

		gin.SetMode(gin.TestMode)
		w := httptest.NewRecorder()
		ctx, _ := gin.CreateTestContext(w)
		ctx.Request = httptest.NewRequest("POST", "/employees/123/profile-picture", nil)
		ctx.Params = gin.Params{{Key: "id", Value: "123"}}

		handler.UploadProfilePicture(ctx)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		assert.Contains(t, w.Body.String(), "No file provided or invalid file")
	})

	tests := []struct {
		name           string
		employeeID     string
		setupMocks     func(*storage.MockAzureBlobService, *service.MockEmployeeService)
		setupRequest   func() *http.Request
		expectedStatus int
		expectedBody   string
	}{
		{
			name:       "it succeeds when uploading valid image",
			employeeID: "123",
			setupMocks: func(mockBlob *storage.MockAzureBlobService, mockEmployee *service.MockEmployeeService) {
				expectedResult := &storage.UploadResult{
					BlobURL:  "https://test.blob.core.windows.net/container/test.jpg",
					BlobName: "employee-123/test.jpg",
					Size:     1024,
				}
				mockBlob.EXPECT().
					UploadProfilePicture(gomock.Any(), gomock.Any(), gomock.Any(), uint(123)).
					Return(expectedResult, nil)

				// Expect employee service to be called to update profile picture URL
				mockEmployee.EXPECT().
					UpdateEmployee(uint(123), gomock.Any()).
					Return(&employeeV1.EmployeeResponse{ID: 123}, nil)
			},
			setupRequest: func() *http.Request {
				body := &bytes.Buffer{}
				writer := multipart.NewWriter(body)
				part, _ := writer.CreateFormFile("file", "test.jpg")
				part.Write([]byte("fake image data"))
				writer.Close()

				req := httptest.NewRequest("POST", "/employees/123/profile-picture", body)
				req.Header.Set("Content-Type", writer.FormDataContentType())
				return req
			},
			expectedStatus: http.StatusOK,
			expectedBody:   `"message":"Profile picture uploaded successfully"`,
		},
		{
			name:       "it fails when blob service returns error",
			employeeID: "123",
			setupMocks: func(mockBlob *storage.MockAzureBlobService, mockEmployee *service.MockEmployeeService) {
				mockBlob.EXPECT().
					UploadProfilePicture(gomock.Any(), gomock.Any(), gomock.Any(), uint(123)).
					Return(nil, errors.New("upload failed"))
				// No expectation for employee service since blob upload fails
			},
			setupRequest: func() *http.Request {
				body := &bytes.Buffer{}
				writer := multipart.NewWriter(body)
				part, _ := writer.CreateFormFile("file", "test.jpg")
				part.Write([]byte("fake image data"))
				writer.Close()

				req := httptest.NewRequest("POST", "/employees/123/profile-picture", body)
				req.Header.Set("Content-Type", writer.FormDataContentType())
				return req
			},
			expectedStatus: http.StatusInternalServerError,
			expectedBody:   `"error":"upload failed"`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockBlobService := storage.NewMockAzureBlobService(ctrl)
			mockEmployeeService := service.NewMockEmployeeService(ctrl)
			tt.setupMocks(mockBlobService, mockEmployeeService)

			log := utils.NewTestLogger()
			handler := NewFileHandler(log, mockBlobService, mockEmployeeService)

			gin.SetMode(gin.TestMode)
			w := httptest.NewRecorder()
			ctx, _ := gin.CreateTestContext(w)
			ctx.Request = tt.setupRequest()
			ctx.Params = gin.Params{{Key: "id", Value: tt.employeeID}}

			handler.UploadProfilePicture(ctx)

			assert.Equal(t, tt.expectedStatus, w.Code)
			assert.Contains(t, w.Body.String(), tt.expectedBody)
		})
	}
}

func TestFileHandler_DeleteProfilePicture(t *testing.T) {
	t.Parallel()

	t.Run("it returns an error when employee ID is invalid", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockBlobService := storage.NewMockAzureBlobService(ctrl)
		mockEmployeeService := service.NewMockEmployeeService(ctrl)
		log := utils.NewTestLogger()
		handler := NewFileHandler(log, mockBlobService, mockEmployeeService)

		gin.SetMode(gin.TestMode)
		w := httptest.NewRecorder()
		ctx, _ := gin.CreateTestContext(w)
		ctx.Request = httptest.NewRequest("DELETE", "/employees/invalid/profile-picture?blobName=test.jpg", nil)
		ctx.Params = gin.Params{{Key: "id", Value: "invalid"}}

		handler.DeleteProfilePicture(ctx)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		assert.Contains(t, w.Body.String(), "Invalid employee ID")
	})

	t.Run("it returns an error when blob name is missing", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockBlobService := storage.NewMockAzureBlobService(ctrl)
		mockEmployeeService := service.NewMockEmployeeService(ctrl)
		log := utils.NewTestLogger()
		handler := NewFileHandler(log, mockBlobService, mockEmployeeService)

		gin.SetMode(gin.TestMode)
		w := httptest.NewRecorder()
		ctx, _ := gin.CreateTestContext(w)
		ctx.Request = httptest.NewRequest("DELETE", "/employees/123/profile-picture", nil)
		ctx.Params = gin.Params{{Key: "id", Value: "123"}}

		handler.DeleteProfilePicture(ctx)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		assert.Contains(t, w.Body.String(), "Blob name is required")
	})

	tests := []struct {
		name           string
		employeeID     string
		blobName       string
		setupMocks     func(*storage.MockAzureBlobService, *service.MockEmployeeService)
		expectedStatus int
		expectedBody   string
	}{
		{
			name:       "it succeeds when deleting existing blob",
			employeeID: "123",
			blobName:   "employee-123/test.jpg",
			setupMocks: func(mockBlob *storage.MockAzureBlobService, mockEmployee *service.MockEmployeeService) {
				mockBlob.EXPECT().
					DeleteProfilePicture(gomock.Any(), "employee-123/test.jpg").
					Return(nil)
				// No employee service call needed for delete
			},
			expectedStatus: http.StatusOK,
			expectedBody:   `"message":"Profile picture deleted successfully"`,
		},
		{
			name:       "it fails when blob service returns error",
			employeeID: "123",
			blobName:   "employee-123/test.jpg",
			setupMocks: func(mockBlob *storage.MockAzureBlobService, mockEmployee *service.MockEmployeeService) {
				mockBlob.EXPECT().
					DeleteProfilePicture(gomock.Any(), "employee-123/test.jpg").
					Return(errors.New("delete failed"))
				// No employee service call needed when delete fails
			},
			expectedStatus: http.StatusInternalServerError,
			expectedBody:   `"error":"Failed to delete profile picture"`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockBlobService := storage.NewMockAzureBlobService(ctrl)
			mockEmployeeService := service.NewMockEmployeeService(ctrl)
			tt.setupMocks(mockBlobService, mockEmployeeService)

			log := utils.NewTestLogger()
			handler := NewFileHandler(log, mockBlobService, mockEmployeeService)

			gin.SetMode(gin.TestMode)
			w := httptest.NewRecorder()
			ctx, _ := gin.CreateTestContext(w)
			ctx.Request = httptest.NewRequest("DELETE", "/employees/"+tt.employeeID+"/profile-picture?blobName="+tt.blobName, nil)
			ctx.Params = gin.Params{{Key: "id", Value: tt.employeeID}}

			handler.DeleteProfilePicture(ctx)

			assert.Equal(t, tt.expectedStatus, w.Code)
			assert.Contains(t, w.Body.String(), tt.expectedBody)
		})
	}
}

func TestFileHandler_GetProfilePictureInfo(t *testing.T) {
	t.Parallel()

	t.Run("it returns an error when blob name is missing", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockBlobService := storage.NewMockAzureBlobService(ctrl)
		mockEmployeeService := service.NewMockEmployeeService(ctrl)
		log := utils.NewTestLogger()
		handler := NewFileHandler(log, mockBlobService, mockEmployeeService)

		gin.SetMode(gin.TestMode)
		w := httptest.NewRecorder()
		ctx, _ := gin.CreateTestContext(w)
		ctx.Request = httptest.NewRequest("GET", "/files/profile-picture/info", nil)

		handler.GetProfilePictureInfo(ctx)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		assert.Contains(t, w.Body.String(), "Blob name is required")
	})

	t.Run("it returns profile picture info when blob name is provided", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockBlobService := storage.NewMockAzureBlobService(ctrl)
		mockEmployeeService := service.NewMockEmployeeService(ctrl)
		log := utils.NewTestLogger()
		handler := NewFileHandler(log, mockBlobService, mockEmployeeService)

		gin.SetMode(gin.TestMode)
		w := httptest.NewRecorder()
		ctx, _ := gin.CreateTestContext(w)
		ctx.Request = httptest.NewRequest("GET", "/files/profile-picture/info?blobName=test.jpg", nil)

		handler.GetProfilePictureInfo(ctx)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.Contains(t, w.Body.String(), "blobName")
		assert.Contains(t, w.Body.String(), "status")
	})
}
