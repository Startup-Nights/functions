package main

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"

	"github.com/pkg/errors"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/gmail/v1"
	"google.golang.org/api/option"
)

func sendMail(ctx context.Context, in Request) error {
	config, err := google.ConfigFromJSON([]byte(credentials), gmail.GmailComposeScope)
	if err != nil {
		return errors.Wrap(err, "get gmail config from json")
	}

	token := &oauth2.Token{}
	if err := json.NewDecoder(bytes.NewBuffer([]byte(gmailCreds))).Decode(token); err != nil {
		return errors.Wrap(err, "decode gmail credentials")
	}

	client := config.Client(ctx, token)
	service, err := gmail.NewService(ctx, option.WithHTTPClient(client))
	if err != nil {
		return errors.Wrap(err, "create new gmail services")
	}

	var message gmail.Message

	buf := bytes.NewBuffer([]byte{})
	if err := tpl.Execute(buf, in); err != nil {
		return errors.Wrap(err, "execute template")
	}

	msg := fmt.Sprintf(
		"To: %s\r\nSubject: %s\n%s\n\n\n%s",
		in.Contact.Email,
		"Startup Nights 2023 Booth Application",
		mime,
		buf.String(),
	)

	message.Raw = base64.URLEncoding.EncodeToString([]byte(msg))
	if _, err := service.Users.Messages.Send("me", &message).Do(); err != nil {
		return errors.Wrap(err, "send mail")
	}

	return nil
}

const BoothRegistrationTemplate = `Hi {{.Contact.FirstName}} {{.Contact.LastName}},

thank you for the registration for a booth for {{.Company.Name}} at the Startup Nights 2023 this November. We will reach out to you soon to confirm your registration.
 
{{if .Varia.Formats}}
Beside the booth, you were interested in:
{{range .Varia.Formats}} - {{.}}
{{end}}
We are going to reach out to you regarding the registration for these formats separately.
{{end}}

We know, it's a long time until November so we have some suggestions for you on how you can pass the time ⏱
- Talk to other founders / startups and invite people to the event 🚀
- Buy the tickets for you and your team and don't forget to use your discount code 20OFFspecial to get 20% off 💸 You can buy the tickets here: https://portal.startup-nights.ch
- Read the FAQ if you have open questions: https://startup-nights.ch/faq

Talk to you soon 😉

Best regards,
the Startup Nights Team`
