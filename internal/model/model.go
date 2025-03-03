package model

type LongURL struct {
	URL string `json:"long_url" binding:"required"`
}

type ShortURL struct {
	URL string `json:"short_url" binding:"required"`
}
