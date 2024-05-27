package config

import (
	"os"
	"strconv"
	"time"

	fiber "github.com/gofiber/fiber/v2"
	"github.com/gofiber/template/html/v2"

	"petrichormud.com/app/internal/layout"
)

func Fiber(e *html.Engine) fiber.Config {
	// TODO: Error handling here
	readTimeoutSecondsCount, _ := strconv.Atoi(os.Getenv("SERVER_READ_TIMEOUT"))
	readTimeout := time.Second * time.Duration(readTimeoutSecondsCount)
	return fiber.Config{
		Views:       e,
		ViewsLayout: layout.Main,
		ReadTimeout: readTimeout,
	}
}
