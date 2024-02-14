package permission

import "petrichormud.com/app/internal/query"

// TODO: Rename PlayerGranted
// TODO: Test

type PlayerGranted struct {
	Permissions     map[string]bool
	PermissionsList []string
	PID             int64
}

func MakePlayerGranted(pid int64, perms []query.PlayerPermission) PlayerGranted {
	filtered := filterInvalidPlayerPermissions(perms)
	list := []string{}
	for _, perm := range filtered {
		list = append(list, perm.Permission)
	}
	return PlayerGranted{
		PID:             pid,
		PermissionsList: list,
		Permissions:     makePermissionMap(filtered),
	}
}

func (p *PlayerGranted) HasPermission(perm string) bool {
	_, ok := p.Permissions[perm]
	return ok
}

func (p *PlayerGranted) HasPermissionInSet(set []string) bool {
	for _, perm := range set {
		_, ok := p.Permissions[perm]
		if ok {
			return true
		}
	}
	return false
}

func (p *PlayerGranted) HasAllPermissionsInSet(set []string) bool {
	for _, perm := range set {
		_, ok := p.Permissions[perm]
		if !ok {
			return false
		}
	}
	return true
}

// TODO: This is to enable adding sub-permissions to grant individual or groups of permissions
func (p *PlayerGranted) CanGrantPermission(name string) bool {
	if !IsValidName(name) {
		return false
	}

	_, ok := RootPlayerByName[name]
	if ok {
		return false
	}

	_, ok = p.Permissions[PlayerGrantAllPermissionsName]
	return ok
}

// TODO: This is to enable adding sub-permissions to grant individual or groups of permissions
func (p *PlayerGranted) CanRevokePermission(name string) bool {
	if !IsValidName(name) {
		return false
	}

	_, ok := RootPlayerByName[name]
	if ok {
		return false
	}

	_, ok = p.Permissions[PlayerGrantAllPermissionsName]
	return ok
}
