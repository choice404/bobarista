package main

import (
	"fmt"
	"log"

	"github.com/charmbracelet/huh"
	"github.com/choice404/bobarista/pkg/bobarista"
)

func main() {
	var (
		name       string
		email      string
		age        string
		terms      bool
		newsletter bool
	)

	boba := bobarista.New("User Registration").
		AddForm(bobarista.NewForm("personal", "Personal Information").
			WithGenerator(func(current *bobarista.FormValues, global *bobarista.FormValues) *huh.Form {
				return huh.NewForm(
					huh.NewGroup(
						huh.NewInput().
							Title("Full Name").
							Value(&name),
						huh.NewInput().
							Title("Email Address").
							Value(&email),
						huh.NewInput().
							Title("Age").
							Value(&age).
							Validate(func(value string) error {
								if value == "" {
									return fmt.Errorf("age cannot be empty")
								}

								if _, err := fmt.Sscanf(value, "%d", new(int)); err != nil {
									return fmt.Errorf("age must be a number")
								}
								return nil
							}),
					),
				)
			}).
			WithOnComplete(func(current *bobarista.FormData, global *bobarista.FormData) error {

				global.Values.Set("full_name", name)
				global.Values.Set("email_address", email)
				global.Values.Set("age", age)
				return nil
			})).
		AddForm(bobarista.NewForm("terms", "Terms and Conditions").
			WithGenerator(func(current *bobarista.FormValues, global *bobarista.FormValues) *huh.Form {
				return huh.NewForm(
					huh.NewGroup(
						huh.NewConfirm().
							Title("Do you accept the terms and conditions?").
							Value(&terms),
						huh.NewConfirm().
							Title("Subscribe to newsletter?").
							Value(&newsletter),
					),
				)
			}).
			WithOnComplete(func(current *bobarista.FormData, global *bobarista.FormData) error {

				global.Values.Set("terms_accepted", fmt.Sprintf("%t", terms))
				global.Values.Set("newsletter_subscription", fmt.Sprintf("%t", newsletter))
				return nil
			})).
		OnComplete(func(boba *bobarista.Bobarista) error {
			fmt.Printf("Registration completed!\n")
			fmt.Printf("Name: %s\n", name)
			fmt.Printf("Email: %s\n", email)
			fmt.Printf("Age: %s\n", age)
			fmt.Printf("Terms accepted: %t\n", terms)
			fmt.Printf("Newsletter: %t\n", newsletter)
			return nil
		}).
		Build()

	if err := boba.Run(); err != nil {
		log.Fatal(err)
	}
}
