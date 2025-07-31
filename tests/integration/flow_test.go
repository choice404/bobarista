package integration

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/charmbracelet/huh"
	"github.com/choice404/bobarista/pkg/bobarista"
	"github.com/stretchr/testify/assert"
)

func TestBasicFlow(t *testing.T) {
	var name string
	var completed bool

	boba := bobarista.New("Test Cupsleeve").
		AddForm(bobarista.NewForm("test", "Test Form").
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
		OnComplete(func(boba *bobarista.Bobarista) error {
			completed = true
			return nil
		}).
		Build()

	assert.NotNil(t, boba)
	assert.False(t, completed)

	globalData := boba.GetGlobalData()
	assert.Equal(t, "global", globalData.ID)
	assert.NotNil(t, globalData.Values)
}

func TestFormValues(t *testing.T) {
	values := bobarista.NewFormValues()

	values.Set("key1", "value1")
	val, exists := values.Get("key1")
	assert.True(t, exists)
	assert.Equal(t, "value1", val)

	assert.True(t, values.Has("key1"))
	assert.False(t, values.Has("nonexistent"))

	values.Delete("key1")
	assert.False(t, values.Has("key1"))

	values.Set("key2", "value2")
	copied := values.Copy()
	val, exists = copied.Get("key2")
	assert.True(t, exists)
	assert.Equal(t, "value2", val)

	other := bobarista.NewFormValues()
	other.Set("key3", "value3")
	values.Merge(other)
	val, exists = values.Get("key3")
	assert.True(t, exists)
	assert.Equal(t, "value3", val)
}

func TestFormData(t *testing.T) {
	data := bobarista.NewFormData("test-form")
	assert.Equal(t, "test-form", data.ID)
	assert.NotNil(t, data.Values)

	data.Values.Set("field1", "value1")
	val, exists := data.Values.Get("field1")
	assert.True(t, exists)
	assert.Equal(t, "value1", val)
}

func TestBobaBuilder(t *testing.T) {
	builder := bobarista.New("Test")
	assert.NotNil(t, builder)

	form := bobarista.NewForm("test", "Test Form").
		WithGenerator(func(current *bobarista.FormValues, global *bobarista.FormValues) *huh.Form {
			var testValue string
			return huh.NewForm(
				huh.NewGroup(
					huh.NewInput().Title("Test").Value(&testValue),
				),
			)
		})

	builder = builder.AddForm(form)
	builder = builder.WithMaxWidth(120)
	builder = builder.WithColorScheme("dark")
	builder = builder.WithDisplayKeys([]string{"test_key"})

	boba := builder.Build()
	assert.NotNil(t, boba)
}

func TestNavigator(t *testing.T) {
	forms := []bobarista.Form{
		bobarista.NewForm("form1", "Form 1").WithGenerator(func(current *bobarista.FormValues, global *bobarista.FormValues) *huh.Form {
			var val string
			return huh.NewForm(huh.NewGroup(huh.NewInput().Value(&val)))
		}),
		bobarista.NewForm("form2", "Form 2").WithGenerator(func(current *bobarista.FormValues, global *bobarista.FormValues) *huh.Form {
			var val string
			return huh.NewForm(huh.NewGroup(huh.NewInput().Value(&val)))
		}),
	}

	navigator := bobarista.NewNavigator(forms)
	assert.NotNil(t, navigator)
	assert.Equal(t, 2, navigator.GetFormCount())
	assert.Equal(t, -1, navigator.GetCurrentIndex())

	assert.True(t, navigator.HasNext())
	assert.False(t, navigator.HasPrevious())

	globalData := bobarista.NewFormData("global")
	err := navigator.MoveToFirstValid(globalData)
	assert.NoError(t, err)
	assert.Equal(t, 0, navigator.GetCurrentIndex())

	assert.True(t, navigator.HasNext())
	assert.False(t, navigator.HasPrevious())

	err = navigator.MoveTo(1)
	assert.NoError(t, err)
	assert.Equal(t, 1, navigator.GetCurrentIndex())

	assert.False(t, navigator.HasNext())
	assert.True(t, navigator.HasPrevious())

	errors := navigator.ValidateNavigation()
	assert.Empty(t, errors)

	form, index, err := navigator.GetFormByID("form1")
	assert.NoError(t, err)
	assert.Equal(t, 0, index)
	assert.Equal(t, "form1", form.ID)

	_, _, err = navigator.GetFormByID("nonexistent")
	assert.Error(t, err)

	navigator.Reset()
	assert.Equal(t, -1, navigator.GetCurrentIndex())
	assert.False(t, navigator.HasPrevious())
}

func TestErrorHandling(t *testing.T) {

	err := bobarista.NewCupSleeveError("test-form", assert.AnError)
	assert.Contains(t, err.Error(), "test-form")
	assert.Contains(t, err.Error(), assert.AnError.Error())

	collector := bobarista.NewErrorCollector()
	assert.False(t, collector.HasErrors())

	collector.Add(assert.AnError)
	collector.AddCupSleeveError("form1", assert.AnError)
	assert.True(t, collector.HasErrors())
	assert.Len(t, collector.Errors(), 2)

	collector.Clear()
	assert.False(t, collector.HasErrors())
}

func TestColorSchemes(t *testing.T) {

	scheme, exists := bobarista.GetColorScheme("default")
	assert.True(t, exists)
	assert.Equal(t, "Default", scheme.Name)

	_, exists = bobarista.GetColorScheme("nonexistent")
	assert.False(t, exists)

	schemes := bobarista.GetAvailableColorSchemes()
	assert.Contains(t, schemes, "default")
	assert.Contains(t, schemes, "dark")
	assert.Contains(t, schemes, "ocean")

	custom := bobarista.CreateCustomColorScheme("test", "#FF0000", "#00FF00", "#0000FF")
	assert.Equal(t, "test", custom.Name)

	bobarista.RegisterColorScheme("test", custom)
	retrieved, exists := bobarista.GetColorScheme("test")
	assert.True(t, exists)
	assert.Equal(t, "test", retrieved.Name)
}

func TestFormValidation(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected bool
	}{
		{"valid input", "test", true},
		{"empty input", "", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

		})
	}
}

func TestFlowExecution(t *testing.T) {
	var name string
	var completed bool

	boba := bobarista.New("Test").
		AddForm(bobarista.NewForm("test", "Test").
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
		OnComplete(func(boba *bobarista.Bobarista) error {
			completed = true
			if completed {
				fmt.Println("Flow completed with name:", name)
			}
			return nil
		}).
		Build()

	cmd := boba.Init()
	assert.NotNil(t, cmd)

	globalData := boba.GetGlobalData()
	assert.NotNil(t, globalData)
}

func TestOnCompleteFileSaveAlternative(t *testing.T) {

	tempDir, err := os.MkdirTemp("", "cupsleeve_test")
	assert.NoError(t, err)
	defer os.RemoveAll(tempDir)

	outputFile := filepath.Join(tempDir, "form_results.json")

	var name, email string

	saveToFileHandler := func(boba *bobarista.Bobarista) error {
		globalData := boba.GetGlobalData()

		saveData := struct {
			Timestamp string            `json:"timestamp"`
			FormData  map[string]string `json:"form_data"`
		}{
			Timestamp: time.Now().Format(time.RFC3339),
			FormData:  make(map[string]string),
		}

		for key, valuePtr := range *globalData.Values {
			if valuePtr != nil {
				saveData.FormData[key] = *valuePtr
			}
		}

		jsonData, err := json.MarshalIndent(saveData, "", "  ")
		if err != nil {
			return err
		}

		return os.WriteFile(outputFile, jsonData, 0644)
	}

	boba := bobarista.New("File Save Test").
		AddForm(bobarista.NewForm("user_info", "User Information").
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
		OnComplete(saveToFileHandler).
		Build()

	name = "Jane Smith"
	email = "jane@example.com"

	globalData := boba.GetGlobalData()
	globalData.Values.Set("name", name)
	globalData.Values.Set("email", email)

	err = saveToFileHandler(boba)
	assert.NoError(t, err)

	assert.FileExists(t, outputFile)

	fileContent, err := os.ReadFile(outputFile)
	assert.NoError(t, err)

	var savedData struct {
		Timestamp string            `json:"timestamp"`
		FormData  map[string]string `json:"form_data"`
	}

	err = json.Unmarshal(fileContent, &savedData)
	assert.NoError(t, err)

	assert.NotEmpty(t, savedData.Timestamp)
	assert.Equal(t, "Jane Smith", savedData.FormData["name"])
	assert.Equal(t, "jane@example.com", savedData.FormData["email"])

	t.Logf("Saved data: %+v", savedData)
}

func TestOnCompleteWithFixtures(t *testing.T) {

	fixturesDir := "../testdata/fixtures"
	err := os.MkdirAll(fixturesDir, 0755)
	assert.NoError(t, err)

	outputFile := filepath.Join(fixturesDir, "test_output.json")
	defer os.Remove(outputFile)

	var name string = "Test User"
	var email string = "test@example.com"

	boba := bobarista.New("Fixtures Test").
		AddForm(bobarista.NewForm("info", "Info").
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
		OnComplete(func(boba *bobarista.Bobarista) error {

			globalData := boba.GetGlobalData()

			result := make(map[string]string)
			for key, valuePtr := range *globalData.Values {
				if valuePtr != nil {
					result[key] = *valuePtr
				}
			}

			jsonData, _ := json.MarshalIndent(result, "", "  ")
			return os.WriteFile(outputFile, jsonData, 0644)
		}).
		Build()

	globalData := boba.GetGlobalData()
	globalData.Values.Set("name", name)
	globalData.Values.Set("email", email)

	assert.NotNil(t, boba)
}
