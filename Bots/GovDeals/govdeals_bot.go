package govdeals

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	dynamicconfiguration "github.com/grep-michael/WebBot/DynamicConfiguration"
	"github.com/grep-michael/WebBot/globals"
)

type GovDealsBot struct {
	InstanceName    string
	PollingInterval time.Duration
	options         dynamicconfiguration.GovDealsBotOptions

	NotificationsDestinations []globals.NotificationDestination
	Cache                     globals.NotificationCache

	ticker *time.Ticker
	http   *http.Client
}

func (bot *GovDealsBot) Name() string {
	return bot.InstanceName
}
func (bot *GovDealsBot) Run(ctx context.Context) error {
	bot.ticker = time.NewTicker(bot.PollingInterval)
	log.Printf("[%s] Searching for \"%s\" every %v\n", bot.InstanceName, bot.options.SearchTerm, bot.PollingInterval)
	if err := bot.checkNewListing(ctx); err != nil {
		return err
	}
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-bot.ticker.C:
			err := bot.checkNewListing(ctx)
			if err != nil {
				return err
			}
		}
	}
}
func (bot *GovDealsBot) checkNewListing(ctx context.Context) error {
	results, err := bot.search(ctx, buildSearchRequest(bot.options))
	if err != nil {
		return err
	}
	log.Printf("[%s] found %d results\n", bot.InstanceName, len(results.Assets))
	for _, asset := range results.Assets {
		notification := assetToNotification(asset)
		if bot.Cache != nil {
			if err := bot.Cache.Cache(notification); err != nil {
				continue
			}
		}
		for _, dest := range bot.NotificationsDestinations {
			err := dest.Send(ctx, notification)
			if err != nil {
				log.Printf("Failed to send Notification:\n\t%v\n", err)
			}
		}
	}
	return nil
}
func (bot *GovDealsBot) search(ctx context.Context, request searchRequest) (*searchResponse, error) {
	body, err := json.Marshal(request)
	if err != nil {
		return nil, err
	}
	httpReq, err := http.NewRequestWithContext(ctx, http.MethodPost, searchURL, bytes.NewReader(body))
	if err != nil {
		return nil, err
	}
	bot.setHeaders(httpReq)
	resp, err := bot.http.Do(httpReq)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, fmt.Errorf("govdeals search: unexpected status %d", resp.StatusCode)
	}

	reader, err := decompressBody(resp)
	if err != nil {
		return nil, err
	}
	defer reader.Close()
	var result searchResponse
	err = json.NewDecoder(reader).Decode(&result)
	if err != nil {
		return nil, err
	}

	return &result, nil
}
func (bot *GovDealsBot) setHeaders(r *http.Request) {
	r.Header.Set("User-Agent", "Mozilla/5.0 (X11; Linux x86_64; rv:148.0) Gecko/20100101 Firefox/148.0")
	r.Header.Set("Accept", "application/json, text/plain, */*")
	r.Header.Set("Accept-Language", "en-US,en;q=0.9")
	r.Header.Set("Accept-Encoding", "gzip, deflate, br, zstd")
	r.Header.Set("Content-Type", "application/json")
	r.Header.Set("Origin", "https://www.govdeals.com")
	r.Header.Set("Referer", "https://www.govdeals.com/")
	r.Header.Set("Connection", "keep-alive")
	r.Header.Set("Sec-Fetch-Dest", "empty")
	r.Header.Set("Sec-Fetch-Mode", "cors")
	r.Header.Set("Sec-Fetch-Site", "cross-site")
	r.Header.Set("Priority", "u=0")
	r.Header.Set("x-api-key", apiKey)
	r.Header.Set("x-api-correlation-id", apiCorrolationId)
	r.Header.Set("Ocp-Apim-Subscription-Key", ocpSubscriptionKey)
	r.Header.Set("x-user-id", "-1")
	r.Header.Set("x-user-timezone", "America/New_York")
}

func init() {
	dynamicconfiguration.RegisterBotType("GovDeals",
		func(bc dynamicconfiguration.BotConfig, bd *dynamicconfiguration.BotDependencies) (globals.Bot, error) {
			var opts dynamicconfiguration.GovDealsBotOptions
			err := json.Unmarshal(bc.Options, &opts)
			if err != nil {
				return nil, err
			}
			interval, err := time.ParseDuration(opts.PollingInterval)
			if err != nil {
				return nil, err
			}
			return &GovDealsBot{
				InstanceName:              bc.InstanceName,
				Cache:                     bd.Cache,
				NotificationsDestinations: bd.Notifications,
				PollingInterval:           interval,
				options:                   opts,
				http:                      &http.Client{Timeout: 30 * time.Second},
			}, nil
		})
}
