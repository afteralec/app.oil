package config

import (
	"github.com/gofiber/fiber/v2/middleware/session"

	"petrichormud.com/app/internal/util"
)

func Session() session.Config {
	if util.IsProd() {
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
