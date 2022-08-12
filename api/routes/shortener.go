package routes

import (
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/Verse1/url-shortener/api/db"
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


	rdb := db.DBinit(1)
	defer rdb.Close()

	val,err :=rdb.Set(db.Ctxt, c.IP()).Result()

	if err == nil {
		_=rdb.Set(db.Ctxt, c.IP(), os.Getenv("QUOTA"), 60*60*time.Second).Err()
	} else{
		val, _= rdb.Get(db.Ctxt, c.IP()).Result()
		val2, _:= strconv.Atoi(val)

		if val2 <= 0 {
			rdb.TTL(db.Ctxt, c.IP()).Result()
			return c.Status(429).JSON(fiber.Map{
				"error": "Rate limit exceeded",
			})
		}
	}


	if !govalidator.IsURL(body.URL) {
		return c.Status(400).JSON(fiber.Map{
			"error": "URL is not valid",
			})
	}

	body.URL = HTTP(body.URL)

	rdb.Decr(db.Ctxt, c.IP())

}