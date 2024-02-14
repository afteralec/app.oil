package app

import (
	"os"

	fiber "github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/csrf"
	"github.com/gofiber/fiber/v2/middleware/logger"

	"petrichormud.com/app/internal/config"
	"petrichormud.com/app/internal/middleware/bind"
	"petrichormud.com/app/internal/middleware/permissions"
	"petrichormud.com/app/internal/middleware/session"
	"petrichormud.com/app/internal/service"
)

func Middleware(a *fiber.App, i *service.Interfaces) {
	if os.Getenv("DISABLE_LOGGING") != "true" {
		a.Use(logger.New())
	}

	// This order is important - if the CSRF middleware loads after bind, the CSRF token isn't sent to the templates
	// TODO: Figure out a way to test with the CSRF token
	if os.Getenv("DISABLE_CSRF") != "true" {
		a.Use(csrf.New(config.CSRF(i.Sessions)))
	}

	a.Use(session.New(i))
	a.Use(permissions.New(i))
	a.Use(bind.New())
}
