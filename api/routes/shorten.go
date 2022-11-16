package routes

import (
	"github.com/asaskevich/govalidator"
	"github.com/go-redis/redis/v8"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/sarkartanmay393/URL-Shortener-Go/database"
	"github.com/sarkartanmay393/URL-Shortener-Go/helpers"
	"os"
	"strconv"
	"time"
)

// Custom Request Structure
type request struct {
	URL         string        `json:"url"`
	CustomShort string        `json:"short"`
	Expiry      time.Duration `json:"expiry"`
}

// Custom Response Structure
type response struct {
	URL             string        `json:"url"`
	CustomShort     string        `json:"short"`
	Expiry          time.Duration `json:"expiry"`
	XRateRemaining  int           `json:"x-rate-remaining"`
	XRateLimitReset time.Duration `json:"x-rate-limit-reset"`
}

func ShortenURL(ctx *fiber.Ctx) error {
	body := new(request)
	if err := ctx.BodyParser(&body); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "cannot request parse body"})
	}

	// Implement rate limiting
	r1 := database.CreateClient(1)
	defer r1.Close()

	val, err := r1.Get(database.Ctx, ctx.IP()).Result()
	if err == redis.Nil {
		_ = r1.Set(database.Ctx, ctx.IP(), os.Getenv("API_QUOTA"), 30*60*time.Second).Err()
	} else {
		valInt, _ := strconv.Atoi(val)
		if valInt <= 0 {
			limit, _ := r1.TTL(database.Ctx, ctx.IP()).Result()
			return ctx.Status(fiber.StatusServiceUnavailable).JSON(fiber.Map{
				"error":      "Rate limit crossed",
				"rate_limit": limit / time.Nanosecond / time.Minute,
			})
		}
	}

	// Check if the URL is valid
	if !govalidator.IsURL(body.URL) {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid URL"})
	}

	// Check for domain error, It is kind of error when the domain is not valid
	if !helpers.RemoveDomainError(body.URL) {
		return ctx.Status(fiber.StatusServiceUnavailable).JSON(fiber.Map{"error": "cannot remove domain error"})
	}

	// Enforce https, ssl, etc.
	body.URL = helpers.EnforceHTTP(body.URL)

	var id string
	if body.CustomShort == "" {
		// uuid is a Google package for random id.
		id = uuid.New().String()[:6]
	} else {
		id = body.CustomShort
	}

	r := database.CreateClient(0)
	defer r.Close()

	// duplication of shorten url is prohibited here.
	val, _ = r.Get(database.Ctx, id).Result()
	if val != "" {
		return ctx.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"error": "custom short is already in use",
		})
	}

	if body.Expiry == 0 {
		body.Expiry = 24
	}

	err = r.Set(database.Ctx, id, body.URL, body.Expiry*3600*time.Second).Err()
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":    "unable to connect to server",
			"rawError": err.Error(),
		})
	}

	resp := response{
		URL:             body.URL,
		CustomShort:     "",
		Expiry:          body.Expiry,
		XRateRemaining:  100,
		XRateLimitReset: 30,
	}

	// Decrementing api quota for user
	r1.Decr(database.Ctx, ctx.IP())

	val, _ = r1.Get(database.Ctx, ctx.IP()).Result()
	resp.XRateRemaining, _ = strconv.Atoi(val)

	ttl, _ := r1.TTL(database.Ctx, ctx.IP()).Result()
	resp.XRateLimitReset = ttl / time.Nanosecond / time.Minute

	resp.CustomShort = os.Getenv("DOMAIN") + "/" + id

	// check if the custom short is valid
	if !govalidator.IsURL(resp.CustomShort) {
		//return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid custom short URL"})
	}

	return ctx.Status(fiber.StatusOK).JSON(resp)
}
