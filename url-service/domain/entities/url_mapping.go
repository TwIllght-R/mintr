package entities

import "time"

type URLMapping struct {
	UUID        string
	OwnerUUID   string
	ShortCode   string
	OriginalURL string
	Title       string
	UTMSource   string
	UTMMedium   string
	UTMCampaign string
	UTMTerm     string
	UTMContent  string
	ExpiresAt   time.Time
	CreatedAt   time.Time
	UpdatedAt   time.Time
}
