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
	app.Get(routes.CharacterApplicationNamePath(routes.ID), handlers.CharacterApplicationNamePage(i))
	app.Get(routes.CharacterApplicationGenderPath(routes.ID), handlers.CharacterApplicationGenderPage(i))
	app.Get(routes.CharacterApplicationShortDescriptionPath(routes.ID), handlers.CharacterApplicationShortDescriptionPage(i))
	app.Get(routes.CharacterApplicationDescriptionPath(routes.ID), handlers.CharacterApplicationDescriptionPage(i))
	app.Get(routes.CharacterApplicationBackstoryPath(routes.ID), handlers.CharacterApplicationBackstoryPage(i))
	app.Get(routes.CharacterApplicationReviewPath(routes.ID), handlers.CharacterApplicationReviewPage(i))

	app.Patch(routes.CharacterApplicationNamePath(routes.ID), handlers.UpdateCharacterApplicationName(i))
	app.Patch(routes.CharacterApplicationGenderPath(routes.ID), handlers.UpdateCharacterApplicationGender(i))
	app.Patch(routes.CharacterApplicationShortDescriptionPath(routes.ID), handlers.UpdateCharacterApplicationShortDescription(i))
	app.Patch(routes.CharacterApplicationDescriptionPath(routes.ID), handlers.UpdateCharacterApplicationDescription(i))
	app.Patch(routes.CharacterApplicationBackstoryPath(routes.ID), handlers.UpdateCharacterApplicationBackstory(i))

	app.Post(routes.SubmitCharacterApplicationPath(routes.ID), handlers.SubmitCharacterApplication(i))

	app.Get(routes.CharacterApplicationSubmittedSuccessPath(routes.ID), handlers.CharacterApplicationSubmittedSuccessPage(i))
	app.Get(routes.CharacterApplicationSubmittedPath(routes.ID), handlers.CharacterApplicationSubmittedPage(i))

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

	app.Get(routes.VerifyEmail, handlers.VerifyEmailPage(i))
	app.Post(routes.VerifyEmail, handlers.VerifyEmail(i))

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

	app.Get(routes.PlayerPermissions, handlers.PlayerPermissionsPage(i))
	app.Get(routes.PlayerPermissionsDetailPath(routes.Username), handlers.PlayerPermissionsDetailPage(i))
	app.Post(routes.PlayerPermissionsPath(routes.ID), handlers.UpdatePlayerPermission(i))

	app.Post(routes.SearchPlayerPath(routes.Destination), handlers.SearchPlayer(i))
}
