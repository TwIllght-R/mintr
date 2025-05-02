package entities

import "time"

type URLClicked struct {
	UrlUUID   string
	OwnerUUID string
	ShortCode string
	IP        string
	UserAgent string
	Platform  string
	Referer   string
	Country   string
	City      string
	Timestamp time.Time
}
