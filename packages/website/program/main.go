package main

import (
	"net/http"

	"golang.org/x/net/html"
)

type Request struct {
	Url string `json:"url"`
}

type Response struct {
	StatusCode int               `json:"statusCode,omitempty"`
	Headers    map[string]string `json:"headers,omitempty"`
	Body       ResponseData      `json:"body"`
}

type ResponseData struct {
	Data  []string `json:"data"`
	Error string   `json:"error"`
}

var items = []string{}

func Main(in Request) (*Response, error) {
	res, err := http.Get(in.Url)
	if err != nil {
		return &Response{Body: ResponseData{
			Error: err.Error(),
		}}, err
	}

	doc, err := html.Parse(res.Body)
	if err != nil {
		return &Response{Body: ResponseData{
			Error: err.Error(),
		}}, err
	}

	// TODO: filter the document for program entries
	// TODO: convert to a sensible format
	extractItems(doc)

	return &Response{
		Body: ResponseData{
			Data: items,
		},
	}, nil
}

func extractItems(n *html.Node) {
	if n.Type == html.ElementNode && n.Data == "a" {
		items = append(items, n.Data)
	}

	for c := n.FirstChild; c != nil; c = c.NextSibling {
		extractItems(c)
	}
}
