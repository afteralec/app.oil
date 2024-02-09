package app

import (
	fiber "github.com/gofiber/fiber/v2"

	"petrichormud.com/app/internal/handler"
	"petrichormud.com/app/internal/interfaces"
	"petrichormud.com/app/internal/routes"
)

func Handlers(app *fiber.App, i *interfaces.Shared) {
	app.Get(routes.Home, handler.HomePage())

	app.Post(routes.ThemePathParam, handler.SetTheme(i))

	app.Get(routes.Help, handler.HelpPage(i))
	app.Get(routes.HelpFilePathParam, handler.HelpFilePage(i))
	app.Post(routes.Help, handler.SearchHelp(i))

	app.Post(routes.Requests, handler.NewRequest(i))
	app.Get(routes.RequestFieldPathParam, handler.RequestFieldPage(i))
	app.Patch(routes.RequestFieldPathParam, handler.UpdateRequestField(i))
	app.Get(routes.RequestPathParam, handler.RequestPage(i))
	app.Post(routes.RequestStatusPathParam, handler.UpdateRequestStatus(i))

	app.Post(routes.CreateRequestCommentPath(routes.ID, routes.Field), handler.CreateRequestComment(i))

	app.Post(routes.Characters, handler.NewCharacterApplication(i))
	app.Get(routes.Characters, handler.CharactersPage(i))

	app.Get(routes.CharacterApplications, handler.CharacterApplicationsQueuePage(i))

	app.Post(routes.Login, handler.Login(i))
	app.Get(routes.Login, handler.LoginPage())
	app.Post(routes.Logout, handler.Logout(i))
	app.Get(routes.Logout, handler.LogoutPage())

	app.Post(routes.Register, handler.Register(i))
	// TODO: Should this be a GET with the search param in the URL?
	app.Post(routes.Reserved, handler.UsernameReserved(i))

	app.Post(routes.NewEmailPath(), handler.AddEmail(i))
	app.Delete(routes.EmailPath(routes.ID), handler.DeleteEmail(i))
	app.Put(routes.EmailPath(routes.ID), handler.EditEmail(i))
	app.Post(routes.ResendEmailVerificationPath(routes.ID), handler.ResendEmailVerification(i))

	app.Get(routes.VerifyEmail, handler.VerifyEmailPage(i))
	app.Post(routes.VerifyEmail, handler.VerifyEmail(i))

	app.Get(routes.Profile, handler.ProfilePage(i))

	app.Put(routes.PlayerPasswordParam, handler.ChangePassword(i))

	app.Get(routes.Recover, handler.RecoverPage())

	app.Get(routes.RecoverUsername, handler.RecoverUsernamePage())
	app.Post(routes.RecoverUsername, handler.RecoverUsername(i))
	app.Get(routes.RecoverUsernameSuccess, handler.RecoverUsernameSuccessPage(i))

	app.Get(routes.RecoverPassword, handler.RecoverPasswordPage())
	app.Post(routes.RecoverPassword, handler.RecoverPassword(i))
	app.Get(routes.RecoverPasswordSuccess, handler.RecoverPasswordSuccessPage(i))

	app.Get(routes.ResetPassword, handler.ResetPasswordPage())
	app.Post(routes.ResetPassword, handler.ResetPassword(i))
	app.Get(routes.ResetPasswordSuccess, handler.ResetPasswordSuccessPage())

	app.Get(routes.PlayerPermissions, handler.PlayerPermissionsPage(i))
	app.Get(routes.PlayerPermissionsDetailPath(routes.Username), handler.PlayerPermissionsDetailPage(i))
	app.Post(routes.PlayerPermissionsTogglePath(routes.ID, routes.Tag), handler.TogglePlayerPermission(i))

	app.Get(routes.Rooms, handler.RoomsPage(i))
	app.Post(routes.Rooms, handler.NewRoom(i))
	app.Get(routes.RoomPathParam, handler.RoomPage(i))
	app.Get(routes.EditRoomPathParam, handler.EditRoomPage(i))
	app.Get(routes.RoomGridPathParam, handler.RoomGrid(i))
	app.Patch(routes.RoomExitsPathParam, handler.EditRoomExit(i))
	app.Delete(routes.RoomExitPathParam, handler.ClearRoomExit(i))
	app.Patch(routes.RoomTitlePathParam, handler.EditRoomTitle(i))
	app.Patch(routes.RoomDescriptionPathParam, handler.EditRoomDescription(i))
	app.Patch(routes.RoomSizePathParam, handler.EditRoomSize(i))

	app.Post(routes.ActorImageReserved, handler.ActorImageNameReserved(i))
	app.Post(routes.ActorImages, handler.NewActorImage(i))
	app.Get(routes.ActorImages, handler.ActorImagesPage(i))
	app.Get(routes.ActorImagePathParam, handler.ActorImagePage(i))
	app.Get(routes.EditActorImagePathParam, handler.EditActorImagePage(i))
	app.Patch(routes.ActorImageShortDescriptionPathParam, handler.EditActorImageShortDescription(i))
	app.Patch(routes.ActorImageDescriptionPathParam, handler.EditActorImageDescription(i))

	app.Post(routes.SearchPlayerPath(routes.Destination), handler.SearchPlayer(i))

	app.Get(routes.DesignDictionary, handler.DesignDictionaryPage())
}
