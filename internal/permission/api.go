package permission

import (
	"petrichormud.com/app/internal/queries"
)

type PlayerGranted struct {
	Permissions     map[string]bool
	PermissionsList []string
	PID             int64
}

func MakePlayerGranted(pid int64, perms []queries.PlayerPermission) PlayerGranted {
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

func (p *PlayerGranted) HasPermissionInSet(set []string) bool {
	for _, perm := range set {
		if p.Permissions[perm] {
			return true
		}
	}
	return false
}

func (p *PlayerGranted) CanGrantPermission(perm string) bool {
	if !IsValid(perm) {
		return false
	}

	_, ok := p.Permissions[PlayerAssignAllPermissionsName]
	return ok
}
