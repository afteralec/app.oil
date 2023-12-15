package permission

type Player struct {
	Name  string
	Tag   string
	Title string
	About string
}

const (
	PlayerAssignAllPermissionsName        string = "AssignAllPermissions"
	PlayerRevokeAllPermissionsName        string = "RevokeAllPermissions"
	PlayerReviewCharacterApplicationsName string = "ReviewCharacterApplications"
)

const (
	PlayerAssignAllPermissionsTag        string = "assign-all"
	PlayerRevokeAllPermissionsTag        string = "revoke-all"
	PlayerReviewCharacterApplicationsTag string = "review-character-applications"
)

const (
	PlayerAssignAllPermissionsTitle        string = "Assign All Permissions"
	PlayerRevokeAllPermissionsTitle        string = "Revoke All Permissions"
	PlayerReviewCharacterApplicationsTitle string = "Review Character Applications"
)

const (
	PlayerAssignAllPermissionsAbout        string = "The root permission. Only one person should have this at a time."
	PlayerRevokeAllPermissionsAbout        string = "The root revocation permission. Only one person should have this at a time."
	PlayerReviewCharacterApplicationsAbout string = "Enable this player to review Character Applications."
)

var PlayerAssignAllPermissions Player = Player{
	Name:  PlayerAssignAllPermissionsName,
	Tag:   PlayerAssignAllPermissionsTag,
	Title: PlayerAssignAllPermissionsTitle,
	About: PlayerAssignAllPermissionsAbout,
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
	PlayerAssignAllPermissionsName,
}

var AllPlayer []Player = []Player{
	PlayerAssignAllPermissions,
	PlayerRevokeAllPermissions,
	PlayerReviewCharacterApplications,
}

var AllPlayerByName map[string]Player = map[string]Player{
	PlayerAssignAllPermissionsName:        PlayerAssignAllPermissions,
	PlayerRevokeAllPermissionsName:        PlayerRevokeAllPermissions,
	PlayerReviewCharacterApplicationsName: PlayerReviewCharacterApplications,
}

var AllPlayerByTag map[string]Player = map[string]Player{
	PlayerAssignAllPermissionsTag:        PlayerAssignAllPermissions,
	PlayerRevokeAllPermissionsTag:        PlayerRevokeAllPermissions,
	PlayerReviewCharacterApplicationsTag: PlayerReviewCharacterApplications,
}

var RootPlayerByName map[string]Player = map[string]Player{
	PlayerAssignAllPermissionsName: PlayerAssignAllPermissions,
	PlayerRevokeAllPermissionsName: PlayerRevokeAllPermissions,
}
