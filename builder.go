// Package bobarista provides a flexible form flow framework built on top of Bubble Tea.
// It allows you to create multi-step forms with navigation, validation, and custom styling.
package bobarista

// BobaBuilder provides a fluent interface for constructing Bobarista form flows.
// It allows you to configure forms, styling, and behavior before building the final Bobarista instance.
type BobaBuilder struct {
	config Recipe
	forms  []Form
}

// AddForm adds a form to the form flow.
// Forms will be processed in the order they are added.
func (b *BobaBuilder) AddForm(form Form) *BobaBuilder {
	b.forms = append(b.forms, form)
	return b
}

// WithMaxWidth sets the maximum width for the form display.
// This controls how wide the form can be rendered on screen.
func (b *BobaBuilder) WithMaxWidth(width int) *BobaBuilder {
	b.config.MaxWidth = width
	return b
}

// WithColorScheme sets the color scheme for the form display.
// Available schemes include "default", "dark", "ubuntu", "ocean", "forest", "sunset", and "monochrome".
func (b *BobaBuilder) WithColorScheme(scheme string) *BobaBuilder {
	b.config.ColorScheme = scheme
	return b
}

// WithDisplayKeys sets which keys should be displayed in the completion summary.
// If not set, all non-empty values will be displayed.
func (b *BobaBuilder) WithDisplayKeys(keys []string) *BobaBuilder {
	b.config.DisplayKeys = keys
	return b
}

// OnInit sets a callback function that is called when the form flow initializes.
// The callback receives the Bobarista instance and initial form data for all forms.
func (b *BobaBuilder) OnInit(handler func(*Bobarista, []FormData)) *BobaBuilder {
	b.config.OnInit = handler
	return b
}

// OnComplete sets a callback function that is called when the entire form flow completes.
// The callback receives the Bobarista instance with all collected data.
func (b *BobaBuilder) OnComplete(handler func(*Bobarista) error) *BobaBuilder {
	b.config.OnComplete = handler
	return b
}

// WithDisplayCallback sets a custom function to generate the completion display content.
// If not set, a default summary will be shown based on DisplayKeys or all values.
func (b *BobaBuilder) WithDisplayCallback(callback func() string) *BobaBuilder {
	b.config.DisplayCallback = callback
	return b
}

// WithDebug enables or disables debug mode.
// When enabled, a debug panel shows current form state, values, and navigation info.
func (b *BobaBuilder) WithDebug(enabled bool) *BobaBuilder {
	b.config.Debug = enabled
	return b
}

// Build creates and returns a new Bobarista instance with the configured settings.
// This finalizes the builder and creates the form flow ready for execution.
func (b *BobaBuilder) Build() *Bobarista {
	if b.config.ColorScheme == "" {
		b.config.ColorScheme = "default"
	}

	return &Bobarista{
		config:      b.config,
		forms:       b.forms,
		navigator:   NewNavigator(b.forms),
		renderer:    NewRenderer(b.config),
		globalData:  &FormData{ID: "global", Values: NewFormValues()},
		currentForm: nil,
		state:       StateActive,
		errors:      make([]error, 0),
	}
}
