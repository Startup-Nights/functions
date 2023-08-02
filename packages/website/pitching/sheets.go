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

func saveToSheets(ctx context.Context, in Request) error {
	config, err := google.ConfigFromJSON([]byte(credentials), "https://www.googleapis.com/auth/spreadsheets")
	if err != nil {
		return errors.Wrap(err, "get sheets config from json")
	}

	token := &oauth2.Token{}
	if err := json.NewDecoder(bytes.NewBuffer([]byte(sheetsCreds))).Decode(token); err != nil {
		return errors.Wrap(err, "decode sheets credentials")
	}

	client := config.Client(ctx, token)
	service, err := sheets.NewService(ctx, option.WithHTTPClient(client), option.WithScopes(sheets.SpreadsheetsScope))
	if err != nil {
		return errors.Wrap(err, "create new sheets services")
	}

	// TODO: get new spreadsheetid
	spreadsheetId := "1WX6vvcCJihBJ9tFN-8AixYAyt5i1nSfMeX81gsEEwjs"
	writeRange := "A:AZ"
	valRange := &sheets.ValueRange{}
	valRange.Values = append(valRange.Values, []interface{}{
		// TODO: write the data
	})

	if _, err := service.Spreadsheets.Values.Append(spreadsheetId, writeRange, valRange).ValueInputOption("RAW").Do(); err != nil {
		return errors.Wrap(err, "append data to spreadsheet")
	}

	return nil
}
