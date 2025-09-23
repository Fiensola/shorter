package dto

type ShorterResponse struct {
	ShortUrl  string `json:"short_url"`
	ExpiresAt string `json:"expires_at,omitempty"`
}
