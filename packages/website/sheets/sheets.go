package main

import (
	"bytes"
	"context"
	"encoding/json"

	"github.com/pkg/errors"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/option"
	"google.golang.org/api/sheets/v4"
)

type SheetService struct {
	service *sheets.Service
}

type Credentials struct {
	Config      string
	Credentials string
}

func NewSheetService(ctx context.Context, creds Credentials) (*SheetService, error) {
	svc := &SheetService{}

	config, err := google.ConfigFromJSON([]byte(creds.Config), "https://www.googleapis.com/auth/spreadsheets")
	if err != nil {
		return svc, errors.Wrap(err, "get sheets config from json")
	}

	token := &oauth2.Token{}
	if err := json.NewDecoder(bytes.NewBuffer([]byte(creds.Credentials))).Decode(token); err != nil {
		return svc, errors.Wrap(err, "decode sheets credentials")
	}

	client := config.Client(ctx, token)
	service, err := sheets.NewService(ctx, option.WithHTTPClient(client), option.WithScopes(sheets.SpreadsheetsScope))
	if err != nil {
		return svc, errors.Wrap(err, "create new sheets services")
	}

	svc.service = service

	return svc, nil
}

func (s *SheetService) Save(ctx context.Context, id, writeRange string, data []string) error {
	values := []interface{}{}

	for _, val := range data {
		values = append(values, val)
	}

	valRange := &sheets.ValueRange{}
	valRange.Values = append(valRange.Values, values)

	if _, err := s.service.Spreadsheets.Values.Append(id, writeRange, valRange).ValueInputOption("RAW").Do(); err != nil {
		return errors.Wrap(err, "append data to spreadsheet")
	}

	return nil
}
