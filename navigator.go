package bobarista

// Navigator manages form navigation and flow control within a Bobarista form flow.
// It handles moving between forms, tracking history, and determining valid navigation paths.
type Navigator struct {
	forms      []Form
	currentIdx int
	history    []int
}

// NewNavigator creates a new Navigator with the provided forms.
// The navigator starts with no current form (currentIdx = -1).
func NewNavigator(forms []Form) *Navigator {
	return &Navigator{
		forms:      forms,
		currentIdx: -1,
		history:    make([]int, 0),
	}
}

// Current returns the currently active form, or nil if no form is active.
func (n *Navigator) Current() *Form {
	if n.currentIdx >= 0 && n.currentIdx < len(n.forms) {
		return &n.forms[n.currentIdx]
	}
	return nil
}

// Next determines the index of the next form to navigate to.
// It considers custom navigation handlers and skip conditions.
// Returns -1 if no more forms are available, -2 if the flow should complete.
func (n *Navigator) Next(data FormData) (int, error) {
	current := n.Current()
	if current == nil {
		return n.findNextValidForm(0, data)
	}

	if current.NextForm != nil {
		nextIdx := current.NextForm(&data)
		if nextIdx == -1 {
			// Use default navigation (next sequential form)
			return n.findNextValidForm(n.currentIdx+1, data)
		}
		return nextIdx, nil
	}

	return n.findNextValidForm(n.currentIdx+1, data)
}

// MoveTo navigates to the form at the specified index.
// It adds the current form to the navigation history before moving.
func (n *Navigator) MoveTo(index int) error {
	if index < 0 || index >= len(n.forms) {
		return ErrInvalidFormIndex
	}

	if n.currentIdx >= 0 {
		n.history = append(n.history, n.currentIdx)
	}

	n.currentIdx = index
	return nil
}

// MoveToFirstValid finds and navigates to the first valid (non-skipped) form.
// This is typically called during initialization.
func (n *Navigator) MoveToFirstValid(globalData FormData) error {
	nextIdx, err := n.findNextValidForm(0, globalData)
	if err != nil {
		return err
	}

	if nextIdx == -1 {
		return ErrNoValidForms
	}

	n.currentIdx = nextIdx
	return nil
}

// Back navigates to the previous form in the history.
// Returns true if navigation was successful, false if no history exists.
func (n *Navigator) Back() bool {
	if len(n.history) == 0 {
		return false
	}

	n.currentIdx = n.history[len(n.history)-1]
	n.history = n.history[:len(n.history)-1]
	return true
}

// HasNext returns true if there are more forms after the current one.
func (n *Navigator) HasNext() bool {
	return n.currentIdx < len(n.forms)-1
}

// HasPrevious returns true if there are forms in the navigation history.
func (n *Navigator) HasPrevious() bool {
	return len(n.history) > 0
}

// GetFormCount returns the total number of forms in the navigator.
func (n *Navigator) GetFormCount() int {
	return len(n.forms)
}

// GetCurrentIndex returns the index of the currently active form.
// Returns -1 if no form is active.
func (n *Navigator) GetCurrentIndex() int {
	return n.currentIdx
}

// GetProgress calculates the current progress as a percentage.
// Returns 0.0 if no forms exist or no form is active.
func (n *Navigator) GetProgress() float64 {
	if len(n.forms) == 0 {
		return 0.0
	}
	if n.currentIdx < 0 {
		return 0.0
	}
	return float64(n.currentIdx+1) / float64(len(n.forms)) * 100.0
}

// findNextValidForm searches for the next form that should not be skipped.
// It starts from startIdx and checks each form's skip condition.
// Returns -1 if no valid forms are found.
func (n *Navigator) findNextValidForm(startIdx int, globalData FormData) (int, error) {
	for i := startIdx; i < len(n.forms); i++ {
		form := &n.forms[i]

		tempData := FormData{
			ID:     form.ID,
			Values: NewFormValues(),
		}

		if form.ShouldSkip == nil || !form.ShouldSkip(&tempData, &globalData) {
			return i, nil
		}
	}
	return -1, nil
}

// Reset resets the navigator to its initial state.
// This clears the current form and navigation history.
func (n *Navigator) Reset() {
	n.currentIdx = -1
	n.history = make([]int, 0)
}

// GetFormByID finds a form by its ID and returns the form, its index, and any error.
// Returns ErrFormNotFound if the form doesn't exist.
func (n *Navigator) GetFormByID(id string) (*Form, int, error) {
	for i, form := range n.forms {
		if form.ID == id {
			return &form, i, nil
		}
	}
	return nil, -1, ErrFormNotFound
}

// GetFormsByGroup returns all forms that belong to the specified group.
// Returns an empty slice if no forms match the group.
func (n *Navigator) GetFormsByGroup(group string) []Form {
	var forms []Form
	for _, form := range n.forms {
		if form.Group == group {
			forms = append(forms, form)
		}
	}
	return forms
}

// ValidateNavigation validates the form configuration and returns any errors found.
// It checks for missing forms, empty IDs, duplicate IDs, and missing generators.
func (n *Navigator) ValidateNavigation() []error {
	var errors []error

	if len(n.forms) == 0 {
		errors = append(errors, ErrNoFormsProvided)
		return errors
	}

	idMap := make(map[string]int)
	for i, form := range n.forms {
		if form.ID == "" {
			errors = append(errors, NewCupSleeveError("", ErrEmptyFormID))
			continue
		}

		if existingIdx, exists := idMap[form.ID]; exists {
			errors = append(errors, NewCupSleeveError(form.ID,
				NewDuplicateFormIDError(form.ID, existingIdx, i)))
		} else {
			idMap[form.ID] = i
		}
	}

	for _, form := range n.forms {
		if form.Generator == nil {
			errors = append(errors, NewCupSleeveError(form.ID, ErrNoGenerator))
		}
	}

	return errors
}
