package nopcache

import (
	dynamicconfiguration "github.com/grep-michael/WebBot/DynamicConfiguration"
	"github.com/grep-michael/WebBot/globals"
)

type NopCache struct {
	cacheMap map[string]bool
}

func (cache *NopCache) Cache(notification globals.Notification) error {
	return nil
}

func (cache *NopCache) Reset() error {
	return nil
}

func init() {
	dynamicconfiguration.RegisterCacheType("NopCache", func(cc dynamicconfiguration.CacheConfig) (globals.NotificationCache, error) {
		return &NopCache{}, nil
	})
}
