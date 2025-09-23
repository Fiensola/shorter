package dto

type ShorterRequest struct {
	OriginalUrl string `json:"original_url" validate:"required,url"`
	CustomAlias string `json:"custom_alias,omitempty" validate:"omitempty,alphanum,min=3,max=100"`
	ExpiresIn   *int   `json:"expires_in,omitempty"`
}
