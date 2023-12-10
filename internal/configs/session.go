package configs

import (
	"os"

	"github.com/gofiber/fiber/v2/middleware/session"
)

func Session() session.Config {
	if os.Getenv("PETRICHOR_APP_ENV") == "prod" {
		return session.Config{
			CookieHTTPOnly:    true,
			CookieSameSite:    "strict",
			CookieSecure:      true,
			CookieSessionOnly: true,
		}
	}
	return session.Config{
		CookieHTTPOnly:    true,
		CookieSameSite:    "strict",
		CookieSessionOnly: true,
	}
}
