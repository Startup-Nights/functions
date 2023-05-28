package main

import (
	"encoding/json"
	"net/http"
	"os"

	"github.com/slack-go/slack"
)

var (
	webhook string
)

func init() {
	if webhook = os.Getenv("WEBHOOK"); webhook == "" {
		panic("no webhook configured")
	}
}

type Request struct {
	Company struct {
		Name                 string   `json:"name"`
		Website              string   `json:"website"`
		FoundingDate         string   `json:"founding_date"`
		LinkedIn             []string `json:"linkedin"`
		Employees            string   `json:"employees"`
		Pitch                string   `json:"pitch"`
		Categories           []string `json:"categories"`
		AdditionalCategories string   `json:"additional_categories"`
		Logo                 string   `json:"logo"`
		Address              struct {
			Street  string `json:"street"`
			ZIP     string `json:"zip"`
			City    string `json:"city"`
			Country string `json:"country"`
		} `json:"address"`
		BillingAddress struct {
			Street  string `json:"street"`
			ZIP     string `json:"zip"`
			City    string `json:"city"`
			Country string `json:"country"`
		} `json:"address_billing"`
	} `json:"company"`
	Contact struct {
		FirstName string `json:"firstname"`
		LastName  string `json:"lastname"`
		Email     string `json:"email"`
		Phone     string `json:"phone"`
		Role      string `json:"role"`
	} `json:"contact"`
	Varia struct {
		Package struct {
			Title string `json:"title"`
		} `json:"package"`
		Formats      []string `json:"formats"`
		Accomodation string   `json:"accomodation"`
		Ukraine      string   `json:"ukraine"`
		Equipment    string   `json:"equipment"`
	} `json:"varia"`
}

type Response struct {
	StatusCode int               `json:"statusCode,omitempty"`
	Headers    map[string]string `json:"headers,omitempty"`
	Body       string            `json:"body,omitempty"`
}

func Main(in Request) (*Response, error) {
	data, err := json.MarshalIndent(in, "", "  ")
	if err != nil {
		return &Response{StatusCode: http.StatusInternalServerError, Body: "{}"}, nil
	}

	if err := slack.PostWebhook(webhook, &slack.WebhookMessage{
		Channel: "10_01_booth_notification",
		Text:    string(data),
	}); err != nil {
		return &Response{StatusCode: http.StatusInternalServerError, Body: "{}"}, nil
	}

	return &Response{
		StatusCode: http.StatusOK,
		Headers: map[string]string{
			"Content-Type": "application/json",
		},
		Body: "{}",
	}, nil
}
