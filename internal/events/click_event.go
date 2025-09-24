package events

import "time"

type ClickEvent struct {
	Alias     string    `json:"alias"`
	Timestamp time.Time `json:"timestamp"`
	IP        string    `json:"ip"`
	UserAgent string    `json:"user_agent"`
	Referer   string    `json:"referer,omitempty"`
}
