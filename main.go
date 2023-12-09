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

	app.Get(handlers.CharactersRoute, handlers.CharactersPage(&i))
	app.Get(handlers.CharacterApplicationNameRoute, handlers.CharacterNamePage(&i))
	app.Get(handlers.CharacterApplicationGenderRoute, handlers.CharacterGenderPage(&i))
	app.Get(handlers.CharacterApplicationSdescRoute, handlers.CharacterSdescPage(&i))
	app.Get(handlers.CharacterApplicationDescriptionRoute, handlers.CharacterDescriptionPage(&i))
	app.Get(handlers.CharacterApplicationBackstoryRoute, handlers.CharacterBackstoryPage(&i))
	app.Post(handlers.NewCharacterApplicationRoute, handlers.NewCharacterApplication(&i))
	app.Put(handlers.CharacterApplicationRoute, handlers.UpdateCharacterApplication(&i))
	app.Get(handlers.CharacterApplicationNameRoute, handlers.CharacterNamePage(&i))
	app.Get(handlers.CharacterApplicationGenderRoute, handlers.CharacterGenderPage(&i))
	app.Get(handlers.CharacterApplicationSdescRoute, handlers.CharacterSdescPage(&i))
	app.Get(handlers.CharacterApplicationDescriptionRoute, handlers.CharacterDescriptionPage(&i))
	app.Get(handlers.CharacterApplicationBackstoryRoute, handlers.CharacterBackstoryPage(&i))
	app.Patch(handlers.CharacterApplicationNameRoute, handlers.UpdateCharacterApplicationName(&i))
	app.Patch(handlers.CharacterApplicationGenderRoute, handlers.UpdateCharacterApplicationGender(&i))
	app.Patch(handlers.CharacterApplicationSdescRoute, handlers.UpdateCharacterApplicationSdesc(&i))
	app.Patch(handlers.CharacterApplicationDescriptionRoute, handlers.UpdateCharacterApplicationDescription(&i))
	app.Patch(handlers.CharacterApplicationBackstoryRoute, handlers.UpdateCharacterApplicationBackstory(&i))

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

	app.Get(handlers.RecoverRoute, handlers.RecoverPage())

	app.Get(handlers.RecoverUsernameRoute, handlers.RecoverUsernamePage())
	app.Post(handlers.RecoverUsernameRoute, handlers.RecoverUsername(&i))
	app.Get(handlers.RecoverUsernameSuccessRoute, handlers.RecoverUsernameSuccessPage(&i))

	app.Get(handlers.RecoverPasswordRoute, handlers.RecoverPasswordPage())
	app.Post(handlers.RecoverPasswordRoute, handlers.RecoverPassword(&i))
	app.Get(handlers.RecoverPasswordSuccessRoute, handlers.RecoverPasswordSuccessPage())

	app.Get(handlers.ResetPasswordRoute, handlers.ResetPasswordPage())
	app.Post(handlers.ResetPasswordRoute, handlers.ResetPassword(&i))
	app.Get(handlers.ResetPasswordSuccessRoute, handlers.ResetPasswordSuccessPage())

	log.Fatal(app.Listen(":8008"))
}
