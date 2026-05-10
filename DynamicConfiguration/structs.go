package dynamicconfiguration

import (
	"encoding/json"

	"github.com/grep-michael/WebBot/globals"
)

type BotList struct {
	Bots []BotConfig `json:"Bots"`
}

type BotConfig struct {
	Type          string               `json:"Type"`
	InstanceName  string               `json:"InstanceName"`
	Notifications []NotificationConfig `json:"Notifications"`
	Cache         CacheConfig          `json:"Cache"`
	Options       json.RawMessage      `json:"Options"`
}

type BotDependencies struct {
	Notifications []globals.NotificationDestination
	Cache         globals.NotificationCache
}

type NotificationConfig struct {
	Type    string          `json:"Type"`
	Options json.RawMessage `json:"Options"`
}
type CacheConfig struct {
	Type    string          `json:"Type"`
	Options json.RawMessage `json:"Options"`
}
