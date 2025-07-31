package bobarista

import (
	"github.com/charmbracelet/lipgloss"
)

// ColorScheme defines a complete color theme for the Bobarista form flow.
// It provides colors for different UI elements and semantic meanings.
type ColorScheme struct {
	// Name is the human-readable name of the color scheme.
	Name string

	// Primary is the main accent color used for headers and primary elements.
	Primary lipgloss.AdaptiveColor

	// Secondary is used for highlighted text and secondary emphasis.
	Secondary lipgloss.AdaptiveColor

	// Tertiary is used for labels, keys, and tertiary elements.
	Tertiary lipgloss.AdaptiveColor

	// Success is used for positive feedback and success messages.
	Success lipgloss.AdaptiveColor

	// Error is used for error messages and negative feedback.
	Error lipgloss.AdaptiveColor

	// Warning is used for warning messages and cautionary feedback.
	Warning lipgloss.AdaptiveColor

	// Info is used for informational messages and neutral feedback.
	Info lipgloss.AdaptiveColor
}

// Predefined adaptive colors that work well in both light and dark terminals.
var (
	red      = lipgloss.AdaptiveColor{Light: "#FE5F86", Dark: "#FE5F86"}
	indigo   = lipgloss.AdaptiveColor{Light: "#5A56E0", Dark: "#7571F9"}
	green    = lipgloss.AdaptiveColor{Light: "#02BA84", Dark: "#02BF87"}
	cyan     = lipgloss.AdaptiveColor{Light: "#00FFFF", Dark: "#0066aa"}
	white    = lipgloss.AdaptiveColor{Light: "#FFFFFF", Dark: "#FFFFFF"}
	blue     = lipgloss.AdaptiveColor{Light: "#0077FF", Dark: "#0000FF"}
	navy     = lipgloss.AdaptiveColor{Light: "#000080", Dark: "#000080"}
	sapphire = lipgloss.AdaptiveColor{Light: "#0F52BA", Dark: "#0F52BA"}
	emerald  = lipgloss.AdaptiveColor{Light: "#50C878", Dark: "#50C878"}
	orange   = lipgloss.AdaptiveColor{Light: "#FF8C00", Dark: "#FF8C00"}
	purple   = lipgloss.AdaptiveColor{Light: "#9932CC", Dark: "#9932CC"}
	yellow   = lipgloss.AdaptiveColor{Light: "#FFD700", Dark: "#FFD700"}
	gray     = lipgloss.AdaptiveColor{Light: "#808080", Dark: "#A0A0A0"}
)

// colorSchemes contains all built-in color schemes available in Bobarista.
var colorSchemes = map[string]ColorScheme{
	"default": {
		Name:      "Default",
		Primary:   cyan,
		Secondary: emerald,
		Tertiary:  indigo,
		Success:   green,
		Error:     red,
		Warning:   orange,
		Info:      blue,
	},
	"dark": {
		Name:      "Dark",
		Primary:   white,
		Secondary: sapphire,
		Tertiary:  emerald,
		Success:   green,
		Error:     red,
		Warning:   orange,
		Info:      cyan,
	},
	"ubuntu": {
		Name:      "Ubuntu",
		Primary:   lipgloss.AdaptiveColor{Light: "#E95420", Dark: "#E95420"},
		Secondary: lipgloss.AdaptiveColor{Light: "#D3A625", Dark: "#D3A625"},
		Tertiary:  indigo,
		Success:   green,
		Error:     red,
		Warning:   orange,
		Info:      blue,
	},
	"ocean": {
		Name:      "Ocean",
		Primary:   blue,
		Secondary: cyan,
		Tertiary:  navy,
		Success:   emerald,
		Error:     red,
		Warning:   orange,
		Info:      sapphire,
	},
	"forest": {
		Name:      "Forest",
		Primary:   emerald,
		Secondary: green,
		Tertiary:  lipgloss.AdaptiveColor{Light: "#228B22", Dark: "#32CD32"},
		Success:   green,
		Error:     red,
		Warning:   yellow,
		Info:      blue,
	},
	"sunset": {
		Name:      "Sunset",
		Primary:   orange,
		Secondary: red,
		Tertiary:  yellow,
		Success:   green,
		Error:     red,
		Warning:   orange,
		Info:      purple,
	},
	"monochrome": {
		Name:      "Monochrome",
		Primary:   white,
		Secondary: gray,
		Tertiary:  lipgloss.AdaptiveColor{Light: "#404040", Dark: "#C0C0C0"},
		Success:   white,
		Error:     lipgloss.AdaptiveColor{Light: "#000000", Dark: "#FFFFFF"},
		Warning:   gray,
		Info:      white,
	},
}

// GetColorScheme retrieves a color scheme by name.
// Returns the color scheme and true if found, or an empty scheme and false if not found.
func GetColorScheme(name string) (ColorScheme, bool) {
	scheme, exists := colorSchemes[name]
	return scheme, exists
}

// GetAvailableColorSchemes returns a list of all available color scheme names.
// This is useful for providing users with theme selection options.
func GetAvailableColorSchemes() []string {
	schemes := make([]string, 0, len(colorSchemes))
	for name := range colorSchemes {
		schemes = append(schemes, name)
	}
	return schemes
}

// RegisterColorScheme adds a new color scheme to the available schemes.
// This allows users to define and register custom themes.
func RegisterColorScheme(name string, scheme ColorScheme) {
	scheme.Name = name
	colorSchemes[name] = scheme
}

// CreateCustomColorScheme creates a new ColorScheme with the specified colors.
// It uses the provided colors for primary, secondary, and tertiary elements,
// while using standard colors for semantic elements (success, error, warning, info).
func CreateCustomColorScheme(name string, primary, secondary, tertiary string) ColorScheme {
	return ColorScheme{
		Name:      name,
		Primary:   lipgloss.AdaptiveColor{Light: primary, Dark: primary},
		Secondary: lipgloss.AdaptiveColor{Light: secondary, Dark: secondary},
		Tertiary:  lipgloss.AdaptiveColor{Light: tertiary, Dark: tertiary},
		Success:   green,
		Error:     red,
		Warning:   orange,
		Info:      blue,
	}
}
