package app

import (
	fiber "github.com/gofiber/fiber/v2"

	"petrichormud.com/app/internal/handlers"
	"petrichormud.com/app/internal/routes"
	"petrichormud.com/app/internal/shared"
)

func Handlers(app *fiber.App, i *shared.Interfaces) {
	app.Get(handlers.HomeRoute, handlers.HomePage())

	app.Get(handlers.CharactersRoute, handlers.CharactersPage(i))
	app.Get(handlers.CharacterApplicationNameRoute, handlers.CharacterNamePage(i))
	app.Get(handlers.CharacterApplicationGenderRoute, handlers.CharacterGenderPage(i))
	app.Get(handlers.CharacterApplicationSdescRoute, handlers.CharacterSdescPage(i))
	app.Get(handlers.CharacterApplicationDescriptionRoute, handlers.CharacterDescriptionPage(i))
	app.Get(handlers.CharacterApplicationBackstoryRoute, handlers.CharacterBackstoryPage(i))
	app.Post(handlers.NewCharacterApplicationRoute, handlers.NewCharacterApplication(i))
	app.Put(handlers.CharacterApplicationRoute, handlers.UpdateCharacterApplication(i))
	app.Get(handlers.CharacterApplicationNameRoute, handlers.CharacterNamePage(i))
	app.Get(handlers.CharacterApplicationGenderRoute, handlers.CharacterGenderPage(i))
	app.Get(handlers.CharacterApplicationSdescRoute, handlers.CharacterSdescPage(i))
	app.Get(handlers.CharacterApplicationDescriptionRoute, handlers.CharacterDescriptionPage(i))
	app.Get(handlers.CharacterApplicationBackstoryRoute, handlers.CharacterBackstoryPage(i))
	app.Patch(handlers.CharacterApplicationNameRoute, handlers.UpdateCharacterApplicationName(i))
	app.Patch(handlers.CharacterApplicationGenderRoute, handlers.UpdateCharacterApplicationGender(i))
	app.Patch(handlers.CharacterApplicationSdescRoute, handlers.UpdateCharacterApplicationSdesc(i))
	app.Patch(handlers.CharacterApplicationDescriptionRoute, handlers.UpdateCharacterApplicationDescription(i))
	app.Patch(handlers.CharacterApplicationBackstoryRoute, handlers.UpdateCharacterApplicationBackstory(i))

	app.Post(handlers.LoginRoute, handlers.Login(i))
	app.Get(handlers.LoginRoute, handlers.LoginPage())
	app.Post(handlers.LogoutRoute, handlers.Logout(i))
	app.Get(handlers.LogoutRoute, handlers.LogoutPage())

	app.Post(handlers.RegisterRoute, handlers.Register(i))
	app.Post(handlers.ReservedRoute, handlers.Reserved(i))

	app.Post(routes.NewEmailPath(), handlers.AddEmail(i))
	app.Delete(routes.EmailPath(), handlers.DeleteEmail(i))
	app.Put(routes.EmailPath(), handlers.EditEmail(i))
	app.Post(handlers.ResendRoute, handlers.Resend(i))

	app.Get(handlers.VerifyRoute, handlers.VerifyPage(i))
	app.Post(handlers.VerifyRoute, handlers.Verify(i))

	app.Get(handlers.ProfileRoute, handlers.ProfilePage(i))
	app.Get("/me", handlers.ProfilePage(i))

	app.Get(handlers.RecoverRoute, handlers.RecoverPage())

	app.Get(handlers.RecoverUsernameRoute, handlers.RecoverUsernamePage())
	app.Post(handlers.RecoverUsernameRoute, handlers.RecoverUsername(i))
	app.Get(handlers.RecoverUsernameSuccessRoute, handlers.RecoverUsernameSuccessPage(i))

	app.Get(handlers.RecoverPasswordRoute, handlers.RecoverPasswordPage())
	app.Post(handlers.RecoverPasswordRoute, handlers.RecoverPassword(i))
	app.Get(handlers.RecoverPasswordSuccessRoute, handlers.RecoverPasswordSuccessPage())

	app.Get(handlers.ResetPasswordRoute, handlers.ResetPasswordPage())
	app.Post(handlers.ResetPasswordRoute, handlers.ResetPassword(i))
	app.Get(handlers.ResetPasswordSuccessRoute, handlers.ResetPasswordSuccessPage())
}
