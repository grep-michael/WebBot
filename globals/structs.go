package globals

import "time"

type Notification struct {
	Id       string
	Name     string
	Source   string
	Message  string
	ImageUrl string
	MetaData NotificationMetaData
}

type NotificationMetaData struct {
	Timestamp time.Time
	BotName   string
	Tags      map[string]string
}
