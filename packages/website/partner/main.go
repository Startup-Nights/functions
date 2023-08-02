package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"text/template"

	"github.com/pkg/errors"
	"github.com/slack-go/slack"
)

var (
	webhook, credentials, gmailCreds string
	tpl                              *template.Template
)

const (
	// mime data for the email header
	mime = "MIME-version: 1.0;\nContent-Type: text/plain; charset=\"UTF-8\";"
)

func init() {
	if webhook = os.Getenv("WEBHOOK"); webhook == "" {
		panic("no webhook configured")
	}
	if credentials = os.Getenv("CREDENTIALS"); credentials == "" {
		panic("no webhook configured")
	}
	if gmailCreds = os.Getenv("GMAIL"); gmailCreds == "" {
		panic("no webhook configured")
	}
	tpl = template.Must(template.New("booth registration mail").Parse(BoothRegistrationTemplate))
}

// TODO: adjust the data structure
type Address struct {
	Street  string `json:"street"`
	ZIP     string `json:"zip"`
	City    string `json:"city"`
	Country string `json:"country"`
}

type Company struct {
	Name                 string   `json:"name"`
	Website              string   `json:"website"`
	FoundingDate         string   `json:"founding_date"`
	LinkedIn             []string `json:"linkedin"`
	Employees            string   `json:"employees"`
	Pitch                string   `json:"pitch"`
	Categories           []string `json:"categories"`
	AdditionalCategories string   `json:"additional_categories"`
	Logo                 string   `json:"logo"`
	Address              Address  `json:"address"`
	BillingAddress       Address  `json:"address_billing"`
}

type Contact struct {
	FirstName string `json:"firstname"`
	LastName  string `json:"lastname"`
	Email     string `json:"email"`
	Phone     string `json:"phone"`
	Role      string `json:"role"`
}

type Varia struct {
	Package         Package  `json:"package"`
	Formats         []string `json:"formats"`
	Accomodation    string   `json:"accomodation"`
	PreviousVisitor string   `json:"previous_visitor"`
	Referral        string   `json:"referral"`
	Equipment       string   `json:"equipment"`
}

type Package struct {
	Title string `json:"title"`
}

type Request struct {
	Company Company `json:"company"`
	Contact Contact `json:"contact"`
	Varia   Varia   `json:"varia"`
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

	// send the mail to the applicant
	if err := sendMail(ctx, in); err != nil {
		return &Response{
			StatusCode: http.StatusInternalServerError,
			Headers:    map[string]string{"Content-Type": "application/json"},
			Body: ResponseData{
				Error: fmt.Sprintf("send mail: %v", err),
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
		// TODO: adjust text
		Text: fmt.Sprintf(
			"new submission: %s\n\n%s",
			string(data),
		),
	}); err != nil {
		return errors.Wrap(err, "send message to channel")
	}

	return nil
}
