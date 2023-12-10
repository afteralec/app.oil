package app

import (
	fiber "github.com/gofiber/fiber/v2"

	"petrichormud.com/app/internal/handlers"
	"petrichormud.com/app/internal/routes"
	"petrichormud.com/app/internal/shared"
)

func Handlers(app *fiber.App, i *shared.Interfaces) {
	app.Get(routes.Home, handlers.HomePage())

	app.Get(routes.Characters, handlers.CharactersPage(i))
	app.Post(routes.NewCharacterApplicationPath(), handlers.NewCharacterApplication(i))
	app.Put(routes.CharacterApplicationPath(routes.ID), handlers.UpdateCharacterApplication(i))
	app.Get(routes.CharacterApplicationNamePath(routes.ID), handlers.CharacterApplicationNamePage(i))
	app.Get(routes.CharacterApplicationGenderPath(routes.ID), handlers.CharacterGenderPage(i))
	app.Get(routes.CharacterApplicationShortDescriptionPath(routes.ID), handlers.CharacterShortDescriptionPage(i))
	app.Get(routes.CharacterApplicationDescriptionPath(routes.ID), handlers.CharacterDescriptionPage(i))
	app.Get(routes.CharacterApplicationBackstoryPath(routes.ID), handlers.CharacterBackstoryPage(i))
	app.Patch(routes.CharacterApplicationNamePath(routes.ID), handlers.UpdateCharacterApplicationName(i))
	app.Patch(routes.CharacterApplicationGenderPath(routes.ID), handlers.UpdateCharacterApplicationGender(i))
	app.Patch(routes.CharacterApplicationShortDescriptionPath(routes.ID), handlers.UpdateCharacterApplicationShortDescription(i))
	app.Patch(routes.CharacterApplicationDescriptionPath(routes.ID), handlers.UpdateCharacterApplicationDescription(i))
	app.Patch(routes.CharacterApplicationBackstoryPath(routes.ID), handlers.UpdateCharacterApplicationBackstory(i))

	app.Post(routes.Login, handlers.Login(i))
	app.Get(routes.Login, handlers.LoginPage())
	app.Post(routes.Logout, handlers.Logout(i))
	app.Get(routes.Logout, handlers.LogoutPage())

	app.Post(routes.Register, handlers.Register(i))
	app.Post(routes.Reserved, handlers.Reserved(i))

	app.Post(routes.NewEmailPath(), handlers.AddEmail(i))
	app.Delete(routes.EmailPath(routes.ID), handlers.DeleteEmail(i))
	app.Put(routes.EmailPath(routes.ID), handlers.EditEmail(i))
	app.Post(routes.ResendEmailVerificationPath(routes.ID), handlers.Resend(i))

	app.Get(routes.VerifyEmail, handlers.VerifyPage(i))
	app.Post(routes.VerifyEmail, handlers.Verify(i))

	app.Get(routes.Profile, handlers.ProfilePage(i))
	app.Get(routes.Me, handlers.ProfilePage(i))

	app.Get(routes.Recover, handlers.RecoverPage())

	app.Get(routes.RecoverUsername, handlers.RecoverUsernamePage())
	app.Post(routes.RecoverUsername, handlers.RecoverUsername(i))
	app.Get(routes.RecoverUsernameSuccess, handlers.RecoverUsernameSuccessPage(i))

	app.Get(routes.RecoverPassword, handlers.RecoverPasswordPage())
	app.Post(routes.RecoverPassword, handlers.RecoverPassword(i))
	app.Get(routes.RecoverPasswordSuccess, handlers.RecoverPasswordSuccessPage())

	app.Get(routes.ResetPassword, handlers.ResetPasswordPage())
	app.Post(routes.ResetPassword, handlers.ResetPassword(i))
	app.Get(routes.ResetPasswordSuccess, handlers.ResetPasswordSuccessPage())
}
