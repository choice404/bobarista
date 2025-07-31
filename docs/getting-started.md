# Getting Started with Bobarista

Bobarista is a powerful form flow library that makes it easy to create sequential, interactive terminal forms using Charm's Bubble Tea and Huh libraries.

## Installation

```bash
go get github.com/choice404/bobarista
```

## Quick Start

Here's a simple example to get you started:

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
    
    app := bobarista.New("User Registration").
        AddForm(bobarista.NewForm("info", "User Information").
            WithGenerator(func(current *bobarista.FormValues, global *bobarista.FormValues) *huh.Form {
                return huh.NewForm(
                    huh.NewGroup(
                        huh.NewInput().Title("Name").Value(&name),
                        huh.NewInput().Title("Email").Value(&email),
                    ),
                )
            }).
            WithOnComplete(func(current *bobarista.FormData, global *bobarista.FormData) error {
                global.Values.Set("name", name)
                global.Values.Set("email", email)
                return nil
            })).
        OnComplete(func(app *bobarista.Bobarista) error {
            fmt.Printf("Hello %s! Your email is %s\n", name, email)
            return nil
        }).
        Build()
    
    if err := app.Run(); err != nil {
        log.Fatal(err)
    }
}
```

## Core Concepts

### Forms
Forms are the building blocks of your flow. Each form represents a single step in your user interaction.

### Form Generator
The generator function creates the actual huh.Form that will be displayed to the user. It receives both current form values and global values as parameters.

### Bobarista Builder (BobaBuilder)
The fluent API that lets you chain together forms and configure the overall flow behavior.

### Form Values
The system automatically manages form values through the FormValues type, allowing data to persist across forms. Values are stored as pointers to strings for efficient memory usage.

## Key Features

- **Sequential Forms**: Chain multiple forms together
- **Conditional Logic**: Skip forms based on previous answers
- **Custom Navigation**: Control form flow with custom logic
- **Value Persistence**: Automatically maintain form data across steps
- **Theming**: Multiple built-in color schemes
- **Error Handling**: Comprehensive error management with specific error types
- **Completion Callbacks**: Handle form completion with custom logic
- **Debug Mode**: Built-in debugging with visual state inspection
- **Logging**: Integrated logging system for troubleshooting

## Basic Pattern

The typical pattern for using Bobarista:

1. **Create variables** to hold form input
2. **Define forms** with generators that create huh.Form instances
3. **Add completion handlers** to store values in global data
4. **Chain forms together** using the builder pattern
5. **Run the flow** and handle completion

```go
var input string

app := bobarista.New("My Flow").
    AddForm(bobarista.NewForm("step1", "Step 1").
        WithGenerator(func(current *bobarista.FormValues, global *bobarista.FormValues) *huh.Form {
            // Create your huh.Form here
            return huh.NewForm(
                huh.NewGroup(
                    huh.NewInput().Title("Input").Value(&input),
                ),
            )
        }).
        WithOnComplete(func(current *bobarista.FormData, global *bobarista.FormData) error {
            // Store values in global data
            global.Values.Set("key", input)
            return nil
        })).
    Build()

app.Run()
```

## Advanced Features

### Debug Mode
Enable debug mode to see internal state, form values, and navigation information:

```go
app := bobarista.New("Debug Example").
    WithDebug(true).  // Shows debug panel
    AddForm(...).
    Build()
```

### Custom Color Schemes
Choose from built-in themes or create your own:

```go
app := bobarista.New("Styled App").
    WithColorScheme("ocean").  // Built-in theme
    AddForm(...).
    Build()

// Or create a custom theme
customScheme := bobarista.CreateCustomColorScheme("custom", "#FF0000", "#00FF00", "#0000FF")
bobarista.RegisterColorScheme("custom", customScheme)
```

### Conditional Forms
Skip forms based on previous answers:

```go
app := bobarista.New("Conditional Flow").
    AddForm(bobarista.NewForm("condition", "Condition").
        WithGenerator(...).
        WithOnComplete(func(current *bobarista.FormData, global *bobarista.FormData) error {
            global.Values.Set("user_type", userType)
            return nil
        })).
    AddForm(bobarista.NewForm("optional", "Optional Form").
        WithSkipCondition(func(current *bobarista.FormData, global *bobarista.FormData) bool {
            if val, exists := global.Values.Get("user_type"); exists {
                return val != "advanced"  // Skip unless user is advanced
            }
            return true
        }).
        WithGenerator(...)).
    Build()
```

### Error Handling
Bobarista provides comprehensive error handling:

```go
if err := app.Run(); err != nil {
    if cupSleeveErr, ok := err.(bobarista.CupSleeveError); ok {
        fmt.Printf("Form error in %s: %v\n", cupSleeveErr.FormID, cupSleeveErr.Err)
    } else {
        fmt.Printf("General error: %v\n", err)
    }
}
```

### Logging
Built-in logging helps with debugging:

```go
// Enable debug mode to see logs
app := bobarista.New("Logged App").
    WithDebug(true).
    AddForm(...).
    Build()

// Or use logging functions directly
bobarista.LogInfo("Application started")
bobarista.LogError(fmt.Errorf("something went wrong"))
```

## Form Value Management

### Setting Values
```go
.WithOnComplete(func(current *bobarista.FormData, global *bobarista.FormData) error {
    // Store individual values
    global.Values.Set("name", name)
    global.Values.Set("email", email)
    
    // Or merge all current form values
    global.Values.Merge(current.Values)
    return nil
})
```

### Reading Values
```go
.WithGenerator(func(current *bobarista.FormValues, global *bobarista.FormValues) *huh.Form {
    // Check if a value exists
    if name, exists := global.Get("name"); exists {
        // Use the existing value
        defaultName = name
    }
    
    return huh.NewForm(...)
})
```

### Display Control
Control what's shown in the completion screen:

```go
app := bobarista.New("Controlled Display").
    WithDisplayKeys([]string{"name", "email"}).  // Only show these keys
    // Or use a custom display callback
    WithDisplayCallback(func() string {
        return "Custom completion message!"
    }).
    AddForm(...).
    Build()
```

## Best Practices

1. **Always use OnComplete handlers** to store form values in global data
2. **Use meaningful form IDs** for easier debugging and navigation
3. **Enable debug mode during development** to understand form flow
4. **Handle errors gracefully** with proper error checking
5. **Use skip conditions** to create dynamic form flows
6. **Organize related forms with groups** for better structure
7. **Test navigation logic thoroughly** especially with custom navigation handlers

## Common Patterns

### Multi-Step Registration
```go
app := bobarista.New("Registration").
    AddForm(bobarista.NewForm("personal", "Personal Info").WithGenerator(...)).
    AddForm(bobarista.NewForm("contact", "Contact Info").WithGenerator(...)).
    AddForm(bobarista.NewForm("preferences", "Preferences").WithGenerator(...)).
    OnComplete(func(app *bobarista.Bobarista) error {
        // Process complete registration
        return saveUser(app.GetGlobalData())
    }).
    Build()
```

### Configuration Wizard
```go
app := bobarista.New("Setup Wizard").
    AddForm(bobarista.NewForm("type", "Installation Type").WithGenerator(...)).
    AddForm(bobarista.NewForm("database", "Database Config").
        WithSkipCondition(func(current *bobarista.FormData, global *bobarista.FormData) bool {
            // Skip if using embedded database
            if val, exists := global.Values.Get("install_type"); exists {
                return val == "embedded"
            }
            return false
        }).
        WithGenerator(...)).
    Build()
```

## Next Steps

- Check out the [API Reference](api-reference.md) for detailed documentation
- Browse the [Examples](examples.md) for more complex use cases
- See the `examples/` directory for runnable code samples
- Learn about [color schemes and theming](api-reference.md#color-schemes)
