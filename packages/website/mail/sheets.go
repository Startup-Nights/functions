package main

import (
	"bytes"
	"context"
	"encoding/json"
	"strings"

	"github.com/pkg/errors"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/option"
	"google.golang.org/api/sheets/v4"
)

func saveToSheets(ctx context.Context, in Request) error {
	config, err := google.ConfigFromJSON([]byte(credentials), "https://www.googleapis.com/auth/spreadsheets")
	if err != nil {
		return errors.Wrap(err, "get gmail config from json")
	}

	token := &oauth2.Token{}
	if err := json.NewDecoder(bytes.NewBuffer([]byte(sheetsCreds))).Decode(token); err != nil {
		return errors.Wrap(err, "decode gmail credentials")
	}

	client := config.Client(ctx, token)
	service, err := sheets.NewService(ctx, option.WithHTTPClient(client), option.WithScopes(sheets.SpreadsheetsScope))
	if err != nil {
		return errors.Wrap(err, "create new gmail services")
	}

	spreadsheetId := "1WX6vvcCJihBJ9tFN-8AixYAyt5i1nSfMeX81gsEEwjs"
	writeRange := "A:AZ"
	valRange := &sheets.ValueRange{}
	valRange.Values = append(valRange.Values, []interface{}{
		in.Company.Name,
		in.Company.Website,
		in.Company.FoundingDate,
		strings.Join(in.Company.LinkedIn, "\n"),
		in.Company.Employees,
		in.Company.Pitch,
		strings.Join(in.Company.Categories, "\n"),
		in.Company.AdditionalCategories,
		in.Company.Logo,
		in.Company.Address.Street,
		in.Company.Address.ZIP,
		in.Company.Address.City,
		in.Company.Address.Country,
		in.Company.BillingAddress.Street,
		in.Company.BillingAddress.ZIP,
		in.Company.BillingAddress.City,
		in.Company.BillingAddress.Country,
		in.Contact.FirstName,
		in.Contact.LastName,
		in.Contact.Role,
		in.Contact.Email,
		in.Contact.Phone,
		in.Varia.Package.Title,
		strings.Join(in.Varia.Formats, "\n"),
		in.Varia.Accomodation,
		in.Varia.Referral,
		in.Varia.Equipment,
		in.Varia.PreviousVisitor,
	})

	if _, err := service.Spreadsheets.Values.Append(spreadsheetId, writeRange, valRange).ValueInputOption("RAW").Do(); err != nil {
		return errors.Wrap(err, "append data to spreadsheet")
	}

	return nil
}
