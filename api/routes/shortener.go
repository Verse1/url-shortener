package routes

import (
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/Verse1/url-shortener/api/db"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/internal/uuid"
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

	val,err :=rdb.Get(db.Ctxt, c.IP()).Result()

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


	var ID string

	if body.Custom != "" {
		ID = uuid.NewV4().String()[0:6]
	} else {
		ID = body.Custom
	}

	rdb2 := db.DBinit(0)
	defer rdb2.Close()

	val, _= rdb2.Get(db.Ctxt, ID).Result()

	if val != "" {
		return c.Status(400).JSON(fiber.Map{
			"error": "Shortened URL already exists",
		})
	}

	if body.Expire==0 {
		body.Expire = 24*time.Hour
	}

	err=rdb2.Set(db.Ctxt, ID, body.URL, body.Expire).Err()

	if err != nil {

		return c.Status(500).JSON(fiber.Map{
			"error": "Internal server error",
		})
	}

	rdb.Decr(db.Ctxt, c.IP())

}