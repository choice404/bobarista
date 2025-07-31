package bobarista

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/huh"
	"github.com/charmbracelet/lipgloss"
	"github.com/choice404/bobarista/internal"
)

// Renderer handles the visual presentation of the Bobarista form flow.
// It manages layout, styling, and different display states (active, completed, error).
type Renderer struct {
	config   Recipe
	viewport *internal.Viewport
	styles   *Styles
	width    int
	height   int
}

// NewRenderer creates a new Renderer with the specified configuration.
// It initializes the viewport, applies the color scheme, and sets default dimensions.
func NewRenderer(config Recipe) *Renderer {
	colorScheme, exists := GetColorScheme(config.ColorScheme)
	if !exists {
		colorScheme, _ = GetColorScheme("default")
	}

	return &Renderer{
		config:   config,
		viewport: internal.NewViewport(),
		styles:   NewStyles(colorScheme),
		width:    config.MaxWidth,
		height:   24,
	}
}

// UpdateSize updates the renderer's dimensions based on terminal size.
// It respects the maximum width configuration and updates the viewport accordingly.
func (r *Renderer) UpdateSize(width, height int) {
	if width < r.config.MaxWidth {
		r.width = width - r.styles.Base.GetHorizontalFrameSize()
	} else {
		r.width = r.config.MaxWidth - r.styles.Base.GetHorizontalFrameSize()
	}
	r.height = height - 4
	r.viewport.SetSize(r.width, r.height)
}

// Render renders the Bobarista form flow based on its current state.
// It delegates to specific render methods based on the application state.
func (r *Renderer) Render(cupSleeve *Bobarista) string {
	switch cupSleeve.state {
	case StateError:
		return r.renderError(cupSleeve)
	case StateCompleted:
		return r.renderCompleted(cupSleeve)
	default:
		return r.renderActive(cupSleeve)
	}
}

// renderActive renders the form flow in its active state.
// It displays the current form with header, content, and footer.
func (r *Renderer) renderActive(cupSleeve *Bobarista) string {
	current := cupSleeve.navigator.Current()
	if current == nil {
		return r.styles.Base.Render("No forms available")
	}

	progress := cupSleeve.navigator.GetProgress()
	title := fmt.Sprintf("%s - %s (%.0f%%)", cupSleeve.config.Title, current.Name, progress)
	header := r.renderHeader(title)

	var mainContent string
	if cupSleeve.config.Debug {
		mainContent = r.renderWithDebugPanel(cupSleeve)
	} else {
		mainContent = r.renderFormOnly(cupSleeve)
	}

	var footerText string
	if cupSleeve.currentForm != nil && len(cupSleeve.currentForm.Errors()) > 0 {
		footerText = r.renderErrors(cupSleeve.currentForm.Errors())
	} else {
		if cupSleeve.config.Debug {
			footerText = "Press Ctrl+C to quit • Debug mode enabled"
		} else {
			footerText = "Press Ctrl+C to quit"
		}
	}
	footer := r.renderFooter(footerText)

	return lipgloss.JoinVertical(lipgloss.Left, header, mainContent, footer)
}

// renderWithDebugPanel renders the form with a debug panel showing internal state.
// The layout splits the available width between the form and debug information.
func (r *Renderer) renderWithDebugPanel(cupSleeve *Bobarista) string {
	// Split width: 60% for form, 40% for debug panel
	formWidth := int(float64(r.width) * 0.6)
	debugWidth := r.width - formWidth - 2

	var formContent string
	if cupSleeve.currentForm != nil {
		formView := cupSleeve.currentForm.View()
		formContent = lipgloss.NewStyle().
			Width(formWidth).
			Render(formView)
	} else {
		formContent = lipgloss.NewStyle().
			Width(formWidth).
			Render("Loading form...")
	}

	debugContent := r.renderDebugPanel(cupSleeve, debugWidth)

	return lipgloss.JoinHorizontal(
		lipgloss.Top,
		formContent,
		lipgloss.NewStyle().Width(2).Render(""),
		debugContent,
	)
}

// renderFormOnly renders just the form without any debug information.
func (r *Renderer) renderFormOnly(cupSleeve *Bobarista) string {
	if cupSleeve.currentForm != nil {
		formView := cupSleeve.currentForm.View()
		return r.styles.Base.Render(formView)
	}
	return r.styles.Base.Render("Loading form...")
}

// renderDebugPanel creates a debug panel showing form state, values, and navigation info.
// This is displayed when debug mode is enabled.
func (r *Renderer) renderDebugPanel(cupSleeve *Bobarista, width int) string {
	current := cupSleeve.navigator.Current()
	if current == nil {
		return r.styles.Base.Width(width).Render("No current form")
	}

	var content strings.Builder

	content.WriteString(r.styles.Highlight.Render("🐛 DEBUG PANEL"))
	content.WriteString("\n\n")

	content.WriteString(r.styles.KeyText.Render("Current Form:"))
	content.WriteString("\n")
	content.WriteString(fmt.Sprintf("  ID: %s\n", r.styles.ValueText.Render(current.ID)))
	content.WriteString(fmt.Sprintf("  Name: %s\n", r.styles.ValueText.Render(current.Name)))
	content.WriteString(fmt.Sprintf("  Index: %s\n", r.styles.ValueText.Render(fmt.Sprintf("%d/%d",
		cupSleeve.navigator.GetCurrentIndex()+1, cupSleeve.navigator.GetFormCount()))))
	content.WriteString("\n")

	content.WriteString(r.styles.KeyText.Render("Current Form Values:"))
	content.WriteString("\n")
	currentData := cupSleeve.GetCurrentFormData()
	if len(*currentData.Values) == 0 {
		content.WriteString("  " + r.styles.Help.Render("(empty)") + "\n")
	} else {
		for key, valuePtr := range *currentData.Values {
			value := "(nil)"
			if valuePtr != nil {
				value = *valuePtr
				if value == "" {
					value = "(empty)"
				}
			}
			content.WriteString(fmt.Sprintf("  %s: %s\n",
				r.styles.KeyText.Render(key),
				r.styles.ValueText.Render(value)))
		}
	}
	content.WriteString("\n")

	content.WriteString(r.styles.KeyText.Render("Global Values:"))
	content.WriteString("\n")
	globalData := cupSleeve.GetGlobalData()
	if len(*globalData.Values) == 0 {
		content.WriteString("  " + r.styles.Help.Render("(empty)") + "\n")
	} else {
		for key, valuePtr := range *globalData.Values {
			value := "(nil)"
			if valuePtr != nil {
				value = *valuePtr
				if value == "" {
					value = "(empty)"
				}
			}
			content.WriteString(fmt.Sprintf("  %s: %s\n",
				r.styles.KeyText.Render(key),
				r.styles.ValueText.Render(value)))
		}
	}
	content.WriteString("\n")

	content.WriteString(r.styles.KeyText.Render("Navigation:"))
	content.WriteString("\n")
	content.WriteString(fmt.Sprintf("  Has Previous: %s\n",
		r.styles.ValueText.Render(fmt.Sprintf("%t", cupSleeve.navigator.HasPrevious()))))
	content.WriteString(fmt.Sprintf("  Has Next: %s\n",
		r.styles.ValueText.Render(fmt.Sprintf("%t", cupSleeve.navigator.HasNext()))))
	content.WriteString(fmt.Sprintf("  Progress: %s\n",
		r.styles.ValueText.Render(fmt.Sprintf("%.1f%%", cupSleeve.navigator.GetProgress()))))
	content.WriteString("\n")

	content.WriteString(r.styles.KeyText.Render("Form State:"))
	content.WriteString("\n")
	if cupSleeve.currentForm != nil {
		content.WriteString(fmt.Sprintf("  State: %s\n",
			r.styles.ValueText.Render(r.formStateToString(cupSleeve.currentForm.State))))
		content.WriteString(fmt.Sprintf("  Errors: %s\n",
			r.styles.ValueText.Render(fmt.Sprintf("%d", len(cupSleeve.currentForm.Errors())))))
	} else {
		content.WriteString("  " + r.styles.Help.Render("No current form") + "\n")
	}

	debugPanel := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(r.styles.KeyText.GetForeground()).
		Padding(1).
		Width(width).
		Height(r.height - 4).
		Render(content.String())

	return debugPanel
}

// formStateToString converts a huh.FormState to a human-readable string.
func (r *Renderer) formStateToString(state huh.FormState) string {
	switch state {
	case huh.StateNormal:
		return "Normal"
	case huh.StateCompleted:
		return "Completed"
	case huh.StateAborted:
		return "Aborted"
	default:
		return fmt.Sprintf("Unknown(%d)", int(state))
	}
}

// renderCompleted renders the completion screen with collected data.
// It supports scrolling for long content and custom display callbacks.
func (r *Renderer) renderCompleted(cupSleeve *Bobarista) string {
	header := r.renderHeader(cupSleeve.config.Title + " - Completed")

	var content string
	if cupSleeve.config.DisplayCallback != nil {
		content = cupSleeve.config.DisplayCallback()
	} else {
		content = r.renderDefaultCompletion(cupSleeve)
	}

	lines := strings.Split(content, "\n")
	r.viewport.SetContent(lines)

	visibleLines := r.viewport.VisibleContent()
	body := r.styles.Status.Render(strings.Join(visibleLines, "\n"))

	var footerText string
	if r.viewport.CanScrollUp() || r.viewport.CanScrollDown() {
		footerText = "↑/↓ to scroll, Enter to finish, Q to quit"
	} else {
		footerText = "Press Enter to finish, Q to quit"
	}
	footer := r.renderFooter(footerText)

	return lipgloss.JoinVertical(lipgloss.Left, header, body, footer)
}

// renderError renders the error screen displaying all accumulated errors.
func (r *Renderer) renderError(cupSleeve *Bobarista) string {
	header := r.renderHeader(cupSleeve.config.Title + " - Error")

	var content strings.Builder
	content.WriteString("The following errors occurred:\n\n")

	for i, err := range cupSleeve.errors {
		if i > 0 {
			content.WriteString("\n")
		}

		if cupSleeveErr, ok := err.(CupSleeveError); ok && cupSleeveErr.FormID != "" {
			content.WriteString(fmt.Sprintf("Form '%s': %s", cupSleeveErr.FormID, cupSleeveErr.Err.Error()))
		} else {
			content.WriteString(err.Error())
		}
	}

	body := r.styles.Error.Render(content.String())
	footer := r.renderFooter("Press Q to quit")

	return lipgloss.JoinVertical(lipgloss.Left, header, body, footer)
}

// renderDefaultCompletion creates the default completion display.
// It shows either specified display keys or all collected values.
func (r *Renderer) renderDefaultCompletion(cupSleeve *Bobarista) string {
	var content strings.Builder

	globalData := cupSleeve.GetGlobalData()

	if len(r.config.DisplayKeys) > 0 {
		// Show only specified keys
		content.WriteString("Summary:\n\n")
		for _, key := range r.config.DisplayKeys {
			if value, exists := globalData.Values.Get(key); exists && value != "" {
				content.WriteString(fmt.Sprintf("%s: %s\n",
					r.formatKey(key), r.styles.Highlight.Render(value)))
			}
		}
	} else {
		// Show all non-empty values
		content.WriteString("All Values:\n\n")
		hasValues := false
		for key, valuePtr := range *globalData.Values {
			if valuePtr != nil && *valuePtr != "" {
				content.WriteString(fmt.Sprintf("%s: %s\n",
					r.formatKey(key), r.styles.Highlight.Render(*valuePtr)))
				hasValues = true
			}
		}

		if !hasValues {
			content.WriteString("No values to display.")
		}
	}

	return content.String()
}

// renderHeader creates a styled header with the specified title.
func (r *Renderer) renderHeader(title string) string {
	return lipgloss.PlaceHorizontal(
		r.width,
		lipgloss.Left,
		r.styles.HeaderText.Render(title),
		lipgloss.WithWhitespaceChars("="),
		lipgloss.WithWhitespaceForeground(lipgloss.Color("12")),
	)
}

// renderFooter creates a styled footer with help text.
func (r *Renderer) renderFooter(text string) string {
	return r.styles.Help.Render(text)
}

// renderErrors formats multiple errors for display in the footer.
func (r *Renderer) renderErrors(errors []error) string {
	if len(errors) == 0 {
		return ""
	}

	var content strings.Builder
	content.WriteString("Errors: ")

	for i, err := range errors {
		if i > 0 {
			content.WriteString(", ")
		}
		content.WriteString(err.Error())
	}

	return r.styles.Error.Render(content.String())
}

// formatKey converts underscore-separated keys to human-readable format.
// For example, "first_name" becomes "First Name".
func (r *Renderer) formatKey(key string) string {
	if key == "" {
		return ""
	}

	parts := strings.Split(key, "_")
	for i, part := range parts {
		if len(part) > 0 {
			parts[i] = strings.ToUpper(part[:1]) + strings.ToLower(part[1:])
		}
	}

	return strings.Join(parts, " ")
}

// HandleScroll processes scroll events for the viewport.
// Positive direction scrolls down, negative scrolls up.
func (r *Renderer) HandleScroll(direction int) {
	r.viewport.Scroll(direction)
}

// CanScrollUp returns true if the viewport can scroll up.
func (r *Renderer) CanScrollUp() bool {
	return r.viewport.CanScrollUp()
}

// CanScrollDown returns true if the viewport can scroll down.
func (r *Renderer) CanScrollDown() bool {
	return r.viewport.CanScrollDown()
}

// SetColorScheme updates the renderer's color scheme.
// If the scheme doesn't exist, the change is ignored.
func (r *Renderer) SetColorScheme(schemeName string) {
	if colorScheme, exists := GetColorScheme(schemeName); exists {
		r.styles.ApplyColorScheme(colorScheme)
	}
}
