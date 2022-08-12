package routes

import (
	"github.com/Verse1/url-shortener/api/db"
	"github.com/go-redis/redis/v8"
	"github.com/gofiber/fiber/v2"
)

func connect(c *fiber.Ctx) error {

	url := c.Params("url")
	rdb := db.DBinit(0)
	defer rdb.Close()

	// rdb.get(db.Ctx, url)
	val, err := rdb.Get(db.Ctxt, url).Result()

	
	if err != nil || err==redis.Nil {
		return c.Status(404).JSON(fiber.Map{
			"error": "Not found",
		})
	}

	increment:=db.DBinit(1)
	defer increment.Close()
	_=increment.Incr(db.Ctxt, "count")

	return c.Redirect(val, 301)
}