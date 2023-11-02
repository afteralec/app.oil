package configs

import (
	"time"

	"github.com/gofiber/fiber/v2/middleware/csrf"
	"github.com/gofiber/fiber/v2/middleware/session"
	"github.com/gofiber/fiber/v2/utils"
)

const CSRFContextKey = "csrf"

func CSRF(s *session.Store) csrf.Config {
	return csrf.Config{
		KeyLookup:         "header:X-CSRF-Token",
		CookieName:        "csrf_",
		CookieSameSite:    "Lax",
		CookieSessionOnly: true,
		CookieHTTPOnly:    true,
		Expiration:        1 * time.Hour,
		KeyGenerator:      utils.UUIDv4,
		Session:           s,
		SessionKey:        "fiber.csrf.token",
		ContextKey:        CSRFContextKey,
		HandlerContextKey: "fiber.csrf.handler",
	}
}
