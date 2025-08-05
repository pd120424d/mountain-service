package handler

//go:generate mockgen -source=file_handler.go -destination=file_handler_gomock.go -package=handler -imports=gomock=go.uber.org/mock/gomock

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	employeeV1 "github.com/pd120424d/mountain-service/api/contracts/employee/v1"
	"github.com/pd120424d/mountain-service/api/employee/internal/service"
	"github.com/pd120424d/mountain-service/api/shared/storage"
	"github.com/pd120424d/mountain-service/api/shared/utils"
)

type FileHandler interface {
	UploadProfilePicture(ctx *gin.Context)
	DeleteProfilePicture(ctx *gin.Context)
	GetProfilePictureInfo(ctx *gin.Context)
}

type fileHandler struct {
	log             utils.Logger
	blobService     storage.AzureBlobService
	employeeService service.EmployeeService
}

func NewFileHandler(log utils.Logger, blobService storage.AzureBlobService, employeeService service.EmployeeService) FileHandler {
	return &fileHandler{
		log:             log.WithName("FileHandler"),
		blobService:     blobService,
		employeeService: employeeService,
	}
}

type UploadProfilePictureResponse struct {
	BlobURL  string `json:"blobUrl"`
	BlobName string `json:"blobName"`
	Size     int64  `json:"size"`
	Message  string `json:"message"`
}

// UploadProfilePicture handles profile picture upload
// @Summary Upload profile picture
// @Description Upload a profile picture for an employee
// @Tags files
// @Accept multipart/form-data
// @Produce json
// @Param employeeId path int true "Employee ID"
// @Param file formData file true "Profile picture file"
// @Success 200 {object} UploadProfilePictureResponse
// @Failure 400 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /employees/{employeeId}/profile-picture [post]
func (h *fileHandler) UploadProfilePicture(ctx *gin.Context) {
	h.log.Info("Received profile picture upload request")

	employeeIDStr := ctx.Param("employeeId")
	employeeID, err := strconv.ParseUint(employeeIDStr, 10, 32)
	if err != nil {
		h.log.Errorf("Invalid employee ID: %s", employeeIDStr)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid employee ID"})
		return
	}

	file, header, err := ctx.Request.FormFile("file")
	if err != nil {
		h.log.Errorf("Failed to get file from form: %v", err)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "No file provided or invalid file"})
		return
	}
	defer file.Close()

	h.log.Infof("Uploading profile picture for employee %d: %s (size: %d bytes)",
		employeeID, header.Filename, header.Size)

	result, err := h.blobService.UploadProfilePicture(ctx.Request.Context(), file, header, uint(employeeID))
	if err != nil {
		h.log.Errorf("Failed to upload profile picture: %v", err)

		errorMsg := "Failed to upload profile picture"
		if err.Error() != "" {
			errorMsg = err.Error()
		}

		ctx.JSON(http.StatusInternalServerError, gin.H{"error": errorMsg})
		return
	}

	// Update employee record with the new profile picture URL
	updateRequest := employeeV1.EmployeeUpdateRequest{
		ProfilePicture: result.BlobURL,
	}

	_, err = h.employeeService.UpdateEmployee(uint(employeeID), updateRequest)
	if err != nil {
		h.log.Errorf("Failed to update employee profile picture URL: %v", err)
		// Note: We don't return an error here because the upload was successful
		// The frontend will still get the blob URL and can handle this gracefully
	}

	response := UploadProfilePictureResponse{
		BlobURL:  result.BlobURL,
		BlobName: result.BlobName,
		Size:     result.Size,
		Message:  "Profile picture uploaded successfully",
	}

	h.log.Infof("Successfully uploaded profile picture for employee %d", employeeID)
	ctx.JSON(http.StatusOK, response)
}

// DeleteProfilePicture handles profile picture deletion
// @Summary Delete profile picture
// @Description Delete a profile picture for an employee
// @Tags files
// @Produce json
// @Param employeeId path int true "Employee ID"
// @Param blobName query string true "Blob name to delete"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /employees/{employeeId}/profile-picture [delete]
func (h *fileHandler) DeleteProfilePicture(ctx *gin.Context) {
	h.log.Info("Received profile picture delete request")

	employeeIDStr := ctx.Param("employeeId")
	employeeID, err := strconv.ParseUint(employeeIDStr, 10, 32)
	if err != nil {
		h.log.Errorf("Invalid employee ID: %s", employeeIDStr)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid employee ID"})
		return
	}

	blobName := ctx.Query("blobName")
	if blobName == "" {
		h.log.Error("Blob name is required")
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Blob name is required"})
		return
	}

	h.log.Infof("Deleting profile picture for employee %d: %s", employeeID, blobName)

	err = h.blobService.DeleteProfilePicture(ctx.Request.Context(), blobName)
	if err != nil {
		h.log.Errorf("Failed to delete profile picture: %v", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete profile picture"})
		return
	}

	h.log.Infof("Successfully deleted profile picture for employee %d", employeeID)
	ctx.JSON(http.StatusOK, gin.H{"message": "Profile picture deleted successfully"})
}

// GetProfilePictureInfo gets information about a profile picture
// @Summary Get profile picture info
// @Description Get information about a profile picture
// @Tags files
// @Produce json
// @Param blobName query string true "Blob name"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Router /files/profile-picture/info [get]
func (h *fileHandler) GetProfilePictureInfo(ctx *gin.Context) {
	blobName := ctx.Query("blobName")
	if blobName == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Blob name is required"})
		return
	}

	// For now, just return basic info
	ctx.JSON(http.StatusOK, gin.H{
		"blobName": blobName,
		"status":   "exists", // This would be determined by checking Azure
	})
}
