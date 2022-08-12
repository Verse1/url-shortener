package routes

import (
	"os"
	"strconv"
	"strings"
	"time"
	"github.com/Verse1/url-shortener/db"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/go-redis/redis/v8"
	"github.com/asaskevich/govalidator"
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


func Shorten(c *fiber.Ctx) error {

	body := new(req)
	if err := c.BodyParser(&body); err != nil {
		return err
	}


	rdb := db.DBinit(1)
	defer rdb.Close()

	val,err :=rdb.Get(db.Ctxt, c.IP()).Result()

	if err == redis.Nil {
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

	if body.Shortened != "" {
		ID = uuid.New().String()[0:6]
	} else {
		ID = body.Shortened
	}

	rdb2 := db.DBinit(0)
	defer rdb2.Close()

	val, _= rdb2.Get(db.Ctxt, ID).Result()

	if val != "" {
		return c.Status(400).JSON(fiber.Map{
			"error": "Shortened URL already exists",
		})
	}

	if body.Expiration==0 {
		body.Expiration = 24*time.Hour
	}

	err=rdb2.Set(db.Ctxt, ID, body.URL, body.Expiration).Err()

	if err != nil {

		return c.Status(500).JSON(fiber.Map{
			"error": "Internal server error",
		})
	}

	response:=res{
		URL: body.URL,
		Shortened: "",
		Expiration: body.Expiration,
		RateLimit: 20,
		RateReset: 60*60*time.Second,
	}


	rdb.Decr(db.Ctxt, c.IP())

	val,_=rdb.Get(db.Ctxt, c.IP()).Result()

	response.RateLimit, _= strconv.Atoi(val)
	ttl, _:= rdb.TTL(db.Ctxt, c.IP()).Result()
	response.RateReset = ttl
	response.Shortened = os.Getenv("BASE_URL") + "/" + ID

	return c.Status(200).JSON(response)
}