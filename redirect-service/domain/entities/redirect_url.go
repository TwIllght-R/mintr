package entities

import "time"

type RedirectURL struct {
	UUID        string    `bson:"uuid"`
	ShortCode   string    `bson:"short_code"`
	OriginalURL string    `bson:"original_url"`
	ExpiresAt   time.Time `bson:"expires_at"`
	UTMSource   string    `bson:"utm_source"`
	UTMMedium   string    `bson:"utm_medium"`
	UTMCampaign string    `bson:"utm_campaign"`
	UTMTerm     string    `bson:"utm_term"`
	UTMContent  string    `bson:"utm_content"`
	CreatedAt   time.Time `bson:"created_at"`
	UpdatedAt   time.Time `bson:"updated_at"`
}
