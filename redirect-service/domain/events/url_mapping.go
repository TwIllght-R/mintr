package events

import "time"

type URLCreatedEvent struct {
	UUID        string    `json:"uuid"`
	ShortCode   string    `json:"short_code"`
	OriginalURL string    `json:"original_url"`
	ExpiresAt   time.Time `json:"expires_at"`
	UTMSource   string    `json:"utm_source,omitempty"`
	UTMMedium   string    `json:"utm_medium,omitempty"`
	UTMCampaign string    `json:"utm_campaign,omitempty"`
	UTMTerm     string    `json:"utm_term,omitempty"`
	UTMContent  string    `json:"utm_content,omitempty"`
}

type URLUpdatedEvent struct {
	UUID        string    `json:"uuid"`
	ShortCode   string    `json:"short_code"`
	OriginalURL string    `json:"original_url"`
	ExpiresAt   time.Time `json:"expires_at"`
	UTMSource   string    `json:"utm_source,omitempty"`
	UTMMedium   string    `json:"utm_medium,omitempty"`
	UTMCampaign string    `json:"utm_campaign,omitempty"`
	UTMTerm     string    `json:"utm_term,omitempty"`
	UTMContent  string    `json:"utm_content,omitempty"`
}

type URLDeletedEvent struct {
	UUID string `json:"uuid"`
}
