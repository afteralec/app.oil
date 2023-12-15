package permission

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
)

const (
	PlayerGrantAllPermissionsTag         string = "grant-all"
	PlayerRevokeAllPermissionsTag        string = "revoke-all"
	PlayerReviewCharacterApplicationsTag string = "review-character-applications"
)

const (
	PlayerGrantAllPermissionsTitle         string = "Grant All Permissions"
	PlayerRevokeAllPermissionsTitle        string = "Revoke All Permissions"
	PlayerReviewCharacterApplicationsTitle string = "Review Character Applications"
)

const (
	PlayerGrantAllPermissionsAbout         string = "The root permission. Only one person should have this at a time."
	PlayerRevokeAllPermissionsAbout        string = "The root revocation permission. Only one person should have this at a time."
	PlayerReviewCharacterApplicationsAbout string = "Enable this player to review Character Applications."
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

var ShowPermissionViewPermissions []string = []string{
	PlayerGrantAllPermissionsName,
}

var AllPlayer []Player = []Player{
	PlayerGrantAllPermissions,
	PlayerRevokeAllPermissions,
	PlayerReviewCharacterApplications,
}

var AllPlayerByName map[string]Player = map[string]Player{
	PlayerGrantAllPermissionsName:         PlayerGrantAllPermissions,
	PlayerRevokeAllPermissionsName:        PlayerRevokeAllPermissions,
	PlayerReviewCharacterApplicationsName: PlayerReviewCharacterApplications,
}

var AllPlayerByTag map[string]Player = map[string]Player{
	PlayerGrantAllPermissionsTag:         PlayerGrantAllPermissions,
	PlayerRevokeAllPermissionsTag:        PlayerRevokeAllPermissions,
	PlayerReviewCharacterApplicationsTag: PlayerReviewCharacterApplications,
}

var RootPlayerByName map[string]Player = map[string]Player{
	PlayerGrantAllPermissionsName:  PlayerGrantAllPermissions,
	PlayerRevokeAllPermissionsName: PlayerRevokeAllPermissions,
}
