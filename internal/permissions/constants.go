package permissions

type Player struct {
	Name  string
	Tag   string
	Title string
	About string
}

const (
	PlayerGrantAllPermissionsName         string = "GrantAllPermissions"
	PlayerRevokeAllPermissionsName        string = "RevokeAllPermissions"
	PlayerReviewCharacterApplicationsName string = "ReviewCharacterApplications"
	PlayerViewAllRoomImagesName           string = "ViewAllRoomImages"
	PlayerViewAllRoomsName                string = "ViewAllRooms"
	PlayerCreateRoomImageName             string = "CreateRoomImage"
	PlayerEditRoomImageName               string = "EditRoomImage"
	PlayerCreateRoomName                  string = "CreateRoom"
)

const (
	PlayerGrantAllPermissionsTag         string = "grant-all"
	PlayerRevokeAllPermissionsTag        string = "revoke-all"
	PlayerReviewCharacterApplicationsTag string = "review-character-applications"
	PlayerViewAllRoomImagesTag           string = "view-all-room-images"
	PlayerViewAllRoomsTag                string = "view-all-rooms"
	PlayerCreateRoomImageTag             string = "create-room-image"
	PlayerEditRoomImageTag               string = "edit-room-image"
	PlayerCreateRoomTag                  string = "create-room"
)

const (
	PlayerGrantAllPermissionsTitle         string = "Grant All Permissions"
	PlayerRevokeAllPermissionsTitle        string = "Revoke All Permissions"
	PlayerReviewCharacterApplicationsTitle string = "Review Character Applications"
	PlayerViewAllRoomImagesTitle           string = "View All Room Images"
	PlayerViewAllRoomsTitle                string = "View All Rooms"
	PlayerCreateRoomImageTitle             string = "Create Room Images"
	PlayerEditRoomImageTitle               string = "Edit Room Images"
	PlayerCreateRoomTitle                  string = "Create Room"
)

const (
	PlayerGrantAllPermissionsAbout         string = "The root permission. Only one person should have this at a time."
	PlayerRevokeAllPermissionsAbout        string = "The root revocation permission. Only one person should have this at a time."
	PlayerReviewCharacterApplicationsAbout string = "Enable this player to review Character Applications."
	PlayerViewAllRoomImagesAbout           string = "The permission to view (but not edit) all room image data."
	PlayerViewAllRoomsAbout                string = "The permission to view (but not edit) all room data."
	PlayerCreateRoomImageAbout             string = "Create new room images."
	PlayerEditRoomImageAbout               string = "Edit room images."
	PlayerCreateRoomAbout                  string = "Create a new room, but not connect it to the grid."
)

var PlayerGrantAllPermissions Player = Player{
	Name:  PlayerGrantAllPermissionsName,
	Tag:   PlayerGrantAllPermissionsTag,
	Title: PlayerGrantAllPermissionsTitle,
	About: PlayerGrantAllPermissionsAbout,
}

var PlayerRevokeAllPermissions Player = Player{
	Name:  PlayerRevokeAllPermissionsName,
	Tag:   PlayerRevokeAllPermissionsTag,
	Title: PlayerRevokeAllPermissionsTitle,
	About: PlayerRevokeAllPermissionsAbout,
}

var PlayerReviewCharacterApplications Player = Player{
	Name:  PlayerReviewCharacterApplicationsName,
	Tag:   PlayerReviewCharacterApplicationsTag,
	Title: PlayerReviewCharacterApplicationsTitle,
	About: PlayerReviewCharacterApplicationsAbout,
}

var PlayerViewAllRoomImages Player = Player{
	Name:  PlayerViewAllRoomImagesName,
	Tag:   PlayerViewAllRoomImagesTag,
	Title: PlayerViewAllRoomImagesTitle,
	About: PlayerViewAllRoomImagesAbout,
}

var PlayerViewAllRooms Player = Player{
	Name:  PlayerViewAllRoomsName,
	Tag:   PlayerViewAllRoomsTag,
	Title: PlayerViewAllRoomsTitle,
	About: PlayerViewAllRoomsAbout,
}

var PlayerCreateRoomImage Player = Player{
	Name:  PlayerCreateRoomImageName,
	Tag:   PlayerCreateRoomImageTag,
	Title: PlayerCreateRoomImageTitle,
	About: PlayerCreateRoomImageAbout,
}

var PlayerEditRoomImage Player = Player{
	Name:  PlayerEditRoomImageName,
	Tag:   PlayerEditRoomImageTag,
	Title: PlayerEditRoomImageTitle,
	About: PlayerEditRoomImageAbout,
}

var PlayerCreateRoom Player = Player{
	Name:  PlayerCreateRoomName,
	Tag:   PlayerCreateRoomTag,
	Title: PlayerCreateRoomTitle,
	About: PlayerCreateRoomAbout,
}

var ShowPermissionViewPermissions []string = []string{
	PlayerGrantAllPermissionsName,
}

var AllPlayer []Player = []Player{
	PlayerGrantAllPermissions,
	PlayerRevokeAllPermissions,
	PlayerReviewCharacterApplications,
	PlayerViewAllRoomImages,
	PlayerViewAllRooms,
	PlayerCreateRoomImage,
	PlayerEditRoomImage,
	PlayerCreateRoom,
}

var AllPlayerByName map[string]Player = map[string]Player{
	PlayerGrantAllPermissionsName:         PlayerGrantAllPermissions,
	PlayerRevokeAllPermissionsName:        PlayerRevokeAllPermissions,
	PlayerReviewCharacterApplicationsName: PlayerReviewCharacterApplications,
	PlayerViewAllRoomImagesName:           PlayerViewAllRoomImages,
	PlayerViewAllRoomsName:                PlayerViewAllRooms,
	PlayerCreateRoomImageName:             PlayerCreateRoomImage,
	PlayerEditRoomImageName:               PlayerEditRoomImage,
	PlayerCreateRoomName:                  PlayerCreateRoom,
}

var AllPlayerByTag map[string]Player = map[string]Player{
	PlayerGrantAllPermissionsTag:         PlayerGrantAllPermissions,
	PlayerRevokeAllPermissionsTag:        PlayerRevokeAllPermissions,
	PlayerReviewCharacterApplicationsTag: PlayerReviewCharacterApplications,
	PlayerViewAllRoomImagesTag:           PlayerViewAllRoomImages,
	PlayerViewAllRoomsTag:                PlayerViewAllRooms,
	PlayerCreateRoomImageTag:             PlayerCreateRoomImage,
	PlayerEditRoomImageTag:               PlayerEditRoomImage,
	PlayerCreateRoomTag:                  PlayerCreateRoom,
}

var RootPlayerByName map[string]Player = map[string]Player{
	PlayerGrantAllPermissionsName:  PlayerGrantAllPermissions,
	PlayerRevokeAllPermissionsName: PlayerRevokeAllPermissions,
}
