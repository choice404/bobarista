package main

import (
	"fmt"
	"log"

	"github.com/charmbracelet/huh"
	"github.com/choice404/bobarista/pkg/bobarista"
)

func main() {
	var (
		userType     string
		companyName  string
		personalName string
	)

	boba := bobarista.New("Advanced Example").
		WithDebug(true).
		AddForm(bobarista.NewForm("type", "User Type").
			WithGenerator(func(current *bobarista.FormValues, global *bobarista.FormValues) *huh.Form {
				return huh.NewForm(
					huh.NewGroup(
						huh.NewSelect[string]().
							Title("Are you signing up as?").
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
		AddForm(bobarista.NewForm("company", "Company Information").
			WithSkipCondition(func(current *bobarista.FormData, global *bobarista.FormData) bool {
				if val, exists := global.Values.Get("user_type"); exists {
					shouldSkip := val != "company"
					return shouldSkip
				}
				return true
			}).
			WithGenerator(func(current *bobarista.FormValues, global *bobarista.FormValues) *huh.Form {
				fmt.Printf("DEBUG Company Generator: Called\n")
				return huh.NewForm(
					huh.NewGroup(
						huh.NewInput().
							Title("Company Name").
							Value(&companyName),
					),
				)
			}).
			WithOnComplete(func(current *bobarista.FormData, global *bobarista.FormData) error {
				global.Values.Set("company_name", companyName)
				return nil
			})).
		AddForm(bobarista.NewForm("personal", "Personal Information").
			WithSkipCondition(func(current *bobarista.FormData, global *bobarista.FormData) bool {
				if val, exists := global.Values.Get("user_type"); exists {
					shouldSkip := val != "individual"
					return shouldSkip
				}
				return true
			}).
			WithGenerator(func(current *bobarista.FormValues, global *bobarista.FormValues) *huh.Form {
				return huh.NewForm(
					huh.NewGroup(
						huh.NewInput().
							Title("Full Name").
							Value(&personalName),
					),
				)
			}).
			WithOnComplete(func(current *bobarista.FormData, global *bobarista.FormData) error {
				global.Values.Set("personal_name", personalName)
				return nil
			})).
		OnComplete(func(boba *bobarista.Bobarista) error {
			fmt.Printf("DEBUG Final: userType variable = %s\n", userType)
			if userType == "company" {
				fmt.Printf("Company: %s\n", companyName)
			} else {
				fmt.Printf("Individual: %s\n", personalName)
			}
			return nil
		}).
		Build()

	if err := boba.Run(); err != nil {
		log.Fatal(err)
	}
}
