# 🧋 Bobarista

*A delightfully smooth form flow framework for Go terminal applications*

[![Go Version](https://img.shields.io/badge/go-%3E%3D1.21-blue.svg)](https://golang.org/)
[![License](https://img.shields.io/badge/license-MIT-green.svg)](LICENSE)
[![Go Report Card](https://goreportcard.com/badge/github.com/choice404/bobarista)](https://goreportcard.com/report/github.com/choice404/bobarista)

Bobarista is a powerful, flexible form flow framework built on top of [Charm's Bubble Tea](https://github.com/charmbracelet/bubbletea) and [Huh](https://github.com/charmbracelet/huh). Create beautiful, interactive multi-step terminal forms with conditional logic, custom navigation, and delightful theming - all with a fluent, easy-to-use API.

## Features

- **Sequential Forms** - Chain multiple forms together seamlessly
-  **Conditional Logic** - Skip forms based on user input
-  **Custom Navigation** - Control form flow with custom logic
-  **Value Persistence** - Automatic data management across forms
-  **Beautiful Themes** - 7 built-in color schemes + custom themes
-  **Debug Mode** - Visual debugging with state inspection
-  **Comprehensive Logging** - Built-in logging for troubleshooting
-  **Error Handling** - Robust error management with specific error types
-  **Completion Callbacks** - Handle form completion with custom logic
-  **Responsive** - Adapts to terminal size with scrolling support

##  Quick Start

```bash
go get github.com/choice404/bobarista
```

### Your First Bobarista Form

```go
package main

import (
    "fmt"
    "log"
    
    "github.com/charmbracelet/huh"
    "github.com/choice404/bobarista/pkg/bobarista"
)

func main() {
    var name, email string
    
    // Create your form flow
    app := bobarista.New("🧋 Welcome to Bobarista").
        WithColorScheme("ocean").
        AddForm(bobarista.NewForm("welcome", "User Information").
            WithGenerator(func(current *bobarista.FormValues, global *bobarista.FormValues) *huh.Form {
                return huh.NewForm(
                    huh.NewGroup(
                        huh.NewInput().
                            Title("What's your name?").
                            Placeholder("Enter your name").
                            Value(&name),
                        huh.NewInput().
                            Title("What's your email?").
                            Placeholder("you@example.com").
                            Value(&email),
                    ),
                )
            }).
            WithOnComplete(func(current *bobarista.FormData, global *bobarista.FormData) error {
                global.Values.Set("name", name)
                global.Values.Set("email", email)
                return nil
            })).
        OnComplete(func(app *bobarista.Bobarista) error {
            fmt.Printf("Welcome %s! We'll contact you at %s\n", name, email)
            return nil
        }).
        Build()
    
    if err := app.Run(); err != nil {
        log.Fatal(err)
    }
}
```

##  Beautiful Themes

Bobarista comes with 7 stunning built-in themes:

```go
app := bobarista.New("Themed App").
    WithColorScheme("ocean").    // Ocean blue
    // WithColorScheme("forest").   // Forest green  
    // WithColorScheme("sunset").   // Sunset orange
    // WithColorScheme("ubuntu").   // Ubuntu orange
    // WithColorScheme("dark").     // Dark mode
    // WithColorScheme("monochrome"). // Black & white
    // WithColorScheme("default").  // Default cyan
    Build()
```

Or create your own custom theme:

```go
customTheme := bobarista.CreateCustomColorScheme("neon", "#FF00FF", "#00FFFF", "#FFFF00")
bobarista.RegisterColorScheme("neon", customTheme)

app := bobarista.New("Neon App").WithColorScheme("neon")
```

## Conditional Forms & Smart Navigation

Create dynamic form flows that adapt based on user input:

```go
var userType, companyName, personalDetails string

app := bobarista.New("Smart Registration").
    AddForm(bobarista.NewForm("type", "User Type").
        WithGenerator(func(current *bobarista.FormValues, global *bobarista.FormValues) *huh.Form {
            return huh.NewForm(
                huh.NewGroup(
                    huh.NewSelect[string]().
                        Title("Are you registering as...").
                        Options(
                            huh.NewOption("Company", "company"),
                            huh.NewOption("Individual", "individual"),
                        ).
                        Value(&userType),
                ),
            )
        }).
        WithOnComplete(func(current *bobarista.FormData, global *bobarista.FormData) error {
            global.Values.Set("user_type", userType)
            return nil
        })).
    AddForm(bobarista.NewForm("company", "Company Details").
        WithSkipCondition(func(current *bobarista.FormData, global *bobarista.FormData) bool {
            // Skip this form if user is not a company
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
        })).
    Build()
```

## Debug Mode

Enable debug mode to see what's happening under the hood:

```go
app := bobarista.New("🔍 Debug Example").
    WithDebug(true).  // Shows a debug panel with form state, values, and navigation
    AddForm(...).
    Build()
```

The debug panel shows:
- Current form state and values
- Global data inspection  
- Navigation history and progress
- Form completion status
- Real-time error information

## Advanced Examples

### Multi-Step Project Setup Wizard

```go
var projectType, language, framework, features string

wizard := bobarista.New("Project Setup Wizard").
    WithColorScheme("forest").
    WithMaxWidth(100).
    AddForm(bobarista.NewForm("project", "Project Type").
        WithGenerator(func(current *bobarista.FormValues, global *bobarista.FormValues) *huh.Form {
            return huh.NewForm(
                huh.NewGroup(
                    huh.NewSelect[string]().
                        Title("What type of project?").
                        Options(
                            huh.NewOption("Web Application", "web"),
                            huh.NewOption("⚡ CLI Tool", "cli"),
                            huh.NewOption("📱 Mobile App", "mobile"),
                        ).
                        Value(&projectType),
                ),
            )
        })).
    AddForm(bobarista.NewForm("language", "Programming Language").
        WithGenerator(func(current *bobarista.FormValues, global *bobarista.FormValues) *huh.Form {
            var options []huh.Option[string]
            
            // Dynamic options based on project type
            if projectType == "web" {
                options = []huh.Option[string]{
                    huh.NewOption("Go", "go"),
                    huh.NewOption("Python", "python"),
                    huh.NewOption("Node.js", "nodejs"),
                }
            } else {
                options = []huh.Option[string]{
                    huh.NewOption("Go", "go"),
                    huh.NewOption("Rust", "rust"),
                    huh.NewOption("C++", "cpp"),
                }
            }
            
            return huh.NewForm(
                huh.NewGroup(
                    huh.NewSelect[string]().
                        Title("Choose your language").
                        Options(options...).
                        Value(&language),
                ),
            )
        })).
    OnComplete(func(wizard *bobarista.Bobarista) error {
        fmt.Printf("Project created: %s %s application!\n", language, projectType)
        return nil
    }).
    Build()
```

### Custom Navigation Logic

```go
app := bobarista.New("Custom Navigation").
    AddForm(bobarista.NewForm("start", "Getting Started").
        WithNavigation(func(current *bobarista.FormData) int {
            // Custom navigation logic
            if someCondition {
                return 3 // Jump to form index 3
            }
            if anotherCondition {
                return -2 // Complete the flow early
            }
            return -1 // Use default navigation (next form)
        })).
    Build()
```

## Form Value Management

Bobarista automatically manages form values with a powerful and flexible system:

```go
// Setting values in completion handlers
.WithOnComplete(func(current *bobarista.FormData, global *bobarista.FormData) error {
    // Set individual values
    global.Values.Set("username", username)
    global.Values.Set("email", email)
    
    // Or merge all current form values
    global.Values.Merge(current.Values)
    
    return nil
})

// Reading values in generators
.WithGenerator(func(current *bobarista.FormValues, global *bobarista.FormValues) *huh.Form {
    // Check if a value exists and use it as default
    var defaultName string
    if name, exists := global.Get("name"); exists {
        defaultName = name
    }
    
    return huh.NewForm(
        huh.NewGroup(
            huh.NewInput().
                Title("Name").
                Value(&defaultName),
        ),
    )
})
```

## Completion & Display Control

Customize what happens when your form flow completes:

```go
app := bobarista.New("Survey Complete").
    WithDisplayKeys([]string{"name", "email", "rating"}). // Only show these in summary
    AddForm(...).
    OnComplete(func(app *bobarista.Bobarista) error {
        // Save to database, send email, etc.
        return saveResults(app.GetGlobalData())
    }).
    WithDisplayCallback(func() string {
        return `
Thank you for your feedback!

Your responses have been saved and our team will review them shortly.
You should receive a confirmation email within the next few minutes.

Have a great day!
        `
    }).
    Build()
```

## Error Handling

Bobarista provides comprehensive error handling with specific error types:

```go
if err := app.Run(); err != nil {
    switch e := err.(type) {
    case bobarista.CupSleeveError:
        fmt.Printf("Form error in '%s': %v\n", e.FormID, e.Err)
    case bobarista.ValidationError:
        fmt.Printf("Validation error in form '%s', field '%s': %v\n", 
            e.FormID, e.Field, e.Err)
    case bobarista.NavigationError:
        fmt.Printf("Navigation error from '%s' to '%s': %v\n", 
            e.FromFormID, e.ToFormID, e.Err)
    default:
        fmt.Printf("Unexpected error: %v\n", err)
    }
}
```

## 📝 Built-in Logging

Bobarista includes comprehensive logging for debugging and monitoring:

```go
// Logs are automatically written to ~/.config/bobarista/bobarista.log
// Or customize the location:
bobarista.LogFilename = "my-app.log"

// Use logging functions directly
bobarista.LogInfo("Application started")
bobarista.LogError(fmt.Errorf("something went wrong"))
bobarista.LogDebug("Debug information")
bobarista.LogWarning("Warning message")
```

## Builder API Reference

The fluent builder API makes it easy to construct complex form flows:

```go
app := bobarista.New("My App").                    // Create new builder
    WithMaxWidth(120).                             // Set max display width
    WithColorScheme("ocean").                      // Choose color theme
    WithDisplayKeys([]string{"name", "email"}).    // Control completion display
    WithDebug(true).                               // Enable debug mode
    OnInit(func(app *bobarista.Bobarista, forms []bobarista.FormData) {
        // Initialize global data
    }).
    OnComplete(func(app *bobarista.Bobarista) error {
        // Handle completion
        return nil
    }).
    WithDisplayCallback(func() string {
        return "Custom completion message"
    }).
    AddForm(bobarista.NewForm("form1", "First Form").
        WithGenerator(func(current *bobarista.FormValues, global *bobarista.FormValues) *huh.Form {
            // Create huh.Form
        }).
        WithOnComplete(func(current *bobarista.FormData, global *bobarista.FormData) error {
            // Handle form completion
            return nil
        }).
        WithSkipCondition(func(current *bobarista.FormData, global *bobarista.FormData) bool {
            // Return true to skip this form
            return false
        }).
        WithNavigation(func(current *bobarista.FormData) int {
            // Custom navigation logic
            return -1 // Use default navigation
        }).
        WithoutStatus()). // Don't show progress for this form
    Build()
```

## Common Use Cases

- **User Registration & Onboarding** - Multi-step user signup flows
- **Configuration Wizards** - Application setup and configuration
- **Surveys & Forms** - Data collection with conditional questions  
- **Project Generators** - CLI tools for scaffolding projects
- **Installation Wizards** - Software installation and setup
- **Data Entry Applications** - Complex data input workflows
- **Interactive CLI Tools** - User-friendly command-line interfaces

## Documentation

- [Getting Started Guide](getting-started.md) - Detailed tutorial and concepts
- [API Reference](api-reference.md) - Complete API documentation  
- [Examples](examples.md) - More complex examples and patterns
- [Theming Guide](api-reference.md#color-schemes) - Color schemes and customization

## Contributing

We love contributions! Please see our [Contributing Guide](CONTRIBUTING.md) for details.

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add some amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## Acknowledgments

- [Charm](https://charm.sh/) for the amazing Bubble Tea and Huh libraries
- The Go community for inspiration and feedback
- All contributors who help make Bobarista better

## Star History

If you find Bobarista useful, please consider giving it a star! ⭐
