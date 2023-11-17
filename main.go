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
	"petrichormud.com/app/internal/middleware/session"
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
	app.Use(session.New(&i))
	app.Use(bind.New())

	app.Static("/", "./web/static")
	app.Static("/loaders", "./web/svg/loaders")

	app.Get(handlers.HomeRoute, handlers.HomePage())

	app.Post(handlers.LoginRoute, handlers.Login(&i))
	app.Get(handlers.LoginRoute, handlers.LoginPage())
	app.Post(handlers.LogoutRoute, handlers.Logout(&i))
	app.Get(handlers.LogoutRoute, handlers.LogoutPage())

	app.Post(handlers.RegisterRoute, handlers.Register(&i))
	app.Post(handlers.ReservedRoute, handlers.Reserved(&i))

	app.Post(handlers.AddEmailRoute, handlers.AddEmail(&i))
	app.Delete(handlers.EmailRoute, handlers.DeleteEmail(&i))
	app.Put(handlers.EmailRoute, handlers.EditEmail(&i))
	app.Post(handlers.ResendRoute, handlers.Resend(&i))

	app.Get(handlers.VerifyRoute, handlers.VerifyPage(&i))
	app.Post(handlers.VerifyRoute, handlers.Verify(&i))

	app.Get(handlers.ProfileRoute, handlers.ProfilePage(&i))
	app.Get("/me", handlers.ProfilePage(&i))

	app.Get(handlers.RecoverUsernameRoute, handlers.RecoverUsernamePage(&i))
	app.Get(handlers.RecoverPasswordRoute, handlers.RecoverPasswordPage(&i))

	log.Fatal(app.Listen(":8008"))
}
