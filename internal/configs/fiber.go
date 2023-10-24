package configs

import (
	"embed"
	"net/http"
	"os"
	"strconv"
	"time"

	fiber "github.com/gofiber/fiber/v2"
	html "github.com/gofiber/template/html/v2"
)

var viewsfs embed.FS

func Fiber() fiber.Config {
	readTimeoutSecondsCount, _ := strconv.Atoi(os.Getenv("SERVER_READ_TIMEOUT"))
	readTimeout := time.Second * time.Duration(readTimeoutSecondsCount)

	views := html.NewFileSystem(http.FS(viewsfs), ".html")

	return fiber.Config{
		Views:       views,
		ReadTimeout: readTimeout,
	}
}
