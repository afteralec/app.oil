package handler

import (
	fiber "github.com/gofiber/fiber/v2"

	"petrichormud.com/app/internal/layouts"
	"petrichormud.com/app/internal/views"
)

func DesignDictionaryPage() fiber.Handler {
	return func(c *fiber.Ctx) error {
		return c.Render(views.DesignDictionary, views.Bind(c), layouts.Main)
	}
}
