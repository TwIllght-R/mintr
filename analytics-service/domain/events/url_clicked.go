package events

import "time"

type URLClickedEvent struct {
	UrlUUID   string
	ShortCode string
	IP        string
	UserAgent string
	Referer   string
	Timestamp time.Time
	OwnerUUID string
}
