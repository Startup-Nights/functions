package main

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"text/template"

	"github.com/slack-go/slack"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/gmail/v1"
	"google.golang.org/api/option"
)

var (
	webhook, credentials, gmailCreds, sheetsCreds string
	tpl                                           *template.Template
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
	if sheetsCreds = os.Getenv("SHEETS"); sheetsCreds == "" {
		panic("no webhook configured")
	}
	tpl = template.Must(template.New("booth registration mail").Parse(BoothRegistrationTemplate))
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
	Body       ResponseData      `json:"body,omitempty"`
}

type ResponseData struct {
	Error string `json:"error"`
}

func Main(in Request) (*Response, error) {
	ctx := context.Background()

	// indent data for better readability in slack
	data, err := json.MarshalIndent(in, "", "  ")
	if err != nil {
		return &Response{StatusCode: http.StatusInternalServerError, Body: ResponseData{
			Error: err.Error(),
		}}, nil
	}

	// backup - send the data in slack in case the google tokens are not valid
	// anymore
	if err := slack.PostWebhook(webhook, &slack.WebhookMessage{
		Channel: "10_01_booth_notification",
		Text:    string(data),
	}); err != nil {
		return &Response{StatusCode: http.StatusInternalServerError, Body: ResponseData{
			Error: err.Error(),
		}}, nil
	}

	config, err := google.ConfigFromJSON([]byte(credentials), gmail.GmailComposeScope)
	if err != nil {
		return &Response{StatusCode: http.StatusInternalServerError, Body: ResponseData{
			Error: err.Error(),
		}}, nil
	}

	token := &oauth2.Token{}
	if err := json.NewDecoder(bytes.NewBufferString(gmailCreds)).Decode(token); err != nil {
		return &Response{StatusCode: http.StatusInternalServerError, Body: ResponseData{
			Error: err.Error(),
		}}, nil
	}

	client := config.Client(ctx, token)
	service, err := gmail.NewService(ctx, option.WithHTTPClient(client))
	if err != nil {
		return &Response{StatusCode: http.StatusInternalServerError, Body: ResponseData{
			Error: err.Error(),
		}}, nil
	}

	var message gmail.Message

	buf := bytes.NewBuffer([]byte{})
	if err := tpl.Execute(buf, in); err != nil {
		return &Response{StatusCode: http.StatusInternalServerError, Body: ResponseData{
			Error: err.Error(),
		}}, nil
	}

	// TODO: might have to check where linebreaks should be an where not
	msg := fmt.Sprintf(
		"To: %s\r\nSubject: %s\n%s\n\n\n%s",
		in.Contact.Email,
		"Startup Nights 2023 Booth Application",
		mime,
		buf.String(),
	)

	message.Raw = base64.URLEncoding.EncodeToString([]byte(msg))

	if _, err := service.Users.Messages.Send("me", &message).Do(); err != nil {
		return &Response{StatusCode: http.StatusInternalServerError, Body: ResponseData{
			Error: err.Error(),
		}}, nil
	}

	return &Response{
		StatusCode: http.StatusOK,
		Headers: map[string]string{
			"Content-Type": "application/json",
		},
		Body: ResponseData{},
	}, nil
}

const BoothRegistrationTemplate = `Hi {{.Contact.FirstName}} {{.Contact.LastName}},

thank you for the registration for a booth for {{.Company.Name}} at the Startup Nights 2023 this November. We will reach out to you soon to confirm your registration.
 
{{if .Varia.Formats }}
Beside the booth, you were interested in:
{{range .Varia.Formats }} - {{.Name}}
{{end}}
We are going to reach out to you regarding the registration for these formats separately.
{{end}}

We know, it's a long time until November so we have some suggestions for you on how you can pass the time ⏱
- Talk to other founders / startups and invite people to the event 🚀
- Buy the tickets for you and your team and don't forget to use your discount code 20OFFspecial to get 20% off 💸 You can buy the tickets here: https://portal.startup-nights.ch
- Read the FAQ if you have open questions: https://startup-nights.ch/faq

Talk to you soon 😉

Best regards,
the Startup Nights Team`
