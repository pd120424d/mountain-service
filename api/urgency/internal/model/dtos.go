package model

import (
	"fmt"
	"net/mail"
	"strings"
)

// UrgencyCreateRequest DTO for creating a new urgency
// swagger:model
type UrgencyCreateRequest struct {
	Name         string       `json:"name" binding:"required"`
	Email        string       `json:"email" binding:"required,email"`
	ContactPhone string       `json:"contactPhone" binding:"required"`
	Location     string       `json:"location" binding:"required"`
	Description  string       `json:"description" binding:"required"`
	Level        UrgencyLevel `json:"level"`
}

// UrgencyUpdateRequest DTO for updating an urgency
// swagger:model
type UrgencyUpdateRequest struct {
	Name         string       `json:"name"`
	Email        string       `json:"email" binding:"email"`
	ContactPhone string       `json:"contactPhone"`
	Location     string       `json:"location"`
	Description  string       `json:"description"`
	Level        UrgencyLevel `json:"level"`
	Status       Status       `json:"status"`
}

// UrgencyResponse DTO for returning an urgency
// swagger:model
type UrgencyResponse struct {
	ID           uint         `json:"id"`
	Name         string       `json:"name"`
	Email        string       `json:"email"`
	ContactPhone string       `json:"contactPhone"`
	Location     string       `json:"location"`
	Description  string       `json:"description"`
	Level        UrgencyLevel `json:"level"`
	Status       Status       `json:"status"`
	CreatedAt    string       `json:"createdAt"`
	UpdatedAt    string       `json:"updatedAt"`
}

// UrgencyList DTO for returning a list of urgencies
// swagger:model
type UrgencyList struct {
	Urgencies []UrgencyResponse `json:"urgencies"`
}

func (r *UrgencyCreateRequest) Validate() error {
	if strings.TrimSpace(r.Name) == "" {
		return fmt.Errorf("name is required")
	}
	if strings.TrimSpace(r.Email) == "" {
		return fmt.Errorf("email is required")
	}
	if _, err := mail.ParseAddress(r.Email); err != nil {
		return fmt.Errorf("invalid email format")
	}
	if strings.TrimSpace(r.ContactPhone) == "" {
		return fmt.Errorf("contact phone is required")
	}
	if strings.TrimSpace(r.Location) == "" {
		return fmt.Errorf("location is required")
	}
	if strings.TrimSpace(r.Description) == "" {
		return fmt.Errorf("description is required")
	}
	if r.Level != "" && !r.Level.Valid() {
		return fmt.Errorf("invalid urgency level")
	}
	return nil
}

func (r *UrgencyUpdateRequest) Validate() error {
	if r.Email != "" {
		if _, err := mail.ParseAddress(r.Email); err != nil {
			return fmt.Errorf("invalid email format")
		}
	}
	if r.Level != "" && !r.Level.Valid() {
		return fmt.Errorf("invalid urgency level")
	}
	if r.Status != "" && !r.Status.Valid() {
		return fmt.Errorf("invalid status")
	}
	return nil
}
