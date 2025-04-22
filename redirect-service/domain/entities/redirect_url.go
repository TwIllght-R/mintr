package entities

import "time"

type RedirectURL struct {
	ShortCode   string
	OriginalURL string
	IP          string
	UserAgent   string
	VisitedAt   time.Time
}
