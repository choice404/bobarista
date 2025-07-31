package main

import (
	"fmt"
	"log"

	"github.com/charmbracelet/huh"
	"github.com/choice404/bobarista/pkg/bobarista"
)

func main() {
	var (
		projectType string
		projectName string
		language    string
		framework   string
		database    string
	)

	boba := bobarista.New("Project Setup Wizard").
		AddForm(bobarista.NewForm("type", "Project Type").
			WithGenerator(func(current *bobarista.FormValues, global *bobarista.FormValues) *huh.Form {
				return huh.NewForm(
					huh.NewGroup(
						huh.NewSelect[string]().
							Title("What type of project are you creating?").
							Options(
								huh.NewOption("Web Application", "web"),
								huh.NewOption("CLI Tool", "cli"),
								huh.NewOption("Library", "library"),
								huh.NewOption("API Service", "api"),
							).
							Value(&projectType),
					),
				)
			}).
			WithOnComplete(func(current *bobarista.FormData, global *bobarista.FormData) error {
				global.Values.Set("project_type", projectType)
				return nil
			})).
		AddForm(bobarista.NewForm("details", "Project Details").
			WithGenerator(func(current *bobarista.FormValues, global *bobarista.FormValues) *huh.Form {
				return huh.NewForm(
					huh.NewGroup(
						huh.NewInput().
							Title("Project Name").
							Value(&projectName),
						huh.NewSelect[string]().
							Title("Programming Language").
							Options(
								huh.NewOption("Go", "go"),
								huh.NewOption("Python", "python"),
								huh.NewOption("JavaScript", "javascript"),
								huh.NewOption("TypeScript", "typescript"),
							).
							Value(&language),
					),
				)
			}).
			WithOnComplete(func(current *bobarista.FormData, global *bobarista.FormData) error {
				global.Values.Set("project_name", projectName)
				global.Values.Set("language", language)
				return nil
			})).
		AddForm(bobarista.NewForm("framework", "Framework Selection").
			WithSkipCondition(func(current *bobarista.FormData, global *bobarista.FormData) bool {
				if val, exists := global.Values.Get("project_type"); exists {
					return val == "cli"
				}
				return false
			}).
			WithGenerator(func(current *bobarista.FormValues, global *bobarista.FormValues) *huh.Form {
				var options []huh.Option[string]

				if langVal, exists := global.Get("language"); exists {
					switch langVal {
					case "go":
						options = []huh.Option[string]{
							huh.NewOption("Gin", "gin"),
							huh.NewOption("Echo", "echo"),
							huh.NewOption("Fiber", "fiber"),
						}
					case "python":
						options = []huh.Option[string]{
							huh.NewOption("Django", "django"),
							huh.NewOption("Flask", "flask"),
							huh.NewOption("FastAPI", "fastapi"),
						}
					default:
						options = []huh.Option[string]{
							huh.NewOption("None", "none"),
						}
					}
				} else {
					options = []huh.Option[string]{
						huh.NewOption("None", "none"),
					}
				}

				return huh.NewForm(
					huh.NewGroup(
						huh.NewSelect[string]().
							Title("Choose a framework").
							Options(options...).
							Value(&framework),
					),
				)
			}).
			WithOnComplete(func(current *bobarista.FormData, global *bobarista.FormData) error {
				global.Values.Set("framework", framework)
				return nil
			})).
		AddForm(bobarista.NewForm("database", "Database Selection").
			WithSkipCondition(func(current *bobarista.FormData, global *bobarista.FormData) bool {
				if val, exists := global.Values.Get("project_type"); exists {
					return val == "library"
				}
				return false
			}).
			WithGenerator(func(current *bobarista.FormValues, global *bobarista.FormValues) *huh.Form {
				return huh.NewForm(
					huh.NewGroup(
						huh.NewSelect[string]().
							Title("Choose a database").
							Options(
								huh.NewOption("PostgreSQL", "postgresql"),
								huh.NewOption("MySQL", "mysql"),
								huh.NewOption("SQLite", "sqlite"),
								huh.NewOption("MongoDB", "mongodb"),
								huh.NewOption("None", "none"),
							).
							Value(&database),
					),
				)
			}).
			WithOnComplete(func(current *bobarista.FormData, global *bobarista.FormData) error {
				global.Values.Set("database", database)
				return nil
			})).
		OnComplete(func(boba *bobarista.Bobarista) error {
			fmt.Printf("Project setup completed!\n")
			fmt.Printf("Type: %s\n", projectType)
			fmt.Printf("Name: %s\n", projectName)
			fmt.Printf("Language: %s\n", language)
			if framework != "" {
				fmt.Printf("Framework: %s\n", framework)
			}
			if database != "" {
				fmt.Printf("Database: %s\n", database)
			}
			return nil
		}).
		Build()

	if err := boba.Run(); err != nil {
		log.Fatal(err)
	}
}
