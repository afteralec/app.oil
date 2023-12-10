package app

import (
	fiber "github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/csrf"
	"github.com/gofiber/fiber/v2/middleware/logger"

	"petrichormud.com/app/internal/configs"
	"petrichormud.com/app/internal/middleware/bind"
	"petrichormud.com/app/internal/middleware/session"
	"petrichormud.com/app/internal/shared"
)

func Middleware(a *fiber.App, i *shared.Interfaces) {
	a.Use(logger.New())
	a.Use(csrf.New(configs.CSRF(i.Sessions)))
	a.Use(session.New(i))
	a.Use(bind.New())
}
