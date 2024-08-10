package models

import "errors"

type CreateRequest struct {
	LongURL     string `json:"long_url"`
	CustomAlias string `json:"custom_alias"`
	TTLSeconds  uint   `json:"ttl_seconds"`
}

func (r CreateRequest) Validate() error {
	if r.LongURL == "" {
		return errors.New("missing long_url")
	}

	return nil
}

type CreateResponse struct {
	ShortURL string `json:"short_url"`
}

type Analytics struct {
	Alias       string   `json:"alias"`
	LongURL     string   `json:"long_url"`
	TTLSeconds  uint     `json:"ttl_seconds"`
	AccessCount uint     `json:"access_count"`
	AccessTimes []string `json:"access_times"`
	CreatedAt   string   `json:"-"`
}

type UpdateRequest struct {
	CustomAlias string `json:"custom_alias"`
	TTLSeconds  uint   `json:"ttl_seconds"`
}

func (r UpdateRequest) Validate() error {
	if r.CustomAlias == "" && r.TTLSeconds == 0 {
		return errors.New("invalid request")
	}

	return nil
}
