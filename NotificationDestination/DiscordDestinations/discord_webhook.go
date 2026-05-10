package discorddestinations

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	dynamicconfiguration "github.com/grep-michael/WebBot/DynamicConfiguration"
	"github.com/grep-michael/WebBot/globals"
)

type DiscordWebHook struct {
	WebHookUrl string
	token      string
	client     *http.Client
}

func NewDiscordWebhook(webhookURL string) *DiscordWebHook {
	return &DiscordWebHook{
		WebHookUrl: webhookURL,
		client:     &http.Client{Timeout: 10 * time.Second},
	}
}

func (dest *DiscordWebHook) Send(ctx context.Context, notification globals.Notification) error {
	payload := DiscordWebhookPayload{
		Username: notification.MetaData.BotName,
		Embeds:   []DiscordEmbed{dest.buildEmbed(notification)},
	}

	body, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("discord: marshal payload: %w", err)
	}

	result := make(chan error, 1)
	webhookQueue <- webhookJob{
		ctx:      ctx,
		url:      dest.WebHookUrl,
		payload:  body,
		result:   result,
		botToken: dest.token,
	}

	return <-result
}

func (dest *DiscordWebHook) buildEmbed(n globals.Notification) DiscordEmbed {
	fields := []DiscordEmbedField{}
	for k, v := range n.MetaData.Tags {
		fields = append(fields, DiscordEmbedField{
			Name:  k,
			Value: v,
		})
	}
	embed := DiscordEmbed{
		Title:       n.Name,
		Description: n.Message,
		Color:       0x5865F2,
		Fields:      fields,
		URL:         n.Source,
		Thumbnail: DiscordEmbedMedia{
			URL: n.ImageUrl,
		},
		Image: DiscordEmbedMedia{
			URL: n.ImageUrl,
		},
		Footer: DiscordEmbedFooter{
			Text: fmt.Sprintf("Generated At %s", n.MetaData.Timestamp.Format(time.RFC822)),
		},
	}
	return embed
}

func init() {
	dynamicconfiguration.RegisterNotificationType("DiscordWebhook", func(nc dynamicconfiguration.NotificationConfig) (globals.NotificationDestination, error) {
		var opts dynamicconfiguration.DiscordWebhookOptions
		err := json.Unmarshal(nc.Options, &opts)
		if err != nil {
			return nil, err
		}
		if strings.TrimSpace(opts.WebhookUrl) == "" {
			return nil, fmt.Errorf("Webhook Url can not be empty")
		}

		return &DiscordWebHook{
			WebHookUrl: opts.WebhookUrl,
			token:      opts.BotToken,
			client:     &http.Client{Timeout: 10 * time.Second},
		}, nil
	})
}
