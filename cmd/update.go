package cmd

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sync"

	"github.com/go-chi/chi"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/gmail/v1"
)

var (
	renewGmailToken  bool
	renewSheetsToken bool
	config           *oauth2.Config

	updateCmd = &cobra.Command{
		Use:   "update",
		Short: "Renew the tokens and update them in the config file",
		Run: func(cmd *cobra.Command, args []string) {
			credentialsKey := "credentials"

			wg := &sync.WaitGroup{}

			if !viper.IsSet(credentialsKey) {
				cobra.CheckErr(errors.New("no sheets token file configured"))
			}

			b := []byte(viper.GetString(credentialsKey))

			r := chi.NewRouter()
			r.Get("/", func(w http.ResponseWriter, r *http.Request) {
				defer wg.Done()

				scope := r.URL.Query().Get("scope")
				code := r.URL.Query().Get("code")

				token, err := config.Exchange(context.TODO(), code)
				if err != nil {
					_, _ = w.Write([]byte("failed to get token from code: " + err.Error()))
				}

				var buf bytes.Buffer
				if err := json.NewEncoder(&buf).Encode(token); err != nil {
					_, _ = w.Write([]byte("failed to encode token to json: " + err.Error()))
				}

				switch scope {
				case "https://www.googleapis.com/auth/gmail.compose":
					viper.Set("gmail_token", buf.String())
					_, _ = w.Write([]byte("updated gmail token"))

				case "https://www.googleapis.com/auth/spreadsheets":
					viper.Set("sheets_token", buf.String())
					_, _ = w.Write([]byte("updated sheets token"))

				default:
					_, _ = w.Write([]byte("unknown scope: " + scope))
				}
			})

			go func() {
				_ = http.ListenAndServe(":3333", r)
			}()

			if renewSheetsToken {
				var err error
				wg.Add(1)
				config, err = google.ConfigFromJSON(b, "https://www.googleapis.com/auth/spreadsheets")
				if err != nil {
					cobra.CheckErr(errors.Wrap(err, "parse sheets client secret"))
				}
				fmt.Println("=> open this link in your browser: " + config.AuthCodeURL("state-token", oauth2.AccessTypeOffline))
				wg.Wait()
				log.Println("=> generated new sheets token")
			}

			if renewGmailToken {
				var err error
				wg.Add(1)
				config, err = google.ConfigFromJSON(b, gmail.GmailComposeScope)
				if err != nil {
					cobra.CheckErr(errors.Wrap(err, "parse gmail client secret"))
				}
				fmt.Println("=> open this link in your browser: " + config.AuthCodeURL("state-token", oauth2.AccessTypeOffline))
				wg.Wait()
				log.Println("=> generated new gmail token")
			}

			if err := viper.WriteConfig(); err != nil {
				cobra.CheckErr(errors.Wrap(err, "update config with new tokens"))
			}

			fmt.Println("=> token can be updated here: " + viper.GetString("secrets_url"))
		},
	}
)

func init() {
	tokenCmd.AddCommand(updateCmd)
	updateCmd.PersistentFlags().BoolVar(&renewGmailToken, "gmail", false, "A help for foo")
	updateCmd.PersistentFlags().BoolVar(&renewSheetsToken, "sheets", false, "A help for foo")
}
