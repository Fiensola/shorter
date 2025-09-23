package model

import "time"

type Link struct {
	Alias       string     `json:"alias" db:"aliaZ"`
	OriginalUrl string     `json:"original_url"`
	CreatedAt   time.Time  `json:"created_at"`
	ExpiresAt   *time.Time `json:"expires_at"`
	ClickCount  int        `json:"click_count"`
}
