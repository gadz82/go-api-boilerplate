package validation

import (
	"testing"

	"github.com/gadz82/go-api-boilerplate/internal/domain"
	"github.com/stretchr/testify/assert"
)

func TestNewValidator(t *testing.T) {
	v := NewValidator()
	assert.NotNil(t, v)
}

func TestPlaygroundValidator_Validate_Success(t *testing.T) {
	v := NewValidator()

	item := &domain.Item{
		ID:          "550e8400-e29b-41d4-a716-446655440000",
		Title:       "Test Item",
		Description: "A test description",
	}

	errors := v.Validate(item)
	assert.Nil(t, errors)
}

func TestPlaygroundValidator_Validate_RequiredField(t *testing.T) {
	v := NewValidator()

	item := &domain.Item{
		ID:          "550e8400-e29b-41d4-a716-446655440000",
		Title:       "", // Required field is empty
		Description: "A test description",
	}

	errors := v.Validate(item)
	assert.NotNil(t, errors)
	assert.Len(t, errors, 1)
	assert.Equal(t, "title", errors[0].Field)
	assert.Contains(t, errors[0].Message, "required")
}

func TestPlaygroundValidator_Validate_InvalidUUID(t *testing.T) {
	v := NewValidator()

	item := &domain.Item{
		ID:          "invalid-uuid",
		Title:       "Test Item",
		Description: "A test description",
	}

	errors := v.Validate(item)
	assert.NotNil(t, errors)
	assert.Len(t, errors, 1)
	assert.Equal(t, "i_d", errors[0].Field) // toSnakeCase converts "ID" to "i_d"
	assert.Contains(t, errors[0].Message, "UUID")
}

func TestPlaygroundValidator_Validate_EmptyUUID(t *testing.T) {
	v := NewValidator()

	// Empty UUID should be valid (omitempty)
	item := &domain.Item{
		ID:          "",
		Title:       "Test Item",
		Description: "A test description",
	}

	errors := v.Validate(item)
	assert.Nil(t, errors)
}

func TestPlaygroundValidator_Validate_MinLength(t *testing.T) {
	v := NewValidator()

	// Title must be at least 1 character
	item := &domain.Item{
		ID:          "550e8400-e29b-41d4-a716-446655440000",
		Title:       "", // Empty, violates min=1
		Description: "A test description",
	}

	errors := v.Validate(item)
	assert.NotNil(t, errors)
}

func TestPlaygroundValidator_Validate_MaxLength(t *testing.T) {
	v := NewValidator()

	// Create a string longer than 255 characters for title
	longTitle := make([]byte, 256)
	for i := range longTitle {
		longTitle[i] = 'a'
	}

	item := &domain.Item{
		ID:          "550e8400-e29b-41d4-a716-446655440000",
		Title:       string(longTitle),
		Description: "A test description",
	}

	errors := v.Validate(item)
	assert.NotNil(t, errors)
	assert.Equal(t, "title", errors[0].Field)
	assert.Contains(t, errors[0].Message, "at most")
}

func TestPlaygroundValidator_ValidateField_Success(t *testing.T) {
	v := NewValidator()

	errors := v.ValidateField("test@example.com", "email")
	assert.Nil(t, errors)
}

func TestPlaygroundValidator_ValidateField_InvalidEmail(t *testing.T) {
	v := NewValidator()

	errors := v.ValidateField("invalid-email", "email")
	assert.NotNil(t, errors)
	assert.Contains(t, errors[0].Message, "email")
}

func TestPlaygroundValidator_ValidateField_Required(t *testing.T) {
	v := NewValidator()

	errors := v.ValidateField("", "required")
	assert.NotNil(t, errors)
	assert.Contains(t, errors[0].Message, "required")
}

func TestPlaygroundValidator_ValidateField_URL(t *testing.T) {
	v := NewValidator()

	// Valid URL
	errors := v.ValidateField("https://example.com", "url")
	assert.Nil(t, errors)

	// Invalid URL
	errors = v.ValidateField("not-a-url", "url")
	assert.NotNil(t, errors)
	assert.Contains(t, errors[0].Message, "URL")
}

func TestPlaygroundValidator_Validate_ItemProperty(t *testing.T) {
	v := NewValidator()

	property := &domain.ItemProperty{
		ID:     "550e8400-e29b-41d4-a716-446655440000",
		ItemID: "550e8400-e29b-41d4-a716-446655440001",
		Name:   "color",
		Value:  "red",
	}

	errors := v.Validate(property)
	assert.Nil(t, errors)
}

func TestPlaygroundValidator_Validate_ItemProperty_MissingRequired(t *testing.T) {
	v := NewValidator()

	property := &domain.ItemProperty{
		ID:     "550e8400-e29b-41d4-a716-446655440000",
		ItemID: "550e8400-e29b-41d4-a716-446655440001",
		Name:   "", // Required
		Value:  "", // Required
	}

	errors := v.Validate(property)
	assert.NotNil(t, errors)
	assert.Len(t, errors, 2)
}

func TestToSnakeCase(t *testing.T) {
	// Note: The toSnakeCase function adds underscore before each uppercase letter
	// So "ID" becomes "i_d", "HTTPServer" becomes "h_t_t_p_server"
	tests := []struct {
		input    string
		expected string
	}{
		{"Title", "title"},
		{"ItemID", "item_i_d"},
		{"CreatedAt", "created_at"},
		{"ID", "i_d"},
		{"simple", "simple"},
		{"HTTPServer", "h_t_t_p_server"},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			result := toSnakeCase(tt.input)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestFormatErrorMessage_AllTags(t *testing.T) {
	v := NewValidator()

	// Test various validation tags by triggering them
	type TestStruct struct {
		Required string `validate:"required"`
		Min      string `validate:"min=5"`
		Max      string `validate:"max=3"`
		Email    string `validate:"email"`
		URL      string `validate:"url"`
		UUID     string `validate:"uuid4"`
	}

	// Test required
	ts := &TestStruct{Required: ""}
	errors := v.Validate(ts)
	assert.NotNil(t, errors)

	// Test min
	ts = &TestStruct{Required: "ok", Min: "ab"}
	errors = v.Validate(ts)
	assert.NotNil(t, errors)

	// Test max
	ts = &TestStruct{Required: "ok", Min: "abcde", Max: "toolong"}
	errors = v.Validate(ts)
	assert.NotNil(t, errors)

	// Test email
	ts = &TestStruct{Required: "ok", Min: "abcde", Max: "ok", Email: "invalid"}
	errors = v.Validate(ts)
	assert.NotNil(t, errors)

	// Test url
	ts = &TestStruct{Required: "ok", Min: "abcde", Max: "ok", Email: "test@test.com", URL: "invalid"}
	errors = v.Validate(ts)
	assert.NotNil(t, errors)

	// Test uuid4
	ts = &TestStruct{Required: "ok", Min: "abcde", Max: "ok", Email: "test@test.com", URL: "http://test.com", UUID: "invalid"}
	errors = v.Validate(ts)
	assert.NotNil(t, errors)
}

func TestValidateUUID4_EmptyString(t *testing.T) {
	v := NewValidator()

	type TestStruct struct {
		UUID string `validate:"uuid4"`
	}

	// Empty string should be valid (uuid4 allows empty)
	ts := &TestStruct{UUID: ""}
	errors := v.Validate(ts)
	assert.Nil(t, errors)
}

func TestValidateUUID4_ValidUUID(t *testing.T) {
	v := NewValidator()

	type TestStruct struct {
		UUID string `validate:"uuid4"`
	}

	ts := &TestStruct{UUID: "550e8400-e29b-41d4-a716-446655440000"}
	errors := v.Validate(ts)
	assert.Nil(t, errors)
}

func TestValidateUUID4_InvalidUUID(t *testing.T) {
	v := NewValidator()

	type TestStruct struct {
		UUID string `validate:"uuid4"`
	}

	ts := &TestStruct{UUID: "not-a-valid-uuid"}
	errors := v.Validate(ts)
	assert.NotNil(t, errors)
}
