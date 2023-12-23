package bind

import (
	fiber "github.com/gofiber/fiber/v2"

	"petrichormud.com/app/internal/routes"
)

func CurrentView(b fiber.Map, c *fiber.Ctx) fiber.Map {
	b["HomeView"] = c.Path() == routes.Home
	b["ProfileView"] = c.Path() == routes.Profile || c.Path() == routes.Me
	b["CharactersView"] = c.Path() == routes.Characters
	b["PermissionsView"] = c.Path() == routes.PlayerPermissions
	return b
}
