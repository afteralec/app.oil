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
	PlayerViewAllRoomsName                string = "ViewAllRooms"
	PlayerCreateRoomName                  string = "CreateRoom"
	PlayerViewAllActorImagesName          string = "ViewAllActorImages"
	PlayerCreateActorImageName            string = "CreateActorImage"
)

const (
	PlayerGrantAllPermissionsTag         string = "grant-all"
	PlayerRevokeAllPermissionsTag        string = "revoke-all"
	PlayerReviewCharacterApplicationsTag string = "review-character-applications"
	PlayerViewAllRoomsTag                string = "view-all-rooms"
	PlayerCreateRoomTag                  string = "create-room"
	PlayerViewAllActorImagesTag          string = "view-all-actor-images"
	PlayerCreateActorImageTag            string = "create-actor-image"
)

const (
	PlayerGrantAllPermissionsTitle         string = "Grant All Permissions"
	PlayerRevokeAllPermissionsTitle        string = "Revoke All Permissions"
	PlayerReviewCharacterApplicationsTitle string = "Review Character Applications"
	PlayerViewAllRoomsTitle                string = "View All Rooms"
	PlayerCreateRoomTitle                  string = "Create Room"
	PlayerViewAllActorImagesTitle          string = "View All Actor Images"
	PlayerCreateActorImageTitle            string = "Create Actor Image"
)

const (
	PlayerGrantAllPermissionsAbout         string = "The root permission. Only one person should have this at a time."
	PlayerRevokeAllPermissionsAbout        string = "The root revocation permission. Only one person should have this at a time."
	PlayerReviewCharacterApplicationsAbout string = "Enable this player to review Character Applications."
	PlayerViewAllRoomsAbout                string = "The permission to view (but not edit) all room data."
	PlayerCreateRoomAbout                  string = "Create a new room, but not connect it to the grid."
	PlayerViewAllActorImagesAbout          string = "View all Actor Images, i.e. in the main Actor Images list."
	PlayerCreateActorImageAbout            string = "Create new actor via creating new Actor Images"
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

var PlayerViewAllRooms Player = Player{
	Name:  PlayerViewAllRoomsName,
	Tag:   PlayerViewAllRoomsTag,
	Title: PlayerViewAllRoomsTitle,
	About: PlayerViewAllRoomsAbout,
}

var PlayerCreateRoom Player = Player{
	Name:  PlayerCreateRoomName,
	Tag:   PlayerCreateRoomTag,
	Title: PlayerCreateRoomTitle,
	About: PlayerCreateRoomAbout,
}

var PlayerViewAllActorImages Player = Player{
	Name:  PlayerViewAllActorImagesName,
	Tag:   PlayerViewAllActorImagesTag,
	Title: PlayerViewAllActorImagesTitle,
	About: PlayerViewAllActorImagesAbout,
}

var PlayerCreateActorImage Player = Player{
	Name:  PlayerCreateActorImageName,
	Tag:   PlayerCreateActorImageTag,
	Title: PlayerCreateActorImageTitle,
	About: PlayerCreateActorImageAbout,
}

var ShowPermissionViewPermissions []string = []string{
	PlayerGrantAllPermissionsName,
}

var AllPlayer []Player = []Player{
	PlayerGrantAllPermissions,
	PlayerRevokeAllPermissions,
	PlayerReviewCharacterApplications,
	PlayerViewAllRooms,
	PlayerCreateRoom,
	PlayerViewAllActorImages,
	PlayerCreateActorImage,
}

var AllPlayerByName map[string]Player = map[string]Player{
	PlayerGrantAllPermissionsName:         PlayerGrantAllPermissions,
	PlayerRevokeAllPermissionsName:        PlayerRevokeAllPermissions,
	PlayerReviewCharacterApplicationsName: PlayerReviewCharacterApplications,
	PlayerViewAllRoomsName:                PlayerViewAllRooms,
	PlayerCreateRoomName:                  PlayerCreateRoom,
	PlayerViewAllActorImagesName:          PlayerViewAllActorImages,
	PlayerCreateActorImageName:            PlayerCreateActorImage,
}

var AllPlayerByTag map[string]Player = map[string]Player{
	PlayerGrantAllPermissionsTag:         PlayerGrantAllPermissions,
	PlayerRevokeAllPermissionsTag:        PlayerRevokeAllPermissions,
	PlayerReviewCharacterApplicationsTag: PlayerReviewCharacterApplications,
	PlayerViewAllRoomsTag:                PlayerViewAllRooms,
	PlayerCreateRoomTag:                  PlayerCreateRoom,
	PlayerViewAllActorImagesTag:          PlayerViewAllActorImages,
	PlayerCreateActorImageTag:            PlayerCreateActorImage,
}

var RootPlayerByName map[string]Player = map[string]Player{
	PlayerGrantAllPermissionsName:  PlayerGrantAllPermissions,
	PlayerRevokeAllPermissionsName: PlayerRevokeAllPermissions,
}
