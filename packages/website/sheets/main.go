package main

import (
	"context"
	"net/http"
	"os"
)

type Request struct {
	ID    string `json:"id"`    // sheets id
	Range string `json:"range"` // range to read or write data

	// for save requests: the data to write
	Data []string `json:"data"`

	// https://docs.digitalocean.com/products/functions/reference/parameters-responses/#event-parameter
	RequestData struct {
		Method string `json:"method"`
	} `json:"http"`
}

type Response struct {
	StatusCode int               `json:"statusCode,omitempty"`
	Headers    map[string]string `json:"headers,omitempty"`
	Body       ResponseData      `json:"body,omitempty"`
}

type ResponseData struct {
	// for read requests
	Data  [][]string `json:"data"`
	Error string     `json:"error"`
}

func Main(ctx context.Context, in Request) (*Response, error) {
	credentials := os.Getenv("CREDENTIALS")
	gmailCreds := os.Getenv("SHEETS")

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

	switch in.ID {
	case "10jYxSsWu93D-l0eybcZZmrvonmGVUTTcd8uvdsBmoU8":
		fallthrough
	case "1F4r2nCsQUIE38qOJaBzuyqHOtgVc3KshhucyOQI6zBU":
		fallthrough
	case "1WX6vvcCJihBJ9tFN-8AixYAyt5i1nSfMeX81gsEEwjs":
		fallthrough
	case "1ifvj7KmYvitjVGiKFceve7Zce9KlTijic2MDI_aXn9I":
		// TODO: this switch should be more "official"
		// maybe based on an url parameter?

		// in case there is no data, handle this as a read request
		if len(in.Data) == 0 {
			data, err := svc.Read(ctx, in.ID, in.Range)
			if err != nil {
				return &Response{
					StatusCode: http.StatusInternalServerError,
					Body:       ResponseData{Error: err.Error()},
				}, err
			}

			return &Response{
				StatusCode: http.StatusOK,
				Body: ResponseData{
					Data:  data,
					Error: "",
				},
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

	default:
		return &Response{
			StatusCode: http.StatusInternalServerError,
			Body:       ResponseData{Error: "not authorized"},
		}, err
	}
}
