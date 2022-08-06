package routes

import (
	"github.com/Verse1/url-shortener/db"
	"github.com/go-redis/redis/v8"
	"github.com/gofiber/fiber/v2"
)

func connect(c *fiber.Ctx) error {

	url := c.Params("url")
	rdb := db.DBinit(0)
	defer rdb.Close()

	rdb.get(db.Ctxt, url)
	val, err := rdb.Get(db.ctxt, url).Result()

	
	if err != nil || err==redis.Nil {
		return c.Status(404).JSON(fiber.Map{
			"error": "Not found",
		})
	}
	return c.Redirect(val, 301)
}