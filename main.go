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

	app.Get(handlers.HomeRoute, handlers.Home())

	app.Post(handlers.LoginRoute, handlers.Login(&i))
	app.Post("/logout", handlers.Logout(&i))
	app.Get("/logout", handlers.LogoutPage())

	player := app.Group("/player")
	player.Post("/new", handlers.CreatePlayer(&i))
	player.Post("/reserved", handlers.UsernameReserved(&i))

	email := player.Group("/email")
	email.Post("/new", handlers.AddEmail(&i))
	email.Delete("/:id", handlers.DeleteEmail(&i))
	email.Put("/:id", handlers.EditEmail(&i))
	email.Post("/:id/resend", handlers.ResendEmailVerification(&i))

	// TODO: Move this behind the email group
	app.Get("/verify", handlers.Verify(&i))
	app.Post("/verify", handlers.VerifyEmail(&i))

	app.Get("/profile", handlers.Profile(&i))

	log.Fatal(app.Listen(":8008"))
}
