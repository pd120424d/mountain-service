package v1

// ErrorResponse DTO for returning an error message
// swagger:model
type ErrorResponse struct {
	Error string `json:"error" example:"Error message"`
}

// MessageResponse DTO for returning a success message
// swagger:model
type MessageResponse struct {
	Message string `json:"message" example:"Success message"`
}

// HealthResponse DTO for health check responses
// swagger:model
type HealthResponse struct {
	Status  string `json:"status" example:"healthy"`
	Service string `json:"service" example:"employee-service"`
	Version string `json:"version,omitempty" example:"1.0.0"`
}

// PaginationRequest DTO for pagination parameters
// swagger:model
type PaginationRequest struct {
	Page     int `json:"page" form:"page" binding:"min=1"`
	PageSize int `json:"pageSize" form:"pageSize" binding:"min=1,max=100"`
}

// PaginationResponse DTO for pagination metadata
// swagger:model
type PaginationResponse struct {
	Page       int   `json:"page"`
	PageSize   int   `json:"pageSize"`
	TotalItems int64 `json:"totalItems"`
	TotalPages int   `json:"totalPages"`
}

// ValidationError represents a field validation error
// swagger:model
type ValidationError struct {
	Field   string `json:"field" example:"email"`
	Message string `json:"message" example:"Invalid email format"`
}

// ValidationErrorResponse DTO for returning validation errors
// swagger:model
type ValidationErrorResponse struct {
	Error  string            `json:"error" example:"Validation failed"`
	Fields []ValidationError `json:"fields"`
}
