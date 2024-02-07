package routes

import (
	"context"

	"github.com/gofiber/fiber/v2"
	database "github.com/kiyoshi-87/url-shortner-API/Database"

	"github.com/redis/go-redis/v9"
)

//THIS FILE MAKES THE MAIN CONNECTION TO THE REDIS DATABASE
//ALSO IT MAKES THE KEY VALUE PAIRS IN THE DATABASE WHERE THE SHORTENED URL IS MAPPED TO THE ACTUAL URL

func ResolveURL(c *fiber.Ctx) error {
	url := c.Params("url")

	r := database.CreateClient(0)
	defer r.Close()

	//Now using this client we can search the database for the given URL

	res, err := r.Get(context.Background(), url).Result()

	//FIRST ERROR SIGNIFIES THE URL WAS NOT FOUND IN THE DATABASE
	//SECOND ERROR SIGNIFIES THERER WAS ERROR CONNECTING TO THE DATABASE IN THE FIRST PLACE, GOOD FOR DEBUGGING
	if err == redis.Nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Shortened URL not found in the database"})
	} else if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Internal Server Error"})
	}

	//IN THE 1-INDEXED REDIS DB THE COUNTER KEY IS STORED WHERE WHERE WE ARE COUNTING THE NUMBER OF API CALLS TAKEN
	rInr := database.CreateClient(1)
	defer rInr.Close()

	rInr.Incr(context.Background(), "counter")

	return c.Redirect(res, 301)
}
