package main

import (
	"context"
	"net/http"
	"os"
)

type Request struct {
	ID    string   `json:"id"`    // sheets id
	Range string   `json:"range"` // range to write data
	Data  []string `json:"data"`  // the data
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
	gmailCreds := os.Getenv("SHEETS")

	ctx := context.Background()

	svc, err := NewSheetService(ctx, Credentials{
		Config:      credentials,
		Credentials: gmailCreds,
	})
	if err != nil {
		return &Response{
			StatusCode: http.StatusInternalServerError,
			Body:       ResponseData{Error: err.Error()},
		}, err
	}

	if err := svc.Save(ctx, in.ID, in.Range, in.Data); err != nil {
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
