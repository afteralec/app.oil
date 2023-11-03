package main

import (
	"embed"
	"log"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
	fiber "github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/csrf"
	"github.com/gofiber/fiber/v2/middleware/logger"
	html "github.com/gofiber/template/html/v2"

	"petrichormud.com/app/internal/configs"
	"petrichormud.com/app/internal/handlers"
	"petrichormud.com/app/internal/middleware/bind"
	"petrichormud.com/app/internal/middleware/sessiondata"
	"petrichormud.com/app/internal/shared"
)

//go:embed web/views/*
var viewsfs embed.FS

func main() {
	i := shared.SetupInterfaces()
	defer i.Close()

	views := html.NewFileSystem(http.FS(viewsfs), ".html")
	app := fiber.New(configs.Fiber(views))

	app.Use(logger.New())
	app.Use(csrf.New(configs.CSRF(i.Sessions)))
	app.Use(sessiondata.New(&i))
	app.Use(bind.New())

	app.Static("/", "./web/static")
	app.Static("/loaders", "./web/svg/loaders")

	// TODO: Rename this to HomePage
	app.Get(handlers.HomeRoute, handlers.Home())

	app.Post(handlers.LoginRoute, handlers.Login(&i))
	app.Get(handlers.LoginRoute, handlers.LoginPage())
	app.Post(handlers.LogoutRoute, handlers.Logout(&i))
	app.Get(handlers.LogoutRoute, handlers.LogoutPage())

	app.Post(handlers.RegisterRoute, handlers.Register(&i))
	app.Post(handlers.ReservedRoute, handlers.Reserved(&i))

	app.Post(handlers.AddEmailRoute, handlers.AddEmail(&i))
	app.Delete("/player/email/:id", handlers.DeleteEmail(&i))
	app.Put("player/email/:id", handlers.EditEmail(&i))
	app.Post("/player/email/:id/resend", handlers.ResendEmailVerification(&i))

	// TODO: Move this behind the email group
	// TODO: Rename this to Verify and VerifyPage
	app.Get("/verify", handlers.Verify(&i))
	app.Post("/verify", handlers.VerifyEmail(&i))

	app.Get("/profile", handlers.ProfilePage(&i))

	app.Get(handlers.PermissionsRoute, handlers.PermissionsPage(&i))

	log.Fatal(app.Listen(":8008"))
}
