# API Reference

## Core Types

### Bobarista
The main controller for form sequences.

```go
type Bobarista struct {
    // ... internal fields
}
```

#### Methods
- `Run() error` - Starts the form flow using Bubble Tea
- `GetGlobalData() FormData` - Returns the global form data
- `GetCurrentFormData() FormData` - Returns the current form data
- `GetErrors() []error` - Returns all errors that occurred during the flow

### Form
Represents a single form in the flow.

```go
type Form struct {
    ID         string
    Name       string
    Group      string
    Generator  FormGenerator
    OnComplete CompletionHandler
    ShouldSkip SkipCondition
    NextForm   NavigationHandler
    ShowStatus bool
}
```

### FormValues
A map of form field values.

```go
type FormValues map[string]*string
```

#### Methods
- `Set(key, value string)` - Sets a value
- `Get(key string) (string, bool)` - Gets a value
- `Has(key string) bool` - Checks if key exists
- `Delete(key string)` - Removes a key
- `Copy() FormValues` - Creates a copy
- `Merge(other *FormValues)` - Merges another FormValues

### FormData
Represents form data with ID and values.

```go
type FormData struct {
    ID     string
    Values *FormValues
}
```

## Builder API

### BobaBuilder
Provides a fluent API for building form flows.

#### Methods
- `AddForm(form Form) *BobaBuilder` - Adds a form to the flow
- `WithMaxWidth(width int) *BobaBuilder` - Sets maximum width
- `WithColorScheme(scheme string) *BobaBuilder` - Sets the color scheme
- `WithDisplayKeys(keys []string) *BobaBuilder` - Sets display keys for completion screen
- `OnInit(handler func(*Bobarista, []FormData)) *BobaBuilder` - Sets init callback
- `OnComplete(handler func(*Bobarista) error) *BobaBuilder` - Sets completion callback
- `WithDisplayCallback(callback func() string) *BobaBuilder` - Sets custom display callback
- `WithDebug(enabled bool) *BobaBuilder` - Enables/disables debug mode
- `Build() *Bobarista` - Creates the final Bobarista instance

## Function Types

### FormGenerator
```go
type FormGenerator func(current *FormValues, global *FormValues) *huh.Form
```

### CompletionHandler
```go
type CompletionHandler func(current *FormData, global *FormData) error
```

### SkipCondition
```go
type SkipCondition func(current *FormData, global *FormData) bool
```

### NavigationHandler
```go
type NavigationHandler func(current *FormData) int
```
Returns:
- `-1` for default navigation (next sequential form)
- `-2` for flow completion
- `>= 0` for specific form index

## Configuration

### Recipe
```go
type Recipe struct {
    Title           string
    MaxWidth        int
    DisplayKeys     []string
    ColorScheme     string
    Debug           bool
    OnInit          func(*Bobarista, []FormData)
    OnComplete      func(*Bobarista) error
    DisplayCallback func() string
}
```

## Error Handling

### CupSleeveError
```go
type CupSleeveError struct {
    FormID string
    Err    error
}
```

### Common Errors
- `ErrInvalidFormIndex` - Invalid form index provided
- `ErrFormNotFound` - Form not found
- `ErrNoFormsProvided` - No forms provided to flow
- `ErrNoValidForms` - No valid forms found
- `ErrNoGenerator` - Form generator is required
- `ErrNilForm` - Form generator returned nil
- `ErrEmptyFormID` - Form ID cannot be empty

### Error Types
- `DuplicateFormIDError` - Duplicate form IDs detected
- `ValidationError` - Form validation errors
- `NavigationError` - Navigation-related errors

## Color Schemes

### Available Schemes
- `default` - Default cyan and emerald theme
- `dark` - Dark theme with white primary
- `ubuntu` - Ubuntu orange theme
- `ocean` - Blue ocean theme
- `forest` - Green forest theme
- `sunset` - Orange sunset theme
- `monochrome` - Black and white theme

### Custom Color Schemes
```go
// Create custom color scheme
custom := bobarista.CreateCustomColorScheme("myscheme", "#FF0000", "#00FF00", "#0000FF")
bobarista.RegisterColorScheme("myscheme", custom)

// Use in builder
app := bobarista.New("My App").WithColorScheme("myscheme")
```

## Debug Mode

Enable debug mode to see internal state and navigation information:

```go
app := bobarista.New("Debug Example").
    WithDebug(true).
    AddForm(...).
    Build()
```

Debug mode provides:
- Current form state and values
- Global data inspection
- Navigation history
- Form completion status
- Error details

## Logging

Bobarista includes built-in logging functionality:

```go
// Global logging functions
bobarista.LogError(err)
bobarista.LogInfo("message")
bobarista.LogDebug("debug message")
bobarista.LogWarning("warning message")

// Set custom log filename
bobarista.LogFilename = "my-app.log"
```

Logs are stored in `~/.config/bobarista/` by default.
