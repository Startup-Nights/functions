package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/pkg/errors"
	"github.com/slack-go/slack"
)

var (
	webhook, credentials, sheetsCreds string
)

func init() {
	if webhook = os.Getenv("WEBHOOK"); webhook == "" {
		panic("no webhook configured")
	}
	if credentials = os.Getenv("CREDENTIALS"); credentials == "" {
		panic("no credentials configured")
	}
	if sheetsCreds = os.Getenv("SHEETS"); sheetsCreds == "" {
		panic("no sheets credentials configured")
	}
}

type Request struct {
	Firstname        string   `json:"firstname"`
	Lastname         string   `json:"lastname"`
	Email            string   `json:"email"`
	Startup          string   `json:"startup"`
	Website          string   `json:"website"`
	Pitchdeck        string   `json:"pitchdeck"`
	Stage            string   `json:"seed"`
	Problem          string   `json:"problem"`
	Approach         string   `json:"approach"`
	Unique           string   `json:"unique"`
	User             string   `json:"user"`
	RaisingFunds     string   `json:"raising"`
	AlreadyPitched   string   `json:"pitched"`
	BusinessModel    string   `json:"business_model"`
	LinkedinProfiles []string `json:"linkedin"`
}

type Response struct {
	StatusCode int               `json:"statusCode,omitempty"`
	Headers    map[string]string `json:"headers,omitempty"`
	Body       ResponseData      `json:"body,omitempty"`
}

type ResponseData struct {
	Error string `json:"error"`
}

func Main(in Request) (*Response, error) {
	ctx := context.Background()

	// backup - send the data in slack in case the google tokens are not valid
	// anymore or something else happens
	if err := sendSlackMessage(in); err != nil {
		return &Response{
			StatusCode: http.StatusInternalServerError,
			Headers:    map[string]string{"Content-Type": "application/json"},
			Body: ResponseData{
				Error: fmt.Sprintf("post slack message: %v", err),
			},
		}, nil
	}

	// try to save the data to a google sheet
	if err := saveToSheets(ctx, in); err != nil {
		return &Response{
			StatusCode: http.StatusInternalServerError,
			Headers:    map[string]string{"Content-Type": "application/json"},
			Body: ResponseData{
				Error: fmt.Sprintf("save data to sheets: %v", err),
			},
		}, nil
	}

	return &Response{
		StatusCode: http.StatusOK,
		Headers:    map[string]string{"Content-Type": "application/json"},
		Body:       ResponseData{},
	}, nil
}

func sendSlackMessage(in Request) error {
	// indent data for better readability in slack
	data, err := json.MarshalIndent(in, "", "  ")
	if err != nil {
		return errors.Wrap(err, "indent json")
	}

	if err := slack.PostWebhook(webhook, &slack.WebhookMessage{
		Channel: "10_01_booth_notification",
		// TODO: update the spreadsheet link
		Text: fmt.Sprintf(
			"new submission: %s\n\n%s",
			"https://docs.google.com/spreadsheets/d/1WX6vvcCJihBJ9tFN-8AixYAyt5i1nSfMeX81gsEEwjs/edit#gid=0",
			string(data),
		),
	}); err != nil {
		return errors.Wrap(err, "send message to channel")
	}

	return nil
}
