package middleware

import (
	"time"

	fiber "github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/csrf"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/utils"
)

const HeaderName = "X-CSRF-TOKEN"

func Apply(a *fiber.App) {
	a.Use(cors.New())
	a.Use(logger.New())
	a.Use(Session())
	a.Use(csrf.New(csrf.Config{
		KeyLookup:         "header:" + HeaderName,
		CookieName:        "csrf_",
		CookieSameSite:    "Lax",
		CookieSessionOnly: true,
		CookieHTTPOnly:    true,
		Expiration:        1 * time.Hour,
		KeyGenerator:      utils.UUIDv4,
		Session:           Sessions,
		SessionKey:        "fiber.csrf.token",
		HandlerContextKey: "fiber.csrf.handler",
	}))
}
