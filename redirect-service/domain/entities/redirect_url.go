package entities

import "time"

type RedirectURL struct {
	UUID        string
	ShortCode   string
	OriginalURL string
	ExpiresAt   time.Time
	UTMSource   string
	UTMMedium   string
	UTMCampaign string
	UTMTerm     string
	UTMContent  string
	CreatedAt   time.Time
	UpdatedAt   time.Time
}
