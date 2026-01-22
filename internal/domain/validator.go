package domain

// ValidationError represents a single validation error
type ValidationError struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

// ValidationErrors is a collection of validation errors
type ValidationErrors []ValidationError

func (v ValidationErrors) Error() string {
	if len(v) == 0 {
		return ""
	}
	return v[0].Message
}

// Validator defines the interface for validating domain objects
// Following Interface Segregation Principle - only validation methods are exposed
type Validator interface {
	// Validate validates a struct and returns validation errors if any
	Validate(obj interface{}) ValidationErrors
	// ValidateField validates a single field value against a tag
	ValidateField(field interface{}, tag string) ValidationErrors
}
