package v1

// AppError is a typed application error that can be used across services and handlers.
// Code should be a stable identifier (e.g., SHIFT_ERRORS.CONSECUTIVE_SHIFTS_LIMIT).
// Message is a human-readable message (optional; handlers or clients may localize).
// Details carries optional structured fields (e.g., {"limit": 6}).
type AppError struct {
	Code    string                 `json:"code"`
	Message string                 `json:"message,omitempty"`
	Details map[string]interface{} `json:"details,omitempty"`
}

func (e *AppError) Error() string {
	if e.Message != "" {
		return e.Message
	}
	return e.Code
}

func NewAppError(code, message string, details map[string]interface{}) *AppError {
	return &AppError{Code: code, Message: message, Details: details}
}

