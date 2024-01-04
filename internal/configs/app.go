package configs

import (
	"os"
	"strconv"
	"time"

	fiber "github.com/gofiber/fiber/v2"

	"petrichormud.com/app/internal/layouts"
	"petrichormud.com/app/web"
)

func Fiber() fiber.Config {
	// TODO: Error handling here
	readTimeoutSecondsCount, _ := strconv.Atoi(os.Getenv("SERVER_READ_TIMEOUT"))
	readTimeout := time.Second * time.Duration(readTimeoutSecondsCount)
	return fiber.Config{
		Views:       web.ViewsEngine(),
		ViewsLayout: layouts.Main,
		ReadTimeout: readTimeout,
	}
}
