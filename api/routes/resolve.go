package routes

import (
	"github.com/go-redis/redis/v8"
	"github.com/gofiber/fiber/v2"
	"github.com/sarkartanmay393/URL-Shortener-Go/database"
)

// ResolveURL resolves the short URL to the original URL
func ResolveURL(ctx *fiber.Ctx) error {
	id := ctx.Params("id")

	r := database.CreateClient(0)
	defer r.Close()

	value, err := r.Get(database.Ctx, id).Result()
	if err == redis.Nil {
		return ctx.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "ID not found in database",
		})
	} else if value == "" {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Error while fetching short url from database",
		})
	}

	//rInr := database.CreateClient(1)
	//defer rInr.Close()
	//
	//_ = rInr.Incr(database.Ctx, url)

	return ctx.Redirect(value, fiber.StatusTemporaryRedirect)
}
