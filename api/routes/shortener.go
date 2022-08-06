package routes

import (
	"time"
	"strings"
)

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

func HTTP(url string) string {
	if !strings.HasPrefix(url, "http://") {
		url = "http://" + url
	}
	return url
}


func shorten(c *fiber.Ctx) error {

	body := new(req)
	if err := c.BodyParser(&body); err != nil {
		return err
	}

	if !govalidator.IsURL(body.URL) {
		return c.Status(400).JSON(res{
			URL: body.URL,
			Shortened: "",
			Expiration: 0,
			RateLimit: 0,
			RateReset: 0,
		})
	}

	body.URL = HTTPS(body.URL)

}