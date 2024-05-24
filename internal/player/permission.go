package player

import "petrichormud.com/app/internal/query"

// TODO: Change this Name to Type
type Permission struct {
	Name  string
	Title string
	About string
}

var PermissionGrantAll Permission = Permission{
	Name:  "grant-all",
	Title: "Grant All Permissions",
	About: "The root permission. Only one person should have this at a time.",
}

var PermissionRevokeAll Permission = Permission{
	Name:  "revoke-all",
	Title: "Revoke All Permissions",
	About: "The root revocation permission. Only one person should have this at a time.",
}

var PermissionReviewCharacterApplications Permission = Permission{
	Name:  "review-character-applications",
	Title: "Review Character Applications",
	About: "Enable this player to review Character Applications.",
}

var PermissionViewAllRooms Permission = Permission{
	Name:  "view-all-rooms",
	Title: "View All Rooms",
	About: "The permission to view (but not edit) all room data.",
}

var PermissionCreateRoom Permission = Permission{
	Name:  "create-room",
	Title: "Create Room",
	About: "Create a new room, but not connect it to the grid.",
}

var PermissionViewAllActorImages Permission = Permission{
	Name:  "view-all-actor-images",
	Title: "View All Actor Images",
	About: "View all Actor Images, i.e. in the main Actor Images list.",
}

var PermissionCreateActorImage Permission = Permission{
	Name:  "create-actor-image",
	Title: "Create Actor Image",
	About: "Create new actor via creating new Actor Images",
}

var ShowPermissionViewPermissions []string = []string{
	PermissionGrantAll.Name,
}

var AllPermissions []Permission = []Permission{
	PermissionGrantAll,
	PermissionRevokeAll,
	PermissionReviewCharacterApplications,
	PermissionViewAllRooms,
	PermissionCreateRoom,
	PermissionViewAllActorImages,
	PermissionCreateActorImage,
}

var RootPermissions []Permission = []Permission{
	PermissionGrantAll,
	PermissionRevokeAll,
}

func permissionsByName(permissions []Permission) map[string]Permission {
	permissionsbyname := make(map[string]Permission)
	for _, permission := range permissions {
		permissionsbyname[permission.Name] = permission
	}
	return permissionsbyname
}

var (
	AllPermissionsByName  = permissionsByName(AllPermissions)
	RootPermissionsByName = permissionsByName(RootPermissions)
)

type Permissions struct {
	Permissions     map[string]bool
	PermissionsList []string
	PID             int64
}

func NewPermissions(pid int64, perms []query.PlayerPermission) Permissions {
	filtered := []query.PlayerPermission{}
	for _, perm := range perms {
		if IsValidPermissionName(perm.Name) {
			filtered = append(filtered, perm)
		}
	}
	list := []string{}
	for _, perm := range filtered {
		list = append(list, perm.Name)
	}
	permissionsmap := map[string]bool{}
	for _, perm := range filtered {
		permissionsmap[perm.Name] = true
	}
	return Permissions{
		PID:             pid,
		PermissionsList: list,
		Permissions:     permissionsmap,
	}
}

func (p *Permissions) HasPermission(perm string) bool {
	_, ok := p.Permissions[perm]
	return ok
}

func (p *Permissions) HasPermissionInSet(set []string) bool {
	for _, perm := range set {
		_, ok := p.Permissions[perm]
		if ok {
			return true
		}
	}
	return false
}

func (p *Permissions) HasAllPermissionsInSet(set []string) bool {
	for _, perm := range set {
		_, ok := p.Permissions[perm]
		if !ok {
			return false
		}
	}
	return true
}

// TODO: This is to enable adding sub-permissions to grant individual or groups of permissions
func (p *Permissions) CanGrantPermission(name string) bool {
	if !IsValidPermissionName(name) {
		return false
	}

	_, ok := RootPermissionsByName[name]
	if ok {
		return false
	}

	_, ok = p.Permissions[PermissionGrantAll.Name]
	return ok
}

// TODO: This is to enable adding sub-permissions to grant individual or groups of permissions
func (p *Permissions) CanRevokePermission(name string) bool {
	if !IsValidPermissionName(name) {
		return false
	}

	_, ok := RootPermissionsByName[name]
	if ok {
		return false
	}

	_, ok = p.Permissions[PermissionGrantAll.Name]
	return ok
}

func IsValidPermissionName(name string) bool {
	_, ok := AllPermissionsByName[name]
	return ok
}
