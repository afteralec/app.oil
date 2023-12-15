package permission

import (
	"petrichormud.com/app/internal/queries"
)

type PlayerIssuedPermissions struct {
	Permissions     map[string]bool
	PermissionsList []string
	PID             int64
}

func MakePlayerPermissions(pid int64, perms []queries.PlayerPermission) PlayerIssuedPermissions {
	filtered := filterInvalidPlayerPermissions(perms)
	list := []string{}
	for _, perm := range filtered {
		list = append(list, perm.Permission)
	}
	return PlayerIssuedPermissions{
		PID:             pid,
		PermissionsList: list,
		Permissions:     makePermissionMap(filtered),
	}
}

func (p *PlayerIssuedPermissions) HasPermissionInSet(set []string) bool {
	for _, perm := range set {
		if p.Permissions[perm] {
			return true
		}
	}
	return false
}
