package mtgsecretlair

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"
	"time"

	dynamicconfiguration "github.com/grep-michael/WebBot/DynamicConfiguration"
	"github.com/grep-michael/WebBot/globals"
)

type MTGSecretLairBot struct {
	InstanceName       string
	PollingInterval    time.Duration
	CollectionInterval time.Duration

	NotificationsDestinations []globals.NotificationDestination
	Cache                     globals.NotificationCache

	opts   dynamicconfiguration.MTGSecretLairBotOptions
	client *http.Client
	ctx    context.Context
}

func (bot *MTGSecretLairBot) Name() string {
	return bot.InstanceName
}

func (bot *MTGSecretLairBot) Run(ctx context.Context) error {
	bot.ctx = ctx
	log.Printf("[%s] Searching every %v\n", bot.InstanceName, bot.PollingInterval)

	products, _ := bot.queryAllProducts()
	notifications := bot.convertProductsToNotifications(products)
	bot.handleNotifications(notifications, bot.opts.NotifyInitial)

	return bot.queryLoop()
}
func (bot *MTGSecretLairBot) queryLoop() error {
	ticker := time.NewTicker(bot.PollingInterval)
	for {
		select {
		case <-bot.ctx.Done():
			return bot.ctx.Err()
		case <-ticker.C:
			err := bot.queryAndNotify()
			if err != nil {
				return err
			}
		}
	}
}
func (bot *MTGSecretLairBot) queryAndNotify() error {
	products, err := bot.queryAllProducts()
	if err != nil {
		return err
	}
	notifications := bot.convertProductsToNotifications(products)
	err = bot.handleNotifications(notifications, true)
	return err
}
func (bot *MTGSecretLairBot) handleNotifications(notifications []globals.Notification, notify bool) error {
	for _, notif := range notifications {
		err := bot.Cache.Cache(notif)
		if err != nil {
			continue
		}
		if notify {
			err := globals.SendNotification(bot.ctx, notif, bot.NotificationsDestinations)
			if err != nil {
				return err
			}
		}
	}
	return nil
}
func (bot *MTGSecretLairBot) convertProductsToNotifications(products []Product) []globals.Notification {
	var notifications []globals.Notification
	for _, product := range products {
		notif := productToNotification(product)
		notif.MetaData.BotName = bot.Name()
		notif.Name = bot.Name()
		notifications = append(notifications, notif)
	}
	return notifications
}
func (bot *MTGSecretLairBot) queryAllProducts() ([]Product, error) {
	offset := 0
	total := 99999
	var products []Product
	for offset < total {
		result, err := bot.queryOffset(offset)
		if err != nil {
			return products, fmt.Errorf("Failed to query at offset %d\n: %+v", offset, err)
		}
		total = result.Total
		offset += result.Count
		products = append(products, result.Products...)
		if offset < total {
			select {
			case <-bot.ctx.Done():
				return products, bot.ctx.Err()
			case <-time.After(bot.CollectionInterval):
			}
		}
	}
	log.Printf("[%s] %d products found\n", bot.InstanceName, len(products))
	return products, nil
}
func (bot *MTGSecretLairBot) queryOffset(offset int) (*Filter, error) {
	req, err := bot.buildRequest(offset)
	if err != nil {
		return nil, err
	}
	resp, err := bot.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status %d: %s", resp.StatusCode, string(body))
	}
	var result SearchResponse
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("parse response: %w", err)
	}
	if len(result.Filters) != 1 {
		return nil, fmt.Errorf("Filters not equal to 1: %s", string(body))
	}
	return &result.Filters[0], nil

}
func (bot *MTGSecretLairBot) buildRequest(offset int) (*http.Request, error) {
	req, err := http.NewRequestWithContext(bot.ctx, http.MethodGet, apiUrl, nil)
	if err != nil {
		return nil, err
	}
	query := req.URL.Query()
	query.Add("userID", userID)
	query.Add("locale", "en_US")
	query.Add("currency", "USD")
	query.Add("crit", "ALL")
	query.Add("sort", "availability_wotc,featured")
	query.Add("offset", strconv.Itoa(offset))
	query.Add("count", strconv.Itoa(bot.opts.SearchCount))
	query.Add("env", "prod")
	query.Add("preference", preference)
	query.Add("filters", filters)
	req.URL.RawQuery = query.Encode()

	req.Header.Set("User-Agent", "Mozilla/5.0 (X11; Linux x86_64; rv:148.0) Gecko/20100101 Firefox/148.0")
	req.Header.Set("Accept", "application/json, text/plain, */*")
	req.Header.Set("Accept-Language", "en-US,en;q=0.9")
	req.Header.Set("Accept-Encoding", "gzip, deflate, br, zstd")
	req.Header.Set("Referer", "https://secretlair.wizards.com/")
	req.Header.Set("Origin", "https://secretlair.wizards.com")
	req.Header.Set("Connection", "keep-alive")
	req.Header.Set("Sec-Fetch-Dest", "empty")
	req.Header.Set("Sec-Fetch-Mode", "cors")
	req.Header.Set("Sec-Fetch-Site", "cross-site")

	return req, nil
}

func init() {
	dynamicconfiguration.RegisterBotType("MTGSecretLair",
		func(bc dynamicconfiguration.BotConfig, bd *dynamicconfiguration.BotDependencies) (globals.Bot, error) {

			var opts dynamicconfiguration.MTGSecretLairBotOptions
			err := json.Unmarshal(bc.Options, &opts)
			if err != nil {
				return nil, err
			}
			pollingInterval, err := time.ParseDuration(opts.PollingInterval)
			if err != nil {
				return nil, err
			}
			colInterval, err := time.ParseDuration(opts.CollectionInterval)
			if err != nil {
				return nil, err
			}

			return &MTGSecretLairBot{
				InstanceName:              bc.InstanceName,
				PollingInterval:           pollingInterval,
				CollectionInterval:        colInterval,
				opts:                      opts,
				Cache:                     bd.Cache,
				NotificationsDestinations: bd.Notifications,
				client:                    &http.Client{Timeout: 10 * time.Second},
			}, nil
		})
}
