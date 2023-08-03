package main

import (
	"context"
	"net/http"
	"os"
)

type Request struct {
	Recipient string `json:"recipient"`
	Title     string `json:"title"`
	Content   string `json:"content"`
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
	credentials := os.Getenv("CREDENTIALS")
	gmailCreds := os.Getenv("GMAIL")

	ctx := context.Background()

	svc, err := NewMailService(ctx, Credentials{
		Config:      credentials,
		Credentials: gmailCreds,
	})
	if err != nil {
		return &Response{
			StatusCode: http.StatusInternalServerError,
			Body:       ResponseData{Error: err.Error()},
		}, err
	}

	if err := svc.Send(
		ctx,
		in.Recipient,
		in.Title,
		in.Content,
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
