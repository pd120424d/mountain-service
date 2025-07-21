package model

// UrgencyCreateRequest DTO for creating a new urgency
// swagger:model
type UrgencyCreateRequest struct {
	Name         string `json:"name" binding:"required"`
	Email        string `json:"email" binding:"required,email"`
	ContactPhone string `json:"contactPhone" binding:"required"`
	Description  string `json:"description" binding:"required"`
}

// UrgencyUpdateRequest DTO for updating an urgency
// swagger:model
type UrgencyUpdateRequest struct {
	Name         string `json:"name"`
	Email        string `json:"email" binding:"email"`
	ContactPhone string `json:"contactPhone"`
	Description  string `json:"description"`
}

// UrgencyResponse DTO for returning an urgency
// swagger:model
type UrgencyResponse struct {
	ID           uint   `json:"id"`
	Name         string `json:"name"`
	Email        string `json:"email"`
	ContactPhone string `json:"contactPhone"`
	Description  string `json:"description"`
	CreatedAt    string `json:"createdAt"`
	UpdatedAt    string `json:"updatedAt"`
}

// UrgencyList DTO for returning a list of urgencies
// swagger:model
type UrgencyList struct {
	Urgencies []UrgencyResponse `json:"urgencies"`
}
