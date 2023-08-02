package lib

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

const (
	mime = "MIME-version: 1.0;\nContent-Type: text/plain; charset=\"UTF-8\";"
)

type MailService struct {
	service *gmail.Service
}

func NewMailService(ctx context.Context, creds Credentials) (*MailService, error) {
	svc := &MailService{}

	config, err := google.ConfigFromJSON([]byte(creds.Config), gmail.GmailComposeScope)
	if err != nil {
		return svc, errors.Wrap(err, "get gmail config from json")
	}

	token := &oauth2.Token{}
	if err := json.NewDecoder(bytes.NewBuffer([]byte(creds.Credentials))).Decode(token); err != nil {
		return svc, errors.Wrap(err, "decode gmail credentials")
	}

	client := config.Client(ctx, token)
	service, err := gmail.NewService(ctx, option.WithHTTPClient(client))
	if err != nil {
		return svc, errors.Wrap(err, "create new gmail services")
	}

	svc.service = service

	return svc, nil
}

func (m *MailService) Send(ctx context.Context, subject, title, content string) error {
	var message gmail.Message

	msg := fmt.Sprintf(
		"To: %s\r\nSubject: %s\n%s\n\n\n%s",
		subject,
		title,
		mime,
		content,
	)

	message.Raw = base64.URLEncoding.EncodeToString([]byte(msg))
	if _, err := m.service.Users.Messages.Send("me", &message).Do(); err != nil {
		return errors.Wrap(err, "send mail")
	}

	return nil
}
