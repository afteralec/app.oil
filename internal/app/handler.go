package app

import (
	fiber "github.com/gofiber/fiber/v2"

	"petrichormud.com/app/internal/handler"
	"petrichormud.com/app/internal/route"
	"petrichormud.com/app/internal/service"
)

func Handlers(app *fiber.App, i *service.Interfaces) {
	app.Get(route.Home, handler.HomePage())

	app.Post(route.ThemePathParam, handler.SetTheme(i))

	app.Get(route.Help, handler.HelpPage(i))
	app.Get(route.HelpFilePathParam, handler.HelpFilePage(i))
	app.Post(route.Help, handler.SearchHelp(i))

	app.Post(route.Requests, handler.CreateRequest(i))
	app.Get(route.RequestFieldPathParam, handler.RequestFieldPage(i))
	app.Patch(route.RequestFieldPathParam, handler.UpdateRequestField(i))
	app.Post(route.RequestFieldStatusPathParam, handler.UpdateRequestFieldStatus(i))
	app.Get(route.RequestPathParam, handler.RequestPage(i))
	app.Post(route.RequestStatusPathParam, handler.UpdateRequestStatus(i))

	app.Post(route.RequestChangeRequestFieldPathParam, handler.CreateRequestChangeRequest(i))
	app.Delete(route.RequestChangeRequestPathParam, handler.DeleteRequestChangeRequest(i))
	app.Put(route.RequestChangeRequestPathParam, handler.EditRequestChangeRequest(i))

	app.Post(route.Characters, handler.CreateCharacterApplication(i))
	app.Get(route.Characters, handler.CharactersPage(i))

	app.Get(route.CharacterApplications, handler.CharacterApplicationsQueuePage(i))

	app.Post(route.Login, handler.Login(i))
	app.Get(route.Login, handler.LoginPage())
	app.Post(route.Logout, handler.Logout(i))
	app.Get(route.Logout, handler.LogoutPage())

	app.Post(route.Register, handler.Register(i))
	// TODO: Should this be a GET with the search param in the URL?
	app.Post(route.Reserved, handler.UsernameReserved(i))

	app.Post(route.NewEmailPath(), handler.AddEmail(i))
	app.Delete(route.EmailPath(route.ID), handler.DeleteEmail(i))
	app.Put(route.EmailPath(route.ID), handler.EditEmail(i))
	app.Post(route.ResendEmailVerificationPath(route.ID), handler.ResendEmailVerification(i))

	app.Get(route.VerifyEmail, handler.VerifyEmailPage(i))
	app.Post(route.VerifyEmail, handler.VerifyEmail(i))

	app.Get(route.Profile, handler.ProfilePage(i))

	app.Put(route.PlayerPasswordParam, handler.ChangePassword(i))

	app.Get(route.Recover, handler.RecoverPage())

	app.Get(route.RecoverUsername, handler.RecoverUsernamePage())
	app.Post(route.RecoverUsername, handler.RecoverUsername(i))
	app.Get(route.RecoverUsernameSuccess, handler.RecoverUsernameSuccessPage(i))

	app.Get(route.RecoverPassword, handler.RecoverPasswordPage())
	app.Post(route.RecoverPassword, handler.RecoverPassword(i))
	app.Get(route.RecoverPasswordSuccess, handler.RecoverPasswordSuccessPage(i))

	app.Get(route.ResetPassword, handler.ResetPasswordPage())
	app.Post(route.ResetPassword, handler.ResetPassword(i))
	app.Get(route.ResetPasswordSuccess, handler.ResetPasswordSuccessPage())

	app.Get(route.PlayerPermissions, handler.PlayerPermissionsPage(i))
	app.Get(route.PlayerPermissionsDetailPath(route.Username), handler.PlayerPermissionsDetailPage(i))
	app.Post(route.PlayerPermissionsTogglePath(route.ID, route.Tag), handler.TogglePlayerPermission(i))

	app.Get(route.Rooms, handler.RoomsPage(i))
	app.Post(route.Rooms, handler.NewRoom(i))
	app.Get(route.RoomPathParam, handler.RoomPage(i))
	app.Get(route.EditRoomPathParam, handler.EditRoomPage(i))
	app.Get(route.RoomGridPathParam, handler.RoomGrid(i))
	app.Patch(route.RoomExitsPathParam, handler.EditRoomExit(i))
	app.Delete(route.RoomExitPathParam, handler.ClearRoomExit(i))
	app.Patch(route.RoomTitlePathParam, handler.EditRoomTitle(i))
	app.Patch(route.RoomDescriptionPathParam, handler.EditRoomDescription(i))
	app.Patch(route.RoomSizePathParam, handler.EditRoomSize(i))

	app.Post(route.ActorImageReserved, handler.ActorImageNameReserved(i))
	app.Post(route.ActorImages, handler.NewActorImage(i))
	app.Get(route.ActorImages, handler.ActorImagesPage(i))
	app.Get(route.ActorImagePathParam, handler.ActorImagePage(i))
	app.Get(route.EditActorImagePathParam, handler.EditActorImagePage(i))
	app.Patch(route.ActorImageShortDescriptionPathParam, handler.EditActorImageShortDescription(i))
	app.Patch(route.ActorImageDescriptionPathParam, handler.EditActorImageDescription(i))

	app.Post(route.SearchPlayerPath(route.Destination), handler.SearchPlayer(i))

	// TODO: Make this a conditional route based on the environment
	app.Get(route.DesignDictionary, handler.DesignDictionaryPage())
}
