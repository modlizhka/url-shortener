package model

import "errors"

var (
	ErrEmtpyURL = errors.New("url is empty")
)

type LongURL struct {
	URL string `json:"long_url"`
}

type ShortURL struct {
	URL string `json:"short_url"`
}

func (u LongURL) Validate() error {
	if u.URL == "" {
		return ErrEmtpyURL
	}
	return nil
}

func (u ShortURL) Validate() error {
	if u.URL == "" {
		return ErrEmtpyURL
	}
	return nil
}
