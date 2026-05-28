package logs

import (
	"encoding/json"
	"moyo-gateway-service/utils"
	"os"

	"github.com/gofiber/fiber/v2"
)

// LogEntry represents a single log entry

func Handler(c *fiber.Ctx) error {
	date := c.Params("date")
	if date == "" {
		return c.Status(400).JSON(fiber.Map{
			"IsSuccess": false,
			"Message":   "Date parameter is required",
			"Data":      nil,
			"Status":    "1",
		})
	}

	// Validate date format (YYYYMMDD)
	if len(date) != 8 {
		return c.Status(400).JSON(fiber.Map{
			"IsSuccess": false,
			"Message":   "Invalid date format. Use YYYYMMDD",
			"Data":      nil,
			"Status":    "1",
		})
	}

	// Read log file
	logDir := "./logs"
	logFile := logDir + "/" + date + ".json"

	data, err := os.ReadFile(logFile)
	if err != nil {
		utils.PushLogf("", "GetLogs", err.Error())
		return err
	}

	// logFile has json format. here's the format
	// /{"logs":[{"date":"2025-05-23 06:39:28","log":"[Request] GET /logs/20250523 , "}]}
	var entries LogsResponse

	err = json.Unmarshal(data, &entries)
	if err != nil {
		utils.PushLogf("", "GetLogs", err.Error())
		return err
	}

	return c.Status(200).JSON(fiber.Map{
		"IsSuccess": true,
		"Message":   "Logs retrieved successfully",
		"Data":      entries,
		"Status":    "0",
	})
}
