package main

import (
	"fmt"
	"log"

	"github.com/charmbracelet/huh"
	"github.com/choice404/bobarista/pkg/bobarista"
)

func main() {
	var selectedTheme string
	var name string

	themes := bobarista.GetAvailableColorSchemes()
	themeOptions := make([]huh.Option[string], len(themes))
	for i, theme := range themes {
		themeOptions[i] = huh.NewOption(theme, theme)
	}

	boba := bobarista.New("Theme Demo").
		AddForm(bobarista.NewForm("theme", "Choose Theme").
			WithGenerator(func(current *bobarista.FormValues, global *bobarista.FormValues) *huh.Form {
				return huh.NewForm(
					huh.NewGroup(
						huh.NewSelect[string]().
							Title("Select a color scheme").
							Options(themeOptions...).
							Value(&selectedTheme),
					),
				)
			}).
			WithOnComplete(func(current *bobarista.FormData, global *bobarista.FormData) error {
				global.Values.Set("theme", selectedTheme)
				return nil
			})).
		AddForm(bobarista.NewForm("info", "User Info").
			WithGenerator(func(current *bobarista.FormValues, global *bobarista.FormValues) *huh.Form {
				return huh.NewForm(
					huh.NewGroup(
						huh.NewInput().
							Title("Enter your name").
							Value(&name),
					),
				)
			}).
			WithOnComplete(func(current *bobarista.FormData, global *bobarista.FormData) error {
				global.Values.Set("name", name)
				return nil
			})).
		WithColorScheme("default").
		OnComplete(func(boba *bobarista.Bobarista) error {
			fmt.Printf("Theme: %s, Name: %s\n", selectedTheme, name)
			return nil
		}).
		Build()

	if err := boba.Run(); err != nil {
		log.Fatal(err)
	}
}
