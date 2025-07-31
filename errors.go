package bobarista

import (
	"errors"
	"fmt"
	"strings"
)

// Predefined errors for common Bobarista operations.
var (
	// ErrInvalidFormIndex is returned when an invalid form index is provided.
	ErrInvalidFormIndex = errors.New("invalid form index")

	// ErrFormNotFound is returned when a requested form cannot be found.
	ErrFormNotFound = errors.New("form not found")

	// ErrNoFormsProvided is returned when no forms are provided to the navigator.
	ErrNoFormsProvided = errors.New("no forms provided")

	// ErrNoValidForms is returned when no valid forms are found for navigation.
	ErrNoValidForms = errors.New("no valid forms found")

	// ErrNoGenerator is returned when a form lacks a required generator function.
	ErrNoGenerator = errors.New("form generator is required")

	// ErrNilForm is returned when a form generator returns nil.
	ErrNilForm = errors.New("form generator returned nil")

	// ErrEmptyFormID is returned when a form has an empty ID.
	ErrEmptyFormID = errors.New("form ID cannot be empty")
)

// CupSleeveError represents an error that occurred within a specific form.
// It wraps the underlying error with context about which form caused the issue.
type CupSleeveError struct {
	// FormID identifies the form where the error occurred.
	FormID string
	// Err is the underlying error.
	Err error
}

// Error implements the error interface, providing a formatted error message.
// It includes the form ID if available.
func (fe CupSleeveError) Error() string {
	if fe.FormID != "" {
		return "form '" + fe.FormID + "': " + fe.Err.Error()
	}
	return fe.Err.Error()
}

// Unwrap returns the underlying error, supporting Go's error unwrapping.
func (fe CupSleeveError) Unwrap() error {
	return fe.Err
}

// NewCupSleeveError creates a new CupSleeveError with the specified form ID and underlying error.
func NewCupSleeveError(formID string, err error) CupSleeveError {
	return CupSleeveError{
		FormID: formID,
		Err:    err,
	}
}

// DuplicateFormIDError represents an error when duplicate form IDs are detected.
// This helps identify configuration issues during form flow setup.
type DuplicateFormIDError struct {
	// ID is the duplicate form ID that was found.
	ID string
	// FirstIndex is the index of the first occurrence of the duplicate ID.
	FirstIndex int
	// SecondIndex is the index of the second occurrence of the duplicate ID.
	SecondIndex int
}

// Error implements the error interface for DuplicateFormIDError.
func (e DuplicateFormIDError) Error() string {
	return fmt.Sprintf("duplicate form ID '%s' found at indices %d and %d",
		e.ID, e.FirstIndex, e.SecondIndex)
}

// NewDuplicateFormIDError creates a new DuplicateFormIDError with the specified details.
func NewDuplicateFormIDError(id string, firstIndex, secondIndex int) DuplicateFormIDError {
	return DuplicateFormIDError{
		ID:          id,
		FirstIndex:  firstIndex,
		SecondIndex: secondIndex,
	}
}

// ValidationError represents a validation error that occurred in a specific form field.
// It provides context about both the form and the specific field that failed validation.
type ValidationError struct {
	// FormID identifies the form where validation failed.
	FormID string
	// Field identifies the specific field that failed validation.
	Field string
	// Err is the underlying validation error.
	Err error
}

// Error implements the error interface for ValidationError.
func (ve ValidationError) Error() string {
	if ve.Field != "" {
		return fmt.Sprintf("validation error in form '%s', field '%s': %s",
			ve.FormID, ve.Field, ve.Err.Error())
	}
	return fmt.Sprintf("validation error in form '%s': %s", ve.FormID, ve.Err.Error())
}

// Unwrap returns the underlying error, supporting Go's error unwrapping.
func (ve ValidationError) Unwrap() error {
	return ve.Err
}

// NewValidationError creates a new ValidationError with the specified form ID, field, and underlying error.
func NewValidationError(formID, field string, err error) ValidationError {
	return ValidationError{
		FormID: formID,
		Field:  field,
		Err:    err,
	}
}

// NavigationError represents an error that occurred during form navigation.
// It provides context about the source and destination forms.
type NavigationError struct {
	// FromFormID identifies the form being navigated from.
	FromFormID string
	// ToFormID identifies the form being navigated to.
	ToFormID string
	// Err is the underlying navigation error.
	Err error
}

// Error implements the error interface for NavigationError.
func (ne NavigationError) Error() string {
	return fmt.Sprintf("navigation error from '%s' to '%s': %s",
		ne.FromFormID, ne.ToFormID, ne.Err.Error())
}

// Unwrap returns the underlying error, supporting Go's error unwrapping.
func (ne NavigationError) Unwrap() error {
	return ne.Err
}

// NewNavigationError creates a new NavigationError with the specified form IDs and underlying error.
func NewNavigationError(fromFormID, toFormID string, err error) NavigationError {
	return NavigationError{
		FromFormID: fromFormID,
		ToFormID:   toFormID,
		Err:        err,
	}
}

// ErrorCollector provides a convenient way to collect multiple errors during form processing.
// It's useful for validation and setup phases where multiple errors might occur.
type ErrorCollector struct {
	errors []error
}

// NewErrorCollector creates a new ErrorCollector instance.
func NewErrorCollector() *ErrorCollector {
	return &ErrorCollector{
		errors: make([]error, 0),
	}
}

// Add adds an error to the collector if it's not nil.
func (ec *ErrorCollector) Add(err error) {
	if err != nil {
		ec.errors = append(ec.errors, err)
	}
}

// AddCupSleeveError adds a CupSleeveError to the collector with the specified form ID.
func (ec *ErrorCollector) AddCupSleeveError(formID string, err error) {
	if err != nil {
		ec.errors = append(ec.errors, NewCupSleeveError(formID, err))
	}
}

// HasErrors returns true if the collector contains any errors.
func (ec *ErrorCollector) HasErrors() bool {
	return len(ec.errors) > 0
}

// Errors returns all collected errors.
func (ec *ErrorCollector) Errors() []error {
	return ec.errors
}

// Error implements the error interface, providing a formatted summary of all collected errors.
func (ec *ErrorCollector) Error() string {
	if len(ec.errors) == 0 {
		return "no errors"
	}

	if len(ec.errors) == 1 {
		return ec.errors[0].Error()
	}

	var result strings.Builder
	result.WriteString(fmt.Sprintf("%d errors occurred:\n", len(ec.errors)))

	for i, err := range ec.errors {
		result.WriteString(fmt.Sprintf("  %d. %s\n", i+1, err.Error()))
	}

	return result.String()
}

// Clear removes all errors from the collector, resetting it to an empty state.
func (ec *ErrorCollector) Clear() {
	ec.errors = make([]error, 0)
}
