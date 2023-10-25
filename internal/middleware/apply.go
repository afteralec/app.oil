package middleware

import (
	fiber "github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
)

const HeaderName = "X-CSRF-TOKEN"

func Apply(a *fiber.App) {
	a.Use(cors.New())
	a.Use(logger.New())
	a.Use(Session())
}
