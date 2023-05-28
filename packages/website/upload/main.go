package main

import (
	"fmt"
	"net/http"
	"net/url"
	"os"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
)

var (
	key, secret, region, bucket string
)

func init() {
	if bucket = os.Getenv("BUCKET"); bucket == "" {
		panic("no bucket provided")
	}
	if region = os.Getenv("REGION"); region == "" {
		panic("no region provided")
	}
	if key = os.Getenv("KEY"); key == "" {
		panic("no key provided")
	}
	if secret = os.Getenv("SECRET"); secret == "" {
		panic("no secret provided")
	}
}

type Request struct {
	Filename string `json:"filename"`
}

type Response struct {
	StatusCode int               `json:"statusCode,omitempty"`
	Headers    map[string]string `json:"headers,omitempty"`
	Body       string            `json:"body,omitempty"`
}

func Main(in Request) (*Response, error) {
	config := &aws.Config{
		Credentials: credentials.NewStaticCredentials(key, secret, ""),
		Endpoint:    aws.String(fmt.Sprintf("%s.digitaloceanspaces.com:443", region)),
		Region:      aws.String(region),
	}
	session, err := session.NewSession(config)
	if err != nil {
		return &Response{StatusCode: http.StatusInternalServerError, Body: ""}, err
	}

	client := s3.New(session)

	// make sure that the file has a solid name and does not overwrite others
	// in case multiple people use something like 'logo.png'
	name := url.QueryEscape(in.Filename)
	key := fmt.Sprintf("%d/%s/%d-%s", time.Now().Year(), "startups", time.Now().Nanosecond(), name)

	object := s3.PutObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
		ACL:    aws.String("public-read"),
	}

	req, _ := client.PutObjectRequest(&object)
	url, err := req.Presign(5 * time.Minute)
	if err != nil {
		return &Response{StatusCode: http.StatusInternalServerError, Body: ""}, err
	}

	return &Response{
		StatusCode: http.StatusOK,
		Headers: map[string]string{
			"Content-Type": "application/json",
		},
		Body: fmt.Sprintf(`{"url":"%s"}`, url),
	}, nil
}
