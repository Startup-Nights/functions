package main

import (
	"context"
	"net/http"
	"os"

	"github.com/startup-nights/functions/packages/lib"
)

var (
	webhook, credentials, gmailCreds string
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
}

type Request struct {
	Firstname string   `json:"firstname"`
	Lastname  string   `json:"lastname"`
	Company   string   `json:"company"`
	Interests []string `json:"interests"`
	Email     string   `json:"email"`
	Budget    string   `json:"budget"`
	Idea      string   `json:"idea"`
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

	// backup stuff to slack just to be sure
	if err := lib.SendToSlack(
		webhook,
		"10_bot_notifications",
		"New Partner Registration - "+in.Company+"!",
	); err != nil {
		return &Response{
			StatusCode: http.StatusInternalServerError,
			Body:       ResponseData{Error: err.Error()},
		}, err
	}

	svc, err := lib.NewMailService(ctx, lib.Credentials{
		Config:      credentials,
		Credentials: gmailCreds,
	})
	if err != nil {
		return &Response{
			StatusCode: http.StatusInternalServerError,
			Body:       ResponseData{Error: err.Error()},
		}, err
	}

	// send the mail to the partner team and maybe even a confirmation mail to
	// the person that signed up
	if err := svc.Send(
		ctx,
		"Partner Registration",
		"mischa@ninetyfour.ch",
		"Neue Partner Registration von Startup Nights",
	); err != nil {
		return &Response{
			StatusCode: http.StatusInternalServerError,
			Body:       ResponseData{Error: err.Error()},
		}, err
	}

	return &Response{
		StatusCode: http.StatusOK,
		Headers:    map[string]string{"Content-Type": "application/json"},
		Body:       ResponseData{},
	}, nil
}
