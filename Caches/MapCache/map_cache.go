package mapcache

import (
	"fmt"

	dynamicconfiguration "github.com/grep-michael/WebBot/DynamicConfiguration"
	"github.com/grep-michael/WebBot/globals"
)

type MapCache struct {
	cacheMap map[string]bool
}

func NewMapCache() *MapCache {
	return &MapCache{
		cacheMap: make(map[string]bool),
	}
}
func (cache *MapCache) Cache(notification globals.Notification) error {
	_, ok := cache.cacheMap[notification.Id]
	if ok {
		return fmt.Errorf("Notifcation already in cache")
	}
	cache.cacheMap[notification.Id] = true
	return nil
}

func (cache *MapCache) Reset() error {
	cache.cacheMap = make(map[string]bool)
	return nil
}

func init() {
	dynamicconfiguration.RegisterCacheType("MapCache", func(cc dynamicconfiguration.CacheConfig) (globals.NotificationCache, error) {
		return NewMapCache(), nil
	})
}
