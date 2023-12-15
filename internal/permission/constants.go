package permission

type Player struct {
	Name  string
	Tag   string
	Title string
	About string
}

const (
	PlayerAssignAllPermissionsName        string = "AssignAllPermissions"
	PlayerReviewCharacterApplicationsName string = "ReviewCharacterApplications"
)

const (
	PlayerAssignAllPermissionsTag        string = "assign-all"
	PlayerReviewCharacterApplicationsTag string = "review-character-applications"
)

const (
	PlayerAssignAllPermissionsTitle        string = "Assign All Permissions"
	PlayerReviewCharacterApplicationsTitle string = "Review Character Applications"
)

const (
	PlayerAssignAllPermissionsAbout        string = "The root permission. Only one person should have this at a time."
	PlayerReviewCharacterApplicationsAbout string = "Enable this player to review Character Applications."
)

var PlayerAssignAllPermissions Player = Player{
	Name:  PlayerAssignAllPermissionsName,
	Tag:   PlayerAssignAllPermissionsTag,
	Title: PlayerAssignAllPermissionsTitle,
	About: PlayerAssignAllPermissionsAbout,
}

var PlayerReviewCharacterApplications Player = Player{
	Name:  PlayerReviewCharacterApplicationsName,
	Tag:   PlayerReviewCharacterApplicationsTag,
	Title: PlayerReviewCharacterApplicationsTitle,
	About: PlayerAssignAllPermissionsAbout,
}

var ShowPermissionViewPermissions []string = []string{
	PlayerAssignAllPermissionsName,
}

var AllPlayer []Player = []Player{
	PlayerAssignAllPermissions,
	PlayerReviewCharacterApplications,
}

var AllPlayerByTag map[string]Player = map[string]Player{
	PlayerAssignAllPermissionsTag:        PlayerAssignAllPermissions,
	PlayerReviewCharacterApplicationsTag: PlayerReviewCharacterApplications,
}
