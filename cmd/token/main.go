package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/gmail/v1"
	"gopkg.in/ini.v1"
)

func main() {
	// generate a new token
	if err := os.WriteFile("sheets_token.json", []byte(""), 0644); err != nil {
		log.Fatalf("overwrite token file: %v", err)
	}
	if err := os.WriteFile("gmail_token.json", []byte(""), 0644); err != nil {
		log.Fatalf("overwrite token file: %v", err)
	}

	b, err := os.ReadFile("credentials.json")
	if err != nil {
		log.Fatalf("read client secret file: %v", err)
	}

	{
		config, err := google.ConfigFromJSON(b, "https://www.googleapis.com/auth/spreadsheets")
		if err != nil {
			log.Fatalf("unable to parse client secret file to config: %v", err)
		}
		getClient(config, "sheets_token.json")
		log.Println("=> generated new sheets token")
	}

	{
		config, err := google.ConfigFromJSON(b, gmail.GmailComposeScope)
		if err != nil {
			log.Fatalf("unable to parse client secret file to config: %v", err)
		}
		getClient(config, "gmail_token.json")
		log.Println("=> generated new gmail token")
	}

	credentialsData, err := os.ReadFile("credentials.json")
	if err != nil {
		log.Fatalf("read credentials: %v", err)
	}

	gmailData, err := os.ReadFile("gmail_token.json")
	if err != nil {
		log.Fatalf("read gmail token: %v", err)
	}

	sheetsData, err := os.ReadFile("sheets_token.json")
	if err != nil {
		log.Fatalf("read sheets token: %v", err)
	}

	data, err := ini.Load("data.env")
	if err != nil {
		log.Fatalf("read env file: %v", err)
	}

	// update / overwrite env variables
	data.Section("").Key("CREDENTIALS").SetValue(string(credentialsData))
	data.Section("").Key("GMAIL").SetValue(string(gmailData))
	data.Section("").Key("SHEETS").SetValue(string(sheetsData))

	if err := data.SaveTo("data.env"); err != nil {
		log.Fatalf("save env file: %v", err)
	}
}

// Retrieve a token, saves the token, then returns the generated client.
func getClient(config *oauth2.Config, path string) *http.Client {
	// The file token.json stores the user's access and refresh tokens, and is
	// created automatically when the authorization flow completes for the first
	// time.
	tok, err := tokenFromFile(path)
	if err != nil {
		tok = getTokenFromWeb(config)
		saveToken(path, tok)
	}
	return config.Client(context.Background(), tok)
}

// Request a token from the web, then returns the retrieved token.
func getTokenFromWeb(config *oauth2.Config) *oauth2.Token {
	authURL := config.AuthCodeURL("state-token", oauth2.AccessTypeOffline)
	fmt.Printf("Go to the following link in your browser then type the "+
		"authorization code: \n%v\n", authURL)

	var authCode string
	if _, err := fmt.Scan(&authCode); err != nil {
		log.Fatalf("Unable to read authorization code: %v", err)
	}

	tok, err := config.Exchange(context.TODO(), authCode)
	if err != nil {
		log.Fatalf("Unable to retrieve token from web: %v", err)
	}
	return tok
}

// Retrieves a token from a local file.
func tokenFromFile(file string) (*oauth2.Token, error) {
	f, err := os.Open(file)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	tok := &oauth2.Token{}
	err = json.NewDecoder(f).Decode(tok)
	return tok, err
}

// Saves a token to a file path.
func saveToken(path string, token *oauth2.Token) {
	fmt.Printf("Saving credential file to: %s\n", path)
	f, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		log.Fatalf("Unable to cache oauth token: %v", err)
	}
	defer f.Close()
	json.NewEncoder(f).Encode(token)
}
