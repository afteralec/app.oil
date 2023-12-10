package configs

import (
	"time"

	"github.com/gofiber/fiber/v2/middleware/csrf"
	"github.com/gofiber/fiber/v2/middleware/session"
	"github.com/gofiber/fiber/v2/utils"
)

func CSRF(s *session.Store) csrf.Config {
	return csrf.Config{
		KeyLookup:         "header:X-CSRF-Token",
		CookieName:        "csrf_",
		CookieSameSite:    "Strict",
		CookieSessionOnly: true,
		CookieHTTPOnly:    true,
		Expiration:        1 * time.Hour,
		KeyGenerator:      utils.UUIDv4,
		Session:           s,
		SessionKey:        "petrichor.csrf.token",
		ContextKey:        "csrf",
		HandlerContextKey: "petrichor.csrf.handler",
	}
}
