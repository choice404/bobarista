package bobarista

import "github.com/charmbracelet/lipgloss"

// Styles contains all the styling definitions used throughout the Bobarista form flow.
// It provides consistent visual appearance across different components and states.
type Styles struct {
	// Base is the default style applied to the main content area.
	Base lipgloss.Style

	// HeaderText styles the main header text at the top of the interface.
	HeaderText lipgloss.Style

	// Status styles the status display area, typically used for completion screens.
	Status lipgloss.Style

	// StatusHeader styles headers within status displays.
	StatusHeader lipgloss.Style

	// Error styles error messages and error display areas.
	Error lipgloss.Style

	// ErrorHeader styles headers in error displays.
	ErrorHeader lipgloss.Style

	// Help styles help text and instructions.
	Help lipgloss.Style

	// Highlight styles emphasized or highlighted text.
	Highlight lipgloss.Style

	// KeyText styles form field names and keys in debug displays.
	KeyText lipgloss.Style

	// ValueText styles form field values and data in debug displays.
	ValueText lipgloss.Style

	// FooterText styles footer text and status information.
	FooterText lipgloss.Style

	// Progress styles progress indicators and progress text.
	Progress lipgloss.Style

	// ProgressFilled styles the filled portion of progress indicators.
	ProgressFilled lipgloss.Style

	// Border styles borders around panels and containers.
	Border lipgloss.Style

	// Success styles success messages and positive feedback.
	Success lipgloss.Style

	// Warning styles warning messages and cautionary text.
	Warning lipgloss.Style

	// Info styles informational messages and neutral feedback.
	Info lipgloss.Style
}

// NewStyles creates a new Styles instance with the specified color scheme.
// It initializes all style components with appropriate colors and formatting.
func NewStyles(colorScheme ColorScheme) *Styles {
	return &Styles{
		Base: lipgloss.NewStyle().
			Padding(1, 2),

		HeaderText: lipgloss.NewStyle().
			Foreground(colorScheme.Primary).
			Bold(true).
			Padding(0, 1, 0, 2),

		Status: lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(colorScheme.Primary).
			Padding(1, 2).
			Margin(1, 0),

		StatusHeader: lipgloss.NewStyle().
			Foreground(colorScheme.Secondary).
			Bold(true),

		Error: lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(colorScheme.Error).
			Foreground(colorScheme.Error).
			Padding(1, 2).
			Margin(1, 0),

		ErrorHeader: lipgloss.NewStyle().
			Foreground(colorScheme.Error).
			Bold(true).
			Padding(0, 1, 0, 2),

		Help: lipgloss.NewStyle().
			Foreground(lipgloss.Color("240")),

		Highlight: lipgloss.NewStyle().
			Foreground(colorScheme.Secondary).
			Bold(true),

		KeyText: lipgloss.NewStyle().
			Foreground(colorScheme.Tertiary).
			Bold(true),

		ValueText: lipgloss.NewStyle().
			Foreground(colorScheme.Secondary),

		FooterText: lipgloss.NewStyle().
			Foreground(lipgloss.Color("240")).
			Bold(true).
			Padding(0, 1, 0, 2),

		Progress: lipgloss.NewStyle().
			Foreground(colorScheme.Tertiary),

		ProgressFilled: lipgloss.NewStyle().
			Foreground(colorScheme.Primary).
			Bold(true),

		Border: lipgloss.NewStyle().
			BorderForeground(colorScheme.Primary),

		Success: lipgloss.NewStyle().
			Foreground(colorScheme.Success).
			Bold(true),

		Warning: lipgloss.NewStyle().
			Foreground(colorScheme.Warning).
			Bold(true),

		Info: lipgloss.NewStyle().
			Foreground(colorScheme.Info),
	}
}

// DefaultStyles creates a new Styles instance using the default color scheme.
// This provides a convenient way to get standard styling without specifying a scheme.
func DefaultStyles() *Styles {
	defaultScheme, _ := GetColorScheme("default")
	return NewStyles(defaultScheme)
}

// ApplyColorScheme updates all styles to use the specified color scheme.
// This allows for dynamic theme switching without recreating the entire Styles instance.
func (s *Styles) ApplyColorScheme(colorScheme ColorScheme) {
	s.HeaderText = s.HeaderText.Foreground(colorScheme.Primary)
	s.Status = s.Status.BorderForeground(colorScheme.Primary)
	s.StatusHeader = s.StatusHeader.Foreground(colorScheme.Secondary)
	s.Error = s.Error.BorderForeground(colorScheme.Error).Foreground(colorScheme.Error)
	s.ErrorHeader = s.ErrorHeader.Foreground(colorScheme.Error)
	s.Highlight = s.Highlight.Foreground(colorScheme.Secondary)
	s.KeyText = s.KeyText.Foreground(colorScheme.Tertiary)
	s.ValueText = s.ValueText.Foreground(colorScheme.Secondary)
	s.Progress = s.Progress.Foreground(colorScheme.Tertiary)
	s.ProgressFilled = s.ProgressFilled.Foreground(colorScheme.Primary)
	s.Border = s.Border.BorderForeground(colorScheme.Primary)
	s.Success = s.Success.Foreground(colorScheme.Success)
	s.Warning = s.Warning.Foreground(colorScheme.Warning)
	s.Info = s.Info.Foreground(colorScheme.Info)
}
