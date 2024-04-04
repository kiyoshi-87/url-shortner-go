package routes

import (
	"context"
	"os"
	"strconv"
	"time"

	"github.com/asaskevich/govalidator"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"

	database "github.com/kiyoshi-87/url-shortner-API/Database"
	helpers "github.com/kiyoshi-87/url-shortner-API/Helpers"
	"github.com/redis/go-redis/v9"
)

type Request struct {
	Url         string        `json:"url"`
	CustomShort string        `json:"custom_short"`
	Expiry      time.Duration `json:"expiry"`
}

type Response struct {
	Url             string        `json:"url"`
	CustomShort     string        `json:"custom_short"`
	Expiry          time.Duration `json:"expiry"`
	XRateRemaining  int           `json:"rate_limit"`
	XRateLimitReset time.Duration `json:"rate_limit_reset"`
	ShortUrl        string        `json:"short_url"`
	CreatedAt       time.Time     `json:"created_at"`
}

func ShortenURL(c *fiber.Ctx) error {
	var body Request
	if err := c.BodyParser(&body); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "cannot pass request JSON"})
	}

	//Implement Rate Limiting (Limiting the number of API calls a particular IP can make)   10 calls every 30 mins

	client := database.CreateClient(1)
	defer client.Close()

	res, err := client.Get(context.Background(), c.IP()).Result()
	if err == redis.Nil {
		client.Set(context.Background(), c.IP(), os.Getenv("API_QUOTA"), time.Minute*30)
	} else {
		// val, _ := client.Get(context.Background(), c.IP()).Result()
		valInt, _ := strconv.Atoi(res)

		if valInt <= 0 {
			return c.Status(fiber.StatusServiceUnavailable).JSON(fiber.Map{"error": "Rate limit exceeded"})
		}

	}

	//Checking if the URL is valid

	if !govalidator.IsURL(body.Url) {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid URL"})
	}

	//Domain Error Check
	//This function will check for the domain name and return a boolean value.
	//If the domain name is not available, it will return false.
	//If the domain name is available, it will return true.

	if !helpers.RemoveDomainError(body.Url) {
		return c.Status(fiber.StatusServiceUnavailable).JSON(fiber.Map{"error": "Not available"})
	}

	//Enfore https

	body.Url = helpers.EnforceHttp(body.Url)

	//MAIN LOGIC SHORTENING URL

	var id string

	if body.CustomShort == "" {
		id = uuid.New().String()[:6]
	} else {
		id = body.CustomShort
	}

	r := database.CreateClient(0)
	defer r.Close()

	res1, _ := r.Get(context.Background(), id).Result()

	if res1 != "" {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"error": "Short URL already exists"})
	}

	if body.Expiry == 0 {
		body.Expiry = 24
	}

	_, err = r.Set(context.Background(), id, body.Url, body.Expiry*time.Hour).Result()

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Error saving url"})
	}

	//RESPONSE SENDING PART

	resp := Response{
		Url:             body.Url,
		CustomShort:     "",
		Expiry:          body.Expiry,
		XRateRemaining:  10,
		XRateLimitReset: time.Minute * 30,
		ShortUrl:        id,
		CreatedAt:       time.Now(),
	}

	//DECREMENTING THE API COUNTER

	client.Decr(context.Background(), c.IP())

	val, _ := client.Get(context.Background(), c.IP()).Result()
	resp.XRateRemaining, _ = strconv.Atoi(val)

	ttl, _ := client.TTL(context.Background(), c.IP()).Result()
	resp.XRateLimitReset = ttl * time.Nanosecond / time.Minute

	resp.ShortUrl = "localhost:3000/" + id

	return c.Status(fiber.StatusCreated).JSON(resp)

}
