package model

// EmployeeResponse DTO for returning employee data
type EmployeeResponse struct {
	ID        uint   `json:"id"`
	Username  string `json:"username"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Gender    string `json:"gender"`
	Phone     string `json:"phone"`
	Email     string `json:"email"`
	// this may be represented as a byte array if we read the picture from somewhere for an example
	ProfilePicture string `json:"profile_picture"`
	ProfileType    string `json:"profile_type"`
}

// EmployeeCreateRequest DTO for creating a new employee
type EmployeeCreateRequest struct {
	FirstName      string `json:"first_name" binding:"required"`
	LastName       string `json:"last_name" binding:"required"`
	Username       string `json:"username" binding:"required"`
	Password       string `json:"password" binding:"required"`
	Email          string `json:"email" binding:"required,email"`
	Gender         string `json:"gender" binding:"required"`
	Phone          string `json:"phone" binding:"required"`
	ProfilePicture string `json:"profile_picture"`
	ProfileType    string `json:"profile_type"`
}

// EmployeeUpdateRequest DTO for updating an existing employee
type EmployeeUpdateRequest struct {
	FirstName string `json:"first_name,omitempty"`
	LastName  string `json:"last_name,omitempty"`
	Age       int    `json:"age,omitempty"`
	Email     string `json:"email,omitempty"`
}
