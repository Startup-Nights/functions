package lib

import (
	"github.com/pkg/errors"
	"github.com/slack-go/slack"
)

func SendToSlack(webhook, channel, content string) error {
	if err := slack.PostWebhook(webhook, &slack.WebhookMessage{
		Channel: channel,
		Text:    content,
	}); err != nil {
		return errors.Wrap(err, "send message to channel")
	}

	return nil
}
