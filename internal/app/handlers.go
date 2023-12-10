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
	app.Get(routes.CharacterApplicationNamePath(":id"), handlers.CharacterNamePage(i))
	app.Get(routes.CharacterApplicationGenderPath(":id"), handlers.CharacterGenderPage(i))
	app.Get(routes.CharacterApplicationShortDescriptionPath(":id"), handlers.CharacterSdescPage(i))
	app.Get(routes.CharacterApplicationDescriptionPath(":id"), handlers.CharacterDescriptionPage(i))
	app.Get(routes.CharacterApplicationBackstoryPath(":id"), handlers.CharacterBackstoryPage(i))
	app.Post(routes.NewCharacterApplicationPath(), handlers.NewCharacterApplication(i))
	app.Put(routes.CharacterApplicationPath(":id"), handlers.UpdateCharacterApplication(i))
	app.Get(routes.CharacterApplicationNamePath(":id"), handlers.CharacterNamePage(i))
	app.Get(routes.CharacterApplicationGenderPath(":id"), handlers.CharacterGenderPage(i))
	app.Get(routes.CharacterApplicationShortDescriptionPath(":id"), handlers.CharacterSdescPage(i))
	app.Get(routes.CharacterApplicationDescriptionPath(":id"), handlers.CharacterDescriptionPage(i))
	app.Get(routes.CharacterApplicationBackstoryPath(":id"), handlers.CharacterBackstoryPage(i))
	app.Patch(routes.CharacterApplicationNamePath(":id"), handlers.UpdateCharacterApplicationName(i))
	app.Patch(routes.CharacterApplicationGenderPath(":id"), handlers.UpdateCharacterApplicationGender(i))
	app.Patch(routes.CharacterApplicationShortDescriptionPath(":id"), handlers.UpdateCharacterApplicationSdesc(i))
	app.Patch(routes.CharacterApplicationDescriptionPath(":id"), handlers.UpdateCharacterApplicationDescription(i))
	app.Patch(routes.CharacterApplicationBackstoryPath(":id"), handlers.UpdateCharacterApplicationBackstory(i))

	app.Post(routes.Login, handlers.Login(i))
	app.Get(routes.Login, handlers.LoginPage())
	app.Post(routes.Logout, handlers.Logout(i))
	app.Get(routes.Logout, handlers.LogoutPage())

	app.Post(routes.Register, handlers.Register(i))
	app.Post(routes.Reserved, handlers.Reserved(i))

	app.Post(routes.NewEmailPath(), handlers.AddEmail(i))
	app.Delete(routes.EmailPath(":id"), handlers.DeleteEmail(i))
	app.Put(routes.EmailPath(":id"), handlers.EditEmail(i))
	app.Post(routes.ResendEmailVerificationPath(":id"), handlers.Resend(i))

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
