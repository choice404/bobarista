package bobarista

// Recipe holds the configuration settings for a Bobarista form flow.
// It defines the appearance, behavior, and callback functions for the entire flow.
type Recipe struct {
	// Title is the main title displayed at the top of the form flow.
	Title string

	// MaxWidth sets the maximum width in characters for the form display.
	MaxWidth int

	// DisplayKeys specifies which form values to show in the completion summary.
	// If empty, all non-empty values will be displayed.
	DisplayKeys []string

	// ColorScheme determines the visual theme for the form flow.
	// Available options: "default", "dark", "ubuntu", "ocean", "forest", "sunset", "monochrome".
	ColorScheme string

	// Debug enables debug mode, showing additional information during form flow execution.
	Debug bool

	// OnInit is called when the form flow initializes.
	// It receives the Bobarista instance and initial form data for all forms.
	OnInit func(*Bobarista, []FormData)

	// OnComplete is called when the entire form flow completes successfully.
	// It receives the Bobarista instance with all collected data.
	OnComplete func(*Bobarista) error

	// DisplayCallback provides custom content for the completion screen.
	// If nil, a default summary will be generated based on DisplayKeys.
	DisplayCallback func() string
}

// DefaultConfig returns a Recipe with sensible default values.
// This provides a starting point for customization.
func DefaultConfig() Recipe {
	return Recipe{
		Title:       "Cup Sleeve Form Flow",
		MaxWidth:    180,
		ColorScheme: "default",
	}
}
