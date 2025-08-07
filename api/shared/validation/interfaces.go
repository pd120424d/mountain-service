package validation

import (
	"strings"
)

type Validator interface {
	Validate() error
}

type ValidationError struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

func (ve ValidationError) Error() string {
	return ve.Message
}

type ValidationErrors []ValidationError

func (ves ValidationErrors) Error() string {
	if len(ves) == 0 {
		return ""
	}
	if len(ves) == 1 {
		return ves[0].Error()
	}

	// Concatenate all error messages
	var messages []string
	for _, ve := range ves {
		messages = append(messages, ve.Error())
	}
	return strings.Join(messages, "; ")
}

func (ves ValidationErrors) HasErrors() bool {
	return len(ves) > 0
}

func (ves *ValidationErrors) Add(field, message string) {
	*ves = append(*ves, ValidationError{Field: field, Message: message})
}

func (ves *ValidationErrors) AddError(field string, err error) {
	if err != nil {
		*ves = append(*ves, ValidationError{Field: field, Message: err.Error()})
	}
}
