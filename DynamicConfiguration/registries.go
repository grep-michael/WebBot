package dynamicconfiguration

import (
	"fmt"

	"github.com/grep-michael/WebBot/globals"
)

// bots
type BotFactory func(BotConfig, *BotDependencies) (globals.Bot, error)

var botRegistry = map[string]BotFactory{}

func RegisterBotType(typ string, factory BotFactory) {
	botRegistry[typ] = factory
}
func CreateBot(cfg BotConfig) (globals.Bot, error) {
	factory, ok := botRegistry[cfg.Type]
	if !ok {
		return nil, fmt.Errorf("%s bot type not found", cfg.Type)
	}
	deps, err := resolveBotDependencies(cfg)
	if err != nil {
		return nil, err
	}
	return factory(cfg, deps)
}

func resolveBotDependencies(cfg BotConfig) (*BotDependencies, error) {
	deps := &BotDependencies{
		Notifications: make([]globals.NotificationDestination, 0),
	}
	cache, err := CreateCache(cfg.Cache)
	if err != nil {
		return deps, err
	}
	deps.Cache = cache
	for _, notificationCfg := range cfg.Notifications {
		notificationDest, err := CreateNotification(notificationCfg)
		if err != nil {
			return deps, err
		}
		deps.Notifications = append(deps.Notifications, notificationDest)
	}
	return deps, nil
}

// Caches
type CacheFactory func(CacheConfig) (globals.NotificationCache, error)

var cacheRegistry = map[string]CacheFactory{}

func RegisterCacheType(typ string, factory CacheFactory) {
	cacheRegistry[typ] = factory
}
func CreateCache(cfg CacheConfig) (globals.NotificationCache, error) {
	factory, ok := cacheRegistry[cfg.Type]
	if !ok {
		return nil, fmt.Errorf("%s cache type not found", cfg.Type)
	}
	return factory(cfg)
}

// notifications
type NotificationFactory func(NotificationConfig) (globals.NotificationDestination, error)

var notifyRegistry = map[string]NotificationFactory{}

func RegisterNotificationType(typ string, factory NotificationFactory) {
	notifyRegistry[typ] = factory
}
func CreateNotification(cfg NotificationConfig) (globals.NotificationDestination, error) {
	factory, ok := notifyRegistry[cfg.Type]
	if !ok {
		return nil, fmt.Errorf("%s Notification type not found", cfg.Type)
	}
	return factory(cfg)
}
