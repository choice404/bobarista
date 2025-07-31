package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/charmbracelet/huh"
	"github.com/choice404/bobarista/pkg/bobarista"
)

func main() {
	var (
		satisfaction string
		recommend    bool
		feedback     string
	)

	boba := bobarista.New("Customer Survey").
		AddForm(bobarista.NewForm("rating", "Satisfaction Rating").
			WithGenerator(func(current *bobarista.FormValues, global *bobarista.FormValues) *huh.Form {
				return huh.NewForm(
					huh.NewGroup(
						huh.NewSelect[string]().
							Title("How satisfied are you with our service?").
							Options(
								huh.NewOption("Very Satisfied", "very_satisfied"),
								huh.NewOption("Satisfied", "satisfied"),
								huh.NewOption("Neutral", "neutral"),
								huh.NewOption("Dissatisfied", "dissatisfied"),
								huh.NewOption("Very Dissatisfied", "very_dissatisfied"),
							).
							Value(&satisfaction),
					),
				)
			}).
			WithOnComplete(func(current *bobarista.FormData, global *bobarista.FormData) error {

				global.Values.Set("satisfaction_rating", satisfaction)
				return nil
			})).
		AddForm(bobarista.NewForm("recommendation", "Recommendation").
			WithGenerator(func(current *bobarista.FormValues, global *bobarista.FormValues) *huh.Form {
				return huh.NewForm(
					huh.NewGroup(
						huh.NewConfirm().
							Title("Would you recommend us to others?").
							Value(&recommend),
					),
				)
			}).
			WithOnComplete(func(current *bobarista.FormData, global *bobarista.FormData) error {

				global.Values.Set("would_recommend", fmt.Sprintf("%t", recommend))
				return nil
			})).
		AddForm(bobarista.NewForm("feedback", "Additional Feedback").
			WithGenerator(func(current *bobarista.FormValues, global *bobarista.FormValues) *huh.Form {
				return huh.NewForm(
					huh.NewGroup(
						huh.NewText().
							Title("Any additional feedback?").
							Value(&feedback),
					),
				)
			}).
			WithOnComplete(func(current *bobarista.FormData, global *bobarista.FormData) error {

				global.Values.Set("additional_feedback", feedback)
				return nil
			})).
		OnComplete(func(boba *bobarista.Bobarista) error {

			if err := saveSurveyResults(boba); err != nil {
				fmt.Printf("Error saving survey results: %v\n", err)
				return err
			}

			fmt.Printf("Survey completed!\n")
			fmt.Printf("Satisfaction: %s\n", satisfaction)
			fmt.Printf("Would recommend: %t\n", recommend)
			fmt.Printf("Feedback: %s\n", feedback)
			fmt.Printf("\nResults saved to survey_results.json\n")
			return nil
		}).
		Build()

	if err := boba.Run(); err != nil {
		log.Fatal(err)
	}
}

type SurveyResult struct {
	Timestamp          string `json:"timestamp"`
	SatisfactionRating string `json:"satisfaction_rating"`
	WouldRecommend     string `json:"would_recommend"`
	AdditionalFeedback string `json:"additional_feedback"`
	SurveyID           string `json:"survey_id"`
}

func saveSurveyResults(boba *bobarista.Bobarista) error {
	globalData := boba.GetGlobalData()

	result := SurveyResult{
		Timestamp: time.Now().Format(time.RFC3339),
		SurveyID:  fmt.Sprintf("survey_%d", time.Now().Unix()),
	}

	if val, exists := globalData.Values.Get("satisfaction_rating"); exists {
		result.SatisfactionRating = val
	}

	if val, exists := globalData.Values.Get("would_recommend"); exists {
		result.WouldRecommend = val
	}

	if val, exists := globalData.Values.Get("additional_feedback"); exists {
		result.AdditionalFeedback = val
	}

	jsonData, err := json.MarshalIndent(result, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal survey data: %w", err)
	}

	filename := "survey_results.json"
	err = os.WriteFile(filename, jsonData, 0644)
	if err != nil {
		return fmt.Errorf("failed to write survey results to file: %w", err)
	}

	return nil
}
