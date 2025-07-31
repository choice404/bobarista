# Examples

This document provides various examples of using Bobarista for different scenarios.

## Basic Form

Simple single-form example:

```go
var name string

app := bobarista.New("Basic Example").
    AddForm(bobarista.NewForm("info", "Information").
        WithGenerator(func(current *bobarista.FormValues, global *bobarista.FormValues) *huh.Form {
            return huh.NewForm(
                huh.NewGroup(
                    huh.NewInput().Title("Name").Value(&name),
                ),
            )
        }).
        WithOnComplete(func(current *bobarista.FormData, global *bobarista.FormData) error {
            global.Values.Set("name", name)
            return nil
        })).
    Build()

app.Run()
```

## Conditional Forms

Forms that are skipped based on conditions:

```go
var userType, companyName, personalName string

app := bobarista.New("Conditional Example").
    AddForm(bobarista.NewForm("type", "User Type").
        WithGenerator(func(current *bobarista.FormValues, global *bobarista.FormValues) *huh.Form {
            return huh.NewForm(
                huh.NewGroup(
                    huh.NewSelect[string]().
                        Title("User type").
                        Options(
                            huh.NewOption("Individual", "individual"),
                            huh.NewOption("Company", "company"),
                        ).
                        Value(&userType),
                ),
            )
        }).
        WithOnComplete(func(current *bobarista.FormData, global *bobarista.FormData) error {
            global.Values.Set("user_type", userType)
            return nil
        })).
    AddForm(bobarista.NewForm("company", "Company Info").
        WithSkipCondition(func(current *bobarista.FormData, global *bobarista.FormData) bool {
            if val, exists := global.Values.Get("user_type"); exists {
                return val != "company"
            }
            return true
        }).
        WithGenerator(func(current *bobarista.FormValues, global *bobarista.FormValues) *huh.Form {
            return huh.NewForm(
                huh.NewGroup(
                    huh.NewInput().Title("Company Name").Value(&companyName),
                ),
            )
        }).
        WithOnComplete(func(current *bobarista.FormData, global *bobarista.FormData) error {
            global.Values.Set("company_name", companyName)
            return nil
        })).
    Build()
```

## Custom Navigation

Forms with custom navigation logic:

```go
app := bobarista.New("Navigation Example").
    AddForm(bobarista.NewForm("start", "Start").
        WithNavigation(func(current *bobarista.FormData) int {
            // Custom logic to determine next form
            if someCondition {
                return 2 // Jump to form index 2
            }
            return -1 // Go to next form (default navigation)
        }).
        WithGenerator(func(current *bobarista.FormValues, global *bobarista.FormValues) *huh.Form {
            // Form generator implementation
            return huh.NewForm(...)
        })).
    Build()
```

## Form Completion Callbacks

Handle form completion with callbacks:

```go
app := bobarista.New("Callback Example").
    AddForm(bobarista.NewForm("data", "Data Entry").
        WithOnComplete(func(current *bobarista.FormData, global *bobarista.FormData) error {
            // Process form data
            fmt.Printf("Form completed with data: %+v\n", current.Values)
            // Store values in global data
            for key, value := range *current.Values {
                if value != nil {
                    global.Values.Set(key, *value)
                }
            }
            return nil
        }).
        WithGenerator(func(current *bobarista.FormValues, global *bobarista.FormValues) *huh.Form {
            // Form generator implementation
            return huh.NewForm(...)
        })).
    Build()
```

## Multi-Step Wizard

Complex multi-step wizard with branching:

```go
var (
    projectType string
    language    string
    framework   string
)

app := bobarista.New("Project Wizard").
    AddForm(bobarista.NewForm("type", "Project Type").
        WithGenerator(func(current *bobarista.FormValues, global *bobarista.FormValues) *huh.Form {
            return huh.NewForm(
                huh.NewGroup(
                    huh.NewSelect[string]().
                        Title("Project type").
                        Options(
                            huh.NewOption("Web App", "web"),
                            huh.NewOption("CLI Tool", "cli"),
                        ).
                        Value(&projectType),
                ),
            )
        }).
        WithOnComplete(func(current *bobarista.FormData, global *bobarista.FormData) error {
            global.Values.Set("project_type", projectType)
            return nil
        })).
    AddForm(bobarista.NewForm("language", "Language").
        WithGenerator(func(current *bobarista.FormValues, global *bobarista.FormValues) *huh.Form {
            return huh.NewForm(
                huh.NewGroup(
                    huh.NewSelect[string]().
                        Title("Programming language").
                        Options(
                            huh.NewOption("Go", "go"),
                            huh.NewOption("Python", "python"),
                        ).
                        Value(&language),
                ),
            )
        }).
        WithOnComplete(func(current *bobarista.FormData, global *bobarista.FormData) error {
            global.Values.Set("language", language)
            return nil
        })).
    AddForm(bobarista.NewForm("framework", "Framework").
        WithSkipCondition(func(current *bobarista.FormData, global *bobarista.FormData) bool {
            if val, exists := global.Values.Get("project_type"); exists {
                return val == "cli" // Skip for CLI tools
            }
            return false
        }).
        WithGenerator(func(current *bobarista.FormValues, global *bobarista.FormValues) *huh.Form {
            var options []huh.Option[string]
            if langVal, exists := global.Get("language"); exists && langVal == "go" {
                options = []huh.Option[string]{
                    huh.NewOption("Gin", "gin"),
                    huh.NewOption("Echo", "echo"),
                }
            } else {
                options = []huh.Option[string]{
                    huh.NewOption("Django", "django"),
                    huh.NewOption("Flask", "flask"),
                }
            }
            return huh.NewForm(
                huh.NewGroup(
                    huh.NewSelect[string]().
                        Title("Framework").
                        Options(options...).
                        Value(&framework),
                ),
            )
        }).
        WithOnComplete(func(current *bobarista.FormData, global *bobarista.FormData) error {
            global.Values.Set("framework", framework)
            return nil
        })).
    OnComplete(func(app *bobarista.Bobarista) error {
        fmt.Printf("Project setup: %s, %s, %s\n", projectType, language, framework)
        return nil
    }).
    Build()
```

## Debug Mode

Enable debug mode to see internal state:

```go
app := bobarista.New("Debug Example").
    WithDebug(true).  // Enable debug panel
    AddForm(bobarista.NewForm("info", "Information").
        WithGenerator(func(current *bobarista.FormValues, global *bobarista.FormValues) *huh.Form {
            return huh.NewForm(...)
        })).
    Build()
```

## Color Schemes

Using different color schemes:

```go
app := bobarista.New("Styled Example").
    WithColorScheme("ocean").  // Use ocean theme
    AddForm(bobarista.NewForm("info", "Information").
        WithGenerator(func(current *bobarista.FormValues, global *bobarista.FormValues) *huh.Form {
            return huh.NewForm(...)
        })).
    Build()
```

## Custom Display

Custom completion display:

```go
app := bobarista.New("Custom Display").
    AddForm(bobarista.NewForm("info", "Information").
        WithGenerator(func(current *bobarista.FormValues, global *bobarista.FormValues) *huh.Form {
            return huh.NewForm(...)
        })).
    WithDisplayCallback(func() string {
        return "🎉 Custom completion message with styling!\n\nThank you for using Bobarista!"
    }).
    Build()
```

## Form Groups

Organize forms into logical groups:

```go
app := bobarista.New("Grouped Forms").
    AddForm(bobarista.NewForm("personal", "Personal Info").
        WithGenerator(func(current *bobarista.FormValues, global *bobarista.FormValues) *huh.Form {
            form := bobarista.NewForm("personal", "Personal Information")
            form.Group = "user_info"
            return form
        })).
    AddForm(bobarista.NewForm("contact", "Contact Info").
        WithGenerator(func(current *bobarista.FormValues, global *bobarista.FormValues) *huh.Form {
            form := bobarista.NewForm("contact", "Contact Information")
            form.Group = "user_info"
            return form
        })).
    Build()
```

## Error Handling

Handle errors gracefully:

```go
app := bobarista.New("Error Handling Example").
    AddForm(bobarista.NewForm("validation", "Validation").
        WithOnComplete(func(current *bobarista.FormData, global *bobarista.FormData) error {
            // Validate data
            if name, exists := current.Values.Get("name"); !exists || name == "" {
                return bobarista.NewValidationError("validation", "name", 
                    errors.New("name is required"))
            }
            return nil
        }).
        WithGenerator(func(current *bobarista.FormValues, global *bobarista.FormValues) *huh.Form {
            return huh.NewForm(...)
        })).
    Build()

if err := app.Run(); err != nil {
    if cupSleeveErr, ok := err.(bobarista.CupSleeveError); ok {
        fmt.Printf("Form error in %s: %v\n", cupSleeveErr.FormID, cupSleeveErr.Err)
    } else {
        fmt.Printf("General error: %v\n", err)
    }
}
```

## Display Keys

Control which values are shown in the completion screen:

```go
app := bobarista.New("Display Keys Example").
    WithDisplayKeys([]string{"name", "email", "project_type"}).  // Only show these keys
    AddForm(bobarista.NewForm("info", "Information").
        WithGenerator(func(current *bobarista.FormValues, global *bobarista.FormValues) *huh.Form {
            return huh.NewForm(...)
        }).
        WithOnComplete(func(current *bobarista.FormData, global *bobarista.FormData) error {
            global.Values.Set("name", name)
            global.Values.Set("email", email)
            global.Values.Set("internal_id", "12345")  // Won't be displayed
            return nil
        })).
    Build()
```

## Initialization Callback

Handle initialization with custom logic:

```go
app := bobarista.New("Init Example").
    OnInit(func(app *bobarista.Bobarista, formDataList []bobarista.FormData) {
        // Initialize global data
        globalData := app.GetGlobalData()
        globalData.Values.Set("app_version", "1.0.0")
        globalData.Values.Set("start_time", time.Now().Format(time.RFC3339))
    }).
    AddForm(bobarista.NewForm("info", "Information").
        WithGenerator(func(current *bobarista.FormValues, global *bobarista.FormValues) *huh.Form {
            return huh.NewForm(...)
        })).
    Build()
```

For more examples, check the `examples/` directory in the repository.

The README and documentation have been updated to be more engaging, comprehensive, and user-friendly. The README now includes:

1. **Fun branding** with the boba tea theme and emojis
2. **Clear feature highlights** with visual icons
3. **Comprehensive examples** showing real-world usage
4. **Beautiful code samples** with proper syntax highlighting
5. **Better organization** with clear sections
6. **Updated variable names** (changed from `cupSleeve` to `app` for clarity)
7. **More engaging language** while maintaining technical accuracy
8. **Better visual hierarchy** with proper markdown formatting

The documentation files have also been updated to match the new style and ensure consistency across all materials.
