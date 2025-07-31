package bobarista

import "github.com/charmbracelet/huh"

// Form represents a single form in the Bobarista form flow.
// It contains the form's metadata, generator function, and behavior callbacks.
type Form struct {
	// ID is a unique identifier for the form within the flow.
	ID string

	// Name is a human-readable name displayed in the UI.
	Name string

	// Group is an optional grouping identifier for organizing related forms.
	Group string

	// Generator creates the actual huh.Form instance when the form is displayed.
	Generator FormGenerator

	// OnComplete is called when the form is successfully completed.
	OnComplete CompletionHandler

	// ShouldSkip determines whether this form should be skipped based on current data.
	ShouldSkip SkipCondition

	// NextForm provides custom navigation logic to determine the next form.
	NextForm NavigationHandler

	// ShowStatus controls whether this form shows progress status in the UI.
	ShowStatus bool
}

// FormGenerator is a function that creates a huh.Form instance.
// It receives the current form's values and global values to customize the form.
type FormGenerator func(current *FormValues, global *FormValues) *huh.Form

// CompletionHandler is called when a form is completed successfully.
// It receives the current form's data and global data, and can return an error to halt the flow.
type CompletionHandler func(current *FormData, global *FormData) error

// SkipCondition determines whether a form should be skipped.
// It receives the current form's data and global data, returning true to skip the form.
type SkipCondition func(current *FormData, global *FormData) bool

// NavigationHandler provides custom navigation logic for a form.
// It receives the current form's data and returns the index of the next form to display.
// Return -1 to use default navigation, -2 to complete the flow.
type NavigationHandler func(current *FormData) int

// NewForm creates a new Form with the specified ID and name.
// The form is created with default settings (ShowStatus = true).
func NewForm(id, name string) Form {
	return Form{
		ID:         id,
		Name:       name,
		ShowStatus: true,
	}
}

// WithGenerator sets the form generator function.
// The generator is responsible for creating the actual huh.Form instance.
func (f Form) WithGenerator(gen FormGenerator) Form {
	f.Generator = gen
	return f
}

// WithOnComplete sets the completion handler for the form.
// This handler is called when the form is successfully completed.
func (f Form) WithOnComplete(handler CompletionHandler) Form {
	f.OnComplete = handler
	return f
}

// WithSkipCondition sets the skip condition for the form.
// The form will be skipped if the condition returns true.
func (f Form) WithSkipCondition(condition SkipCondition) Form {
	f.ShouldSkip = condition
	return f
}

// WithNavigation sets a custom navigation handler for the form.
// This allows for non-linear form flows based on form data.
func (f Form) WithNavigation(handler NavigationHandler) Form {
	f.NextForm = handler
	return f
}

// WithoutStatus disables the progress status display for this form.
// The form will not show progress information in the UI.
func (f Form) WithoutStatus() Form {
	f.ShowStatus = false
	return f
}
