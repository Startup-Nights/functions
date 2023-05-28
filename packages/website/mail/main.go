package main

import (
	"net/http"
)

// Request takes in the user's input for the filename they want and if the type is a GET or PUT.
type Request struct {
	// Filename is the name of the file that will be uploaded or downloaded.
	Filename string `json:"filename"`
}

// Response returns back the http code, type of data, and the presigned url to the user.
type Response struct {
	// StatusCode is the http code that will be returned back to the user.
	StatusCode int `json:"statusCode,omitempty"`
	// Headers is the information about the type of data being returned back.
	Headers map[string]string `json:"headers,omitempty"`
	// Body will contain the presigned url to upload or download files.
	Body string `json:"body,omitempty"`
}

func Main(in Request) (*Response, error) {
	return &Response{
		StatusCode: http.StatusOK,
		Body:       "",
	}, nil
}
