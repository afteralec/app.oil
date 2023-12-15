package permission

const (
	PlayerAssignAllPermissions        string = "AssignAllPermissions"
	PlayerReviewCharacterApplications string = "ReviewCharacterApplications"
)

const (
	PlayerAssignAllPermissionsTag        string = "assign-all"
	PlayerReviewCharacterApplicationsTag string = "review-character-applications"
)

const (
	PlayerAssignAllPermissionsAbout        string = "The root permission. Only one person should have this at a time."
	PlayerReviewCharacterApplicationsAbout string = "Enable this player to review Character Applications."
)

type PlayerPermission struct {
	Name  string
	Tag   string
	About string
}

var PlayerPermissions map[string]PlayerPermission = map[string]PlayerPermission{
	PlayerAssignAllPermissions: {
		Name:  PlayerAssignAllPermissions,
		Tag:   PlayerAssignAllPermissionsTag,
		About: "The root permission. Only one person should have this at a time.",
	},
	PlayerReviewCharacterApplications: {
		Name:  PlayerReviewCharacterApplications,
		Tag:   PlayerReviewCharacterApplicationsTag,
		About: "The permission to review Character Applications.",
	},
}

var PlayerPermissionsByTag map[string]PlayerPermission = map[string]PlayerPermission{}

var ShowPermissionViewPermissions []string = []string{
	PlayerAssignAllPermissions,
}

var AllPlayerPermissionDetails [][]string = [][]string{
	{PlayerAssignAllPermissions, "The root permission. Only one person should have this at a time."},
	{PlayerReviewCharacterApplications, "The permission to review Character Applications."},
}
