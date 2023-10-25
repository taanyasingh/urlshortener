package routes

import (
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

	metricData := make(map[string]int)
	for _, val := range result {
		str, _ := val.Member.(string)
		metricData[str] = int(val.Score)
	}
	return context.Status(fiber.StatusOK).JSON(metricData)
}
