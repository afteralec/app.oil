package app

import (
	fiber "github.com/gofiber/fiber/v2"

	"petrichormud.com/app/internal/handlers"
	"petrichormud.com/app/internal/routes"
	"petrichormud.com/app/internal/shared"
)

func Handlers(app *fiber.App, i *shared.Interfaces) {
	app.Get(routes.Home, handlers.HomePage())

	app.Post(routes.ThemePathParam, handlers.SetTheme(i))

	app.Post(routes.Requests, handlers.NewRequest(i))
	app.Get(routes.RequestFieldPathParam, handlers.RequestFieldPage(i))
	app.Patch(routes.RequestFieldPathParam, handlers.UpdateRequestFieldNew(i))
	app.Get(routes.RequestPathParam, handlers.RequestPage(i))
	app.Post(routes.RequestStatusPathParam, handlers.UpdateRequestStatus(i))

	app.Post(routes.CreateRequestCommentPath(routes.ID, routes.Field), handlers.CreateRequestComment(i))

	app.Post(routes.Characters, handlers.NewCharacterApplication(i))
	app.Get(routes.Characters, handlers.CharactersPage(i))

	app.Get(routes.CharacterApplications, handlers.CharacterApplicationsQueuePage(i))

	app.Post(routes.Login, handlers.Login(i))
	app.Get(routes.Login, handlers.LoginPage())
	app.Post(routes.Logout, handlers.Logout(i))
	app.Get(routes.Logout, handlers.LogoutPage())

	app.Post(routes.Register, handlers.Register(i))
	app.Post(routes.Reserved, handlers.Reserved(i))

	app.Post(routes.NewEmailPath(), handlers.AddEmail(i))
	app.Delete(routes.EmailPath(routes.ID), handlers.DeleteEmail(i))
	app.Put(routes.EmailPath(routes.ID), handlers.EditEmail(i))
	app.Post(routes.ResendEmailVerificationPath(routes.ID), handlers.ResendEmailVerification(i))

	app.Get(routes.VerifyEmail, handlers.VerifyEmailPage(i))
	app.Post(routes.VerifyEmail, handlers.VerifyEmail(i))

	app.Get(routes.Profile, handlers.ProfilePage(i))

	app.Get(routes.Recover, handlers.RecoverPage())

	app.Get(routes.RecoverUsername, handlers.RecoverUsernamePage())
	app.Post(routes.RecoverUsername, handlers.RecoverUsername(i))
	app.Get(routes.RecoverUsernameSuccess, handlers.RecoverUsernameSuccessPage(i))

	app.Get(routes.RecoverPassword, handlers.RecoverPasswordPage())
	app.Post(routes.RecoverPassword, handlers.RecoverPassword(i))
	app.Get(routes.RecoverPasswordSuccess, handlers.RecoverPasswordSuccessPage(i))

	app.Get(routes.ResetPassword, handlers.ResetPasswordPage())
	app.Post(routes.ResetPassword, handlers.ResetPassword(i))
	app.Get(routes.ResetPasswordSuccess, handlers.ResetPasswordSuccessPage())

	app.Get(routes.PlayerPermissions, handlers.PlayerPermissionsPage(i))
	app.Get(routes.PlayerPermissionsDetailPath(routes.Username), handlers.PlayerPermissionsDetailPage(i))
	app.Post(routes.PlayerPermissionsTogglePath(routes.ID, routes.Tag), handlers.TogglePlayerPermission(i))

	app.Post(routes.SearchPlayerPath(routes.Destination), handlers.SearchPlayer(i))
}
