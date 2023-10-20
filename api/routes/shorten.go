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

func ShortenURL(context *fiber.Ctx) error {

	//incoming request parsing
	fmt.Println("testing....")
	body := new(request)
	fmt.Println(context)
	err := context.BodyParser(&body)
	if err != nil {
		return context.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "cannot parse JSON",
		})
	}

	//Check url
	if !govalidator.IsURL(body.URL) {
		return context.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid URL",
		})
	}

	// check for the domain error
	// users may abuse the shortener by shorting the domain `localhost:3000` itself
	// leading to a inifite loop, so don't accept the domain for shortening
	if !helpers.RemoveDomainError(body.URL) {
		return context.Status(fiber.StatusServiceUnavailable).JSON(fiber.Map{
			"error": "haha... nice try",
		})
	}

	// enforce https
	// all url will be converted to https before storing in database
	body.URL = helpers.EnforceHTTP(body.URL)

	//shortening
	// LOGIC FOR SHORTENING :
	//check if the user has provided any custom dhort urls
	// if yes, proceed,
	// else, create a new short using the first 6 digits of uuid
	// TODO: collision checks

	var id string
	if body.CustomShort == "" {
		id = uuid.New().String()[:6]
	} else {
		id = body.CustomShort
	}

	r := database.CreateClient(0)
	defer r.Close()

	val, _ := r.Get(database.Ctx, id).Result()

	if val != "" {
		return context.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"error": "URL short already in use",
		})
	}
	if body.Expiry == 0 {
		body.Expiry = 24 // default expiry of 24 hours TODO: logic for expiry
	}

	resp := response{
		URL:         body.URL,
		CustomShort: "",
		Expiry:      body.Expiry,
	}

	resp.CustomShort = os.Getenv("DOMAIN") + "/" + id
	return context.Status(fiber.StatusOK).JSON(resp)
	//TODO:add to db

}
