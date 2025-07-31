package bobarista

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/huh"
)

// Bobarista represents the main form flow application.
// It manages the state, navigation, and rendering of multi-step forms.
type Bobarista struct {
	config      Recipe
	forms       []Form
	navigator   *Navigator
	renderer    *Renderer
	globalData  *FormData
	currentForm *huh.Form
	state       BobaState
	errors      []error
	formValues  map[string]FormValues
}

// BobaState represents the current state of the form flow.
type BobaState int

const (
	// StateActive indicates the form flow is currently active and accepting input.
	StateActive BobaState = iota
	// StateCompleted indicates the form flow has completed successfully.
	StateCompleted
	// StateError indicates the form flow has encountered an error.
	StateError
)

// Run starts the form flow and blocks until completion or error.
// It initializes the Bubble Tea program and handles the main event loop.
func (f *Bobarista) Run() error {
	f.infoLog("Starting Bobarista form flow")
	_, err := tea.NewProgram(f, tea.WithAltScreen()).Run()
	if err != nil {
		f.errorLog(fmt.Errorf("tea program error: %w", err))
	}
	f.infoLog("Bobarista form flow completed")
	return err
}

// New creates a new BobaBuilder with the specified title.
// This is the entry point for creating a new form flow.
func New(title string) *BobaBuilder {
	return &BobaBuilder{
		config: Recipe{
			Title:    title,
			MaxWidth: 180,
		},
		forms: make([]Form, 0),
	}
}

// Init implements the tea.Model interface and initializes the form flow.
// It sets up global data, form values, and navigates to the first valid form.
func (f *Bobarista) Init() tea.Cmd {
	f.infoLog("Initializing Bobarista")

	if f.globalData == nil {
		f.debugLog("Initializing global data")
		f.globalData = &FormData{
			ID:     "global",
			Values: NewFormValues(),
		}
	}

	if f.formValues == nil {
		f.debugLog("Initializing form values map")
		f.formValues = make(map[string]FormValues)
	}

	if f.config.OnInit != nil {
		f.debugLog("Calling OnInit callback")
		formDataList := make([]FormData, len(f.forms))
		for i, form := range f.forms {
			formDataList[i] = FormData{
				ID:     form.ID,
				Values: NewFormValues(),
			}
		}
		f.config.OnInit(f, formDataList)
	}

	f.debugLog("Moving to first valid form")
	if err := f.navigator.MoveToFirstValid(*f.globalData); err != nil {
		f.errorLog(fmt.Errorf("failed to move to first valid form: %w", err))
		f.addError("", err)
		return nil
	}

	return f.initCurrentForm()
}

// Update implements the tea.Model interface and handles incoming messages.
// It processes keyboard input, window resize events, and form state changes.
func (f *Bobarista) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		f.renderer.UpdateSize(msg.Width, msg.Height)
		if f.currentForm != nil {
			form, cmd := f.currentForm.Update(msg)
			if huhForm, ok := form.(*huh.Form); ok {
				f.currentForm = huhForm
			}
			return f, cmd
		}
	case tea.KeyMsg:
		if f.state == StateCompleted || f.state == StateError {
			return f.handleCompletedState(msg)
		}

		switch msg.String() {
		case "ctrl+c":
			f.infoLog("User pressed Ctrl+C, quitting")
			return f, tea.Quit
		case "esc":
			f.infoLog("User pressed Esc, quitting")
			return f, tea.Quit
		}

		if f.state == StateActive && f.currentForm != nil {
			return f.updateCurrentForm(msg)
		}
	}

	if f.state == StateActive && f.currentForm != nil {
		return f.updateCurrentForm(msg)
	}

	return f, nil
}

// View implements the tea.Model interface and renders the current state.
// It delegates to the renderer based on the current form flow state.
func (f *Bobarista) View() string {
	return f.renderer.Render(f)
}

// handleCompletedState processes input when the form flow is in completed or error state.
// It handles scrolling, quitting, and completion callbacks.
func (f *Bobarista) handleCompletedState(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "up":
		f.renderer.HandleScroll(-1)
		return f, nil
	case "down":
		f.renderer.HandleScroll(1)
		return f, nil
	case "q", "esc":
		f.infoLog("User quit from completed/error state")
		return f, tea.Quit
	case "enter":
		if f.state == StateCompleted && f.config.OnComplete != nil {
			f.debugLog("Calling OnComplete callback")
			if err := f.config.OnComplete(f); err != nil {
				f.errorLog(fmt.Errorf("OnComplete callback error: %w", err))
				f.addError("", err)
				return f, nil
			}
		}
		f.infoLog("User finished from completed state")
		return f, tea.Quit
	}
	return f, nil
}

// updateCurrentForm processes messages for the currently active form.
// It handles form completion and transitions to the next form.
func (f *Bobarista) updateCurrentForm(msg tea.Msg) (tea.Model, tea.Cmd) {
	if f.currentForm == nil {
		f.warningLog("updateCurrentForm called with nil currentForm")
		return f, nil
	}

	form, cmd := f.currentForm.Update(msg)
	if huhForm, ok := form.(*huh.Form); ok {
		f.currentForm = huhForm
	}

	if f.currentForm.State == huh.StateCompleted {
		current := f.navigator.Current()
		if current != nil {
			f.infoLog(fmt.Sprintf("Form '%s' completed", current.ID))
		}
		return f.handleFormCompletion()
	}

	return f, cmd
}

// handleFormCompletion processes the completion of the current form.
// It calls completion handlers, merges data, and navigates to the next form.
func (f *Bobarista) handleFormCompletion() (tea.Model, tea.Cmd) {
	current := f.navigator.Current()
	if current == nil {
		f.warningLog("handleFormCompletion called with no current form")
		f.state = StateCompleted
		return f, nil
	}

	f.infoLog(fmt.Sprintf("Handling completion for form '%s'", current.ID))

	var currentValues *FormValues
	if stored, exists := f.formValues[current.ID]; exists {
		currentValues = &stored
		f.debugLog(fmt.Sprintf("Found existing values for form '%s'", current.ID))
	} else {
		currentValues = NewFormValues()
		f.debugLog(fmt.Sprintf("Creating new values for form '%s'", current.ID))
	}

	if f.config.Debug {
		f.debugLog(fmt.Sprintf("Form '%s' values before OnComplete:", current.ID))
		for key, valuePtr := range *currentValues {
			value := "(nil)"
			if valuePtr != nil {
				value = *valuePtr
			}
			f.debugLog(fmt.Sprintf("  %s: %s", key, value))
		}

		f.debugLog("Global values before OnComplete:")
		for key, valuePtr := range *f.globalData.Values {
			value := "(nil)"
			if valuePtr != nil {
				value = *valuePtr
			}
			f.debugLog(fmt.Sprintf("  %s: %s", key, value))
		}
	}

	currentData := FormData{
		ID:     current.ID,
		Values: currentValues,
	}

	if current.OnComplete != nil {
		f.debugLog(fmt.Sprintf("Calling OnComplete for form '%s'", current.ID))
		if err := current.OnComplete(&currentData, f.globalData); err != nil {
			f.errorLog(fmt.Errorf("OnComplete error for form '%s': %w", current.ID, err))
			f.addError(current.ID, err)
			return f, nil
		}
		f.debugLog(fmt.Sprintf("OnComplete for form '%s' completed successfully", current.ID))
	}

	if f.config.Debug {
		f.debugLog(fmt.Sprintf("Form '%s' values after OnComplete:", current.ID))
		for key, valuePtr := range *currentData.Values {
			value := "(nil)"
			if valuePtr != nil {
				value = *valuePtr
			}
			f.debugLog(fmt.Sprintf("  %s: %s", key, value))
		}

		f.debugLog("Global values after OnComplete:")
		for key, valuePtr := range *f.globalData.Values {
			value := "(nil)"
			if valuePtr != nil {
				value = *valuePtr
			}
			f.debugLog(fmt.Sprintf("  %s: %s", key, value))
		}
	}

	f.debugLog(fmt.Sprintf("Merging form '%s' values into global data", current.ID))
	f.globalData.Values.Merge(currentData.Values)

	if f.config.Debug {
		f.debugLog("Global values after merge:")
		for key, valuePtr := range *f.globalData.Values {
			value := "(nil)"
			if valuePtr != nil {
				value = *valuePtr
			}
			f.debugLog(fmt.Sprintf("  %s: %s", key, value))
		}
	}

	f.debugLog(fmt.Sprintf("Determining next form after '%s'", current.ID))
	nextIndex, err := f.navigator.Next(*f.globalData)
	if err != nil {
		f.errorLog(fmt.Errorf("navigation error from form '%s': %w", current.ID, err))
		f.addError(current.ID, err)
		return f, nil
	}

	f.infoLog(fmt.Sprintf("Next form index: %d", nextIndex))

	switch nextIndex {
	case -2:
		f.infoLog("Flow completed (nextIndex = -2)")
		f.state = StateCompleted
		return f, nil
	case -1:
		f.infoLog("No more forms (nextIndex = -1)")
		f.state = StateCompleted
		return f, nil
	default:
		f.infoLog(fmt.Sprintf("Moving to form at index %d", nextIndex))
		if err := f.navigator.MoveTo(nextIndex); err != nil {
			f.errorLog(fmt.Errorf("failed to move to form index %d: %w", nextIndex, err))
			f.addError(current.ID, err)
			return f, nil
		}
		return f, f.initCurrentForm()
	}
}

// initCurrentForm initializes the current form in the navigator.
// It handles skip conditions, generates the form, and prepares it for display.
func (f *Bobarista) initCurrentForm() tea.Cmd {
	current := f.navigator.Current()
	if current == nil {
		f.warningLog("initCurrentForm called with no current form")
		f.state = StateCompleted
		return nil
	}

	f.infoLog(fmt.Sprintf("Initializing form '%s'", current.ID))

	if f.formValues == nil {
		f.formValues = make(map[string]FormValues)
	}
	if _, exists := f.formValues[current.ID]; !exists {
		f.debugLog(fmt.Sprintf("Creating new form values for '%s'", current.ID))
		f.formValues[current.ID] = *NewFormValues()
	}

	currentValues := f.formValues[current.ID]

	currentData := FormData{
		ID:     current.ID,
		Values: &currentValues,
	}

	if f.config.Debug {
		f.debugLog(fmt.Sprintf("Global values before skip condition check for '%s':", current.ID))
		for key, valuePtr := range *f.globalData.Values {
			value := "(nil)"
			if valuePtr != nil {
				value = *valuePtr
			}
			f.debugLog(fmt.Sprintf("  %s: %s", key, value))
		}
	}

	if current.ShouldSkip != nil {
		f.debugLog(fmt.Sprintf("Checking skip condition for form '%s'", current.ID))
		shouldSkip := current.ShouldSkip(&currentData, f.globalData)
		f.infoLog(fmt.Sprintf("Form '%s' skip condition result: %t", current.ID, shouldSkip))

		if shouldSkip {
			f.infoLog(fmt.Sprintf("Skipping form '%s'", current.ID))

			nextIndex, err := f.navigator.Next(*f.globalData)
			if err != nil {
				f.errorLog(fmt.Errorf("navigation error while skipping form '%s': %w", current.ID, err))
				f.addError(current.ID, err)
				return nil
			}

			f.infoLog(fmt.Sprintf("Next form index after skipping '%s': %d", current.ID, nextIndex))

			if nextIndex == -1 || nextIndex == -2 {
				f.infoLog("No more forms after skip, completing flow")
				f.state = StateCompleted
				return nil
			}

			if err := f.navigator.MoveTo(nextIndex); err != nil {
				f.errorLog(fmt.Errorf("failed to move to form index %d after skipping '%s': %w", nextIndex, current.ID, err))
				f.addError(current.ID, err)
				return nil
			}

			return f.initCurrentForm()
		}
	} else {
		f.debugLog(fmt.Sprintf("No skip condition defined for form '%s'", current.ID))
	}

	if current.Generator == nil {
		f.errorLog(fmt.Errorf("no generator for form '%s'", current.ID))
		f.addError(current.ID, NewCupSleeveError(current.ID, ErrNoGenerator))
		return nil
	}

	f.debugLog(fmt.Sprintf("Generating form '%s'", current.ID))
	f.currentForm = current.Generator(&currentValues, f.globalData.Values)

	if f.currentForm == nil {
		f.errorLog(fmt.Errorf("generator returned nil form for '%s'", current.ID))
		f.addError(current.ID, NewCupSleeveError(current.ID, ErrNilForm))
		return nil
	}

	f.infoLog(fmt.Sprintf("Form '%s' initialized successfully", current.ID))
	return f.currentForm.Init()
}

// addError adds an error to the form flow and transitions to error state.
// It wraps the error in a CupSleeveError if it isn't already one.
func (f *Bobarista) addError(formID string, err error) {
	f.errorLog(fmt.Errorf("adding error for form '%s': %w", formID, err))
	if cupSleeveErr, ok := err.(CupSleeveError); ok {
		f.errors = append(f.errors, cupSleeveErr)
	} else {
		f.errors = append(f.errors, NewCupSleeveError(formID, err))
	}
	f.state = StateError
}

// GetGlobalData returns a copy of the global form data.
// This contains values that are shared across all forms in the flow.
func (f *Bobarista) GetGlobalData() FormData {
	if f.globalData == nil {
		return FormData{ID: "global", Values: NewFormValues()}
	}
	return *f.globalData
}

// GetCurrentFormData returns the form data for the currently active form.
// If no form is active, returns empty form data.
func (f *Bobarista) GetCurrentFormData() FormData {
	current := f.navigator.Current()
	if current == nil {
		return FormData{ID: "", Values: NewFormValues()}
	}

	var currentValues FormValues
	if stored, exists := f.formValues[current.ID]; exists {
		currentValues = stored
	} else {
		currentValues = *NewFormValues()
	}

	return FormData{
		ID:     current.ID,
		Values: &currentValues,
	}
}

// GetErrors returns all errors that have occurred during the form flow.
// This is useful for debugging and error handling.
func (f *Bobarista) GetErrors() []error {
	return f.errors
}
