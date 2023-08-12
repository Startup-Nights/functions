package main

import (
	"fmt"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
)

type Request struct {
	Filename string `json:"filename"`
}

type Response struct {
	StatusCode int               `json:"statusCode,omitempty"`
	Headers    map[string]string `json:"headers,omitempty"`
	Body       ResponseData      `json:"body"`
}

type ResponseData struct {
	Upload   string `json:"upload"`   // presigned upload url
	Download string `json:"download"` // convenience: direct url for download
	Filename string `json:"filename"` // convenience: only the filename
	Error    string `json:"error"`
}

func Main(in Request) (*Response, error) {
	bucket := os.Getenv("BUCKET")
	region := os.Getenv("REGION")
	secret := os.Getenv("SECRET")
	key := os.Getenv("KEY")

	config := &aws.Config{
		Credentials: credentials.NewStaticCredentials(key, secret, ""),
		Endpoint:    aws.String(fmt.Sprintf("%s.digitaloceanspaces.com:443", strings.TrimSpace(region))),
		Region:      aws.String(region),
	}
	session, err := session.NewSession(config)
	if err != nil {
		return &Response{StatusCode: http.StatusInternalServerError, Body: ResponseData{Error: err.Error()}}, err
	}

	client := s3.New(session)

	// make sure that the file has a solid name and does not overwrite others
	// in case multiple people use something like 'logo.png'
	name := url.QueryEscape(in.Filename)
	filename := fmt.Sprintf("%d/%s/%d-%s", time.Now().Year(), "startups", time.Now().Nanosecond(), name)

	uploadReq, _ := client.PutObjectRequest(&s3.PutObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(filename),
		ACL:    aws.String("public-read"),
	})
	uploadURL, err := uploadReq.Presign(5 * time.Minute)
	if err != nil {
		return &Response{StatusCode: http.StatusInternalServerError, Body: ResponseData{Error: err.Error()}}, err
	}

	return &Response{
		StatusCode: http.StatusOK,
		Headers: map[string]string{
			"Content-Type": "application/json",
		},
		Body: ResponseData{
			Upload:   uploadURL,
			Filename: filename,
			Download: fmt.Sprintf(
				"https://%s.%s.cdn.digitaloceanspaces.com/%s",
				bucket,
				region,
				filename,
			),
		},
	}, nil
}
