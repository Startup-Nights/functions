package main

import (
	"bytes"
	"errors"
	"fmt"
	"image"
	"image/jpeg"
	"image/png"
	"io"
	"net/http"
	"os"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"golang.org/x/image/draw"
)

const (
	typePNG = iota
	typeJPEG
	typeUnsupported
)

type Request struct {
	Filename string `json:"filename"`
	Width    int    `json:"width"`
	Height   int    `json:"height"`
}

type Response struct {
	StatusCode int               `json:"statusCode,omitempty"`
	Headers    map[string]string `json:"headers,omitempty"`
	Body       ResponseData      `json:"body"`
}

type ResponseData struct {
	Download string `json:"download"`
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
	result, err := client.GetObject(&s3.GetObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(in.Filename),
	})
	if err != nil {
		return &Response{StatusCode: http.StatusInternalServerError, Body: ResponseData{Error: err.Error()}}, err
	}

	var input bytes.Buffer
	_, err = io.Copy(&input, result.Body)
	if err != nil {
		return &Response{StatusCode: http.StatusInternalServerError, Body: ResponseData{Error: err.Error()}}, err
	}

	width := in.Width
	height := in.Height

	if width == 0 {
		width = 600
	}
	if height == 0 {
		height = 300
	}

	resized, err := resize(&input, typeFromFilename(in.Filename), width, height)
	if err != nil {
		return &Response{StatusCode: http.StatusInternalServerError, Body: ResponseData{Error: err.Error()}}, err
	}

	adjusted, err := adjustCanvas(resized, width, height)
	if err != nil {
		return &Response{StatusCode: http.StatusInternalServerError, Body: ResponseData{Error: err.Error()}}, err
	}

	data, err := io.ReadAll(adjusted)
	if err != nil {
		return &Response{StatusCode: http.StatusInternalServerError, Body: ResponseData{Error: err.Error()}}, err
	}

	filename := fmt.Sprintf("%s_%dx%d.png", strings.Split(in.Filename, ".")[0], width, height)
	object := s3.PutObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(filename),
		Body:   strings.NewReader(string(data)),
		ACL:    aws.String("public-read"),
	}
	_, err = client.PutObject(&object)
	if err != nil {
		return &Response{StatusCode: http.StatusInternalServerError, Body: ResponseData{Error: err.Error()}}, err
	}

	return &Response{
		StatusCode: http.StatusOK,
		Headers: map[string]string{
			"Content-Type": "application/json",
		},
		Body: ResponseData{
			Download: fmt.Sprintf(
				"https://%s.%s.cdn.digitaloceanspaces.com/%s",
				bucket,
				region,
				filename,
			),
		},
	}, nil
}

// resizes the image to fit the canvas
func resize(in io.Reader, filetype, width, height int) (io.Reader, error) {
	var (
		output bytes.Buffer
		src    image.Image
		err    error
	)

	switch filetype {
	case typePNG:
		src, err = png.Decode(in)
	case typeJPEG:
		src, err = jpeg.Decode(in)
	default:
		return &output, errors.New("filetype not supported yet")
	}

	if err != nil {
		return &output, err
	}

	x := float64(src.Bounds().Max.X)
	y := float64(src.Bounds().Max.Y)

	factorX := x / float64(width)
	factorY := y / float64(height)

	var dst *image.RGBA

	if factorX > factorY {
		dst = image.NewRGBA(image.Rect(0, 0, width, int(y/factorX)))
	} else {
		dst = image.NewRGBA(image.Rect(0, 0, int(x/factorY), height))
	}

	draw.NearestNeighbor.Scale(dst, dst.Rect, src, src.Bounds(), draw.Over, nil)
	if err := png.Encode(&output, dst); err != nil {
		return &output, err
	}

	return &output, nil
}

// makes sure that the image is in the center of the canvas
func adjustCanvas(in io.Reader, width, height int) (io.Reader, error) {
	var output bytes.Buffer

	src, err := png.Decode(in)
	if err != nil {
		return &output, err
	}

	x := float64(src.Bounds().Max.X)
	y := float64(src.Bounds().Max.Y)

	moveX := (float64(width) - x) / 2
	moveY := (float64(height) - y) / 2

	dst := image.NewRGBA(image.Rect(0, 0, width, height))

	draw.Copy(dst, image.Pt(int(moveX), int(moveY)), src, src.Bounds(), draw.Over, nil)
	if err := png.Encode(&output, dst); err != nil {
		return &output, err
	}

	return &output, nil
}

func typeFromFilename(file string) int {
	switch strings.Split(file, ".")[1] {
	case "png":
		return typePNG
	case "jpeg":
		fallthrough
	case "jpg":
		return typeJPEG
	default:
		return typeUnsupported
	}
}
