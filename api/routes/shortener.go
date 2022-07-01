package routes

import "time"

type res struct { 
	URL string `json:"url"`
	Shortened string `json:"shortened"`
	Expiration time.Duration `json:"expiration"`
	RateLimit int `json:"rate_limit"`
	RateReset time.Duration  `json:"rate_reset"`
}

type req struct {
	URL string `json:"url"`
	Shortened string `json:"shortened"`
	Expiration time.Duration `json:"expiration"`
}