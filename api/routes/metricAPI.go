package routes

import (
	"strconv"

	"github.com/gofiber/fiber/v2"
	database "github.com/tanyasingh/urlshortener/db"
)

func MetricAPI(context *fiber.Ctx) error {
	//metrics API returns top 3 domain names t
	//hat have been shortened the most
	//number of times

	redisClient := database.CreateClient(1)
	defer redisClient.Close()
	result, err := redisClient.ZRevRangeWithScores(database.Ctx, "domains", 0, 3).Result()
	if err != nil {
		return context.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Not able to connect to db",
		})
	}

	metricsData := []string{}

	for _, value := range result {
		str := value.Member.(string)
		domainItem := str + ": " + strconv.Itoa(int(value.Score))
		metricsData = append(metricsData, domainItem)
	}

	return context.Status(fiber.StatusOK).JSON(metricsData)
}
