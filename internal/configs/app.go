package configs

import (
	"os"
	"strconv"
	"time"

	fiber "github.com/gofiber/fiber/v2"
)

func Fiber(views fiber.Views) fiber.Config {
	readTimeoutSecondsCount, _ := strconv.Atoi(os.Getenv("SERVER_READ_TIMEOUT"))
	readTimeout := time.Second * time.Duration(readTimeoutSecondsCount)
	return fiber.Config{
		Views:       views,
		ViewsLayout: "web/views/layouts/main",
		ReadTimeout: readTimeout,
	}
}
