package main

import (
	"log"

	"github.com/charmbracelet/huh"
	"github.com/choice404/bobarista/pkg/bobarista"
)

func main() {
	var name, email string

	boba := bobarista.New("Basic Example").
		WithColorScheme("sky").
		AddForm(bobarista.NewForm("info", "User Information").
			WithGenerator(func(current *bobarista.FormValues, global *bobarista.FormValues) *huh.Form {
				return huh.NewForm(
					huh.NewGroup(
						huh.NewInput().
							Title("What's your name?").
							Value(&name),
						huh.NewInput().
							Title("What's your email?").
							Value(&email),
					),
				)
			}).
			WithOnComplete(func(current *bobarista.FormData, global *bobarista.FormData) error {

				global.Values.Set("name", name)
				global.Values.Set("email", email)
				return nil
			})).
		OnComplete(func(boba *bobarista.Bobarista) error {
			log.Printf("Hello %s! Your email is %s", name, email)
			return nil
		}).
		Build()

	if err := boba.Run(); err != nil {
		log.Fatal(err)
	}
}
