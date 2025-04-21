package entities

import "time"

type URLMapping struct {
	UUID        string
	Title       string
	ShortCode   string
	OriginalURL string
	CreatedAt   time.Time
	UpdatedAt   time.Time
}
