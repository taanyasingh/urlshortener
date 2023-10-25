package routes

import (
	"github.com/go-redis/redis/v8"
	"github.com/gofiber/fiber/v2"
	database "github.com/tanyasingh/urlshortener/db"
)

// Function redirecting the shorten url to original url
func RedirectAPI(context *fiber.Ctx) error {

	url := context.Params("url")
	// query the db to find the original URL
	r := database.CreateClient(0)
	defer r.Close()

	value, err := r.Get(database.Ctx, url).Result()
	if err == redis.Nil {
		return context.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "short url not found on database",
		})
	} else if err != nil {
		return context.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "cannot connect to database",
		})
	}

	return context.Redirect(value, 301)

}
