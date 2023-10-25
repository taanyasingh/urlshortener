package routes

import (
	"fmt"
	"os"
	"time"

	"github.com/asaskevich/govalidator"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	database "github.com/tanyasingh/urlshortener/db"
	"github.com/tanyasingh/urlshortener/helpers"
)

type request struct {
	URL         string        `json:"url"`
	CustomShort string        `json:"short"`
	Expiry      time.Duration `json:"expiry"`
}

type response struct {
	URL         string        `json:"url"`
	CustomShort string        `json:"short`
	Expiry      time.Duration `json:"expiry`
}

// Function for shortening url
func ShortenURL(context *fiber.Ctx) error {

	body := new(request)
	err := context.BodyParser(&body)
	if err != nil {
		return context.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Cannot parse JSON in request",
		})
	}

	//Checking url
	if !govalidator.IsURL(body.URL) {
		return context.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid URL",
		})
	}

	//check for the domain error
	if !helpers.RemoveDomainError(body.URL) {
		return context.Status(fiber.StatusServiceUnavailable).JSON(fiber.Map{
			"error": "Error in domain",
		})
	}

	//enforce https
	//all url will be converted to https before storing in database
	body.URL = helpers.EnforceHTTP(body.URL)

	// TODO: check if shorten url exists in db

	var id string
	if body.CustomShort == "" {
		id = uuid.New().String()[:6]
	} else {
		id = body.CustomShort
	}

	redisClient1 := database.CreateClient(0)
	defer redisClient1.Close()

	val, _ := redisClient1.Get(database.Ctx, id).Result()
	if val != "" {
		return context.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"error": "URL short already in use",
		})
	}
	if body.Expiry == 0 {
		body.Expiry = 24
	}

	err = redisClient1.Set(database.Ctx, id, body.URL, body.Expiry*3600*time.Second).Err()
	if err != nil {
		return context.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Unable to connect to server",
		})
	}
	resp := response{
		URL:         body.URL,
		CustomShort: "",
		Expiry:      body.Expiry,
	}

	//get domain name
	domain, err := helpers.GetDomainFromURL(body.URL)

	redisClient2 := database.CreateClient(1)
	defer redisClient2.Close()

	//Storing domain data for anlytical purpose
	_, err = redisClient2.ZIncrBy(database.Ctx, "domains", 1, domain).Result()
	if err != nil {
		fmt.Println(".........", err)
		return context.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Not able to add the set to databse",
		})
	}

	resp.CustomShort = os.Getenv("DOMAIN") + "/" + id
	return context.Status(fiber.StatusOK).JSON(resp)

}
