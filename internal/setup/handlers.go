package setup

import (
	fiber "github.com/gofiber/fiber/v2"

	"petrichormud.com/app/internal/handlers"
	"petrichormud.com/app/internal/shared"
)

func Handlers(app *fiber.App, i *shared.Interfaces) {
	app.Get(handlers.HomeRoute, handlers.HomePage())

	app.Post(handlers.LoginRoute, handlers.Login(i))
	app.Get(handlers.LoginRoute, handlers.LoginPage())
	app.Post(handlers.LogoutRoute, handlers.Logout(i))
	app.Get(handlers.LogoutRoute, handlers.LogoutPage())

	app.Post(handlers.RegisterRoute, handlers.Register(i))
	app.Post(handlers.ReservedRoute, handlers.Reserved(i))

	app.Post(handlers.AddEmailRoute, handlers.AddEmail(i))
	app.Delete("/player/email/:id", handlers.DeleteEmail(i))
	app.Put("player/email/:id", handlers.EditEmail(i))
	app.Post("/player/email/:id/resend", handlers.ResendEmailVerification(i))

	// TODO: Move this behind the email group
	// TODO: Rename this to Verify and VerifyPage
	app.Get("/verify", handlers.Verify(i))
	app.Post("/verify", handlers.VerifyEmail(i))

	app.Get("/profile", handlers.ProfilePage(i))
}
