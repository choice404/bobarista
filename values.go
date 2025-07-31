package bobarista

// FormValues represents a collection of form field values.
// It maps field names to their string values, with nil indicating unset fields.
type FormValues map[string]*string

// FormData represents the complete data for a form, including its ID and values.
// It's used to pass form information between different parts of the Bobarista flow.
type FormData struct {
	// ID is the unique identifier of the form this data belongs to.
	ID string

	// Values contains all the field values for this form.
	Values *FormValues
}

// NewFormValues creates a new, empty FormValues instance.
// All form values start uninitialized (nil pointers).
func NewFormValues() *FormValues {
	values := make(FormValues)
	return &values
}

// NewFormData creates a new FormData instance with the specified ID.
// The values are initialized as empty but ready to use.
func NewFormData(id string) FormData {
	return FormData{
		ID:     id,
		Values: NewFormValues(),
	}
}

// Set stores a value for the specified key.
// The value is copied to ensure the FormValues owns the data.
func (fv FormValues) Set(key, value string) {
	valueCopy := value
	fv[key] = &valueCopy
}

// Get retrieves the value for the specified key.
// Returns the value and true if the key exists and has a non-nil value,
// or empty string and false if the key doesn't exist or is nil.
func (fv FormValues) Get(key string) (string, bool) {
	if v, exists := fv[key]; exists && v != nil {
		return *v, true
	}
	return "", false
}

// Has checks whether a key exists in the FormValues.
// Returns true if the key exists (even if the value is nil).
func (fv FormValues) Has(key string) bool {
	_, exists := fv[key]
	return exists
}

// Delete removes a key and its value from the FormValues.
func (fv FormValues) Delete(key string) {
	delete(fv, key)
}

// Copy creates a deep copy of the FormValues.
// All string values are copied, ensuring the new FormValues is independent.
func (fv FormValues) Copy() FormValues {
	copy := make(FormValues)
	for k, v := range fv {
		if v != nil {
			value := *v
			copy[k] = &value
		} else {
			copy[k] = nil
		}
	}
	return copy
}

// Merge combines values from another FormValues into this one.
// Values from the other FormValues will overwrite existing values with the same key.
// If other is nil, no changes are made.
func (fv FormValues) Merge(other *FormValues) {
	if other == nil {
		return
	}
	for k, v := range *other {
		if v != nil {
			value := *v
			fv[k] = &value
		} else {
			fv[k] = nil
		}
	}
}
