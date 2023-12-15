package permission

import (
	"petrichormud.com/app/internal/queries"
)

type PlayerIssued struct {
	Permissions     map[string]bool
	PermissionsList []string
	PID             int64
}

func MakePlayerIssued(pid int64, perms []queries.PlayerPermission) PlayerIssued {
	filtered := filterInvalidPlayerPermissions(perms)
	list := []string{}
	for _, perm := range filtered {
		list = append(list, perm.Permission)
	}
	return PlayerIssued{
		PID:             pid,
		PermissionsList: list,
		Permissions:     makePermissionMap(filtered),
	}
}

func (p *PlayerIssued) HasPermissionInSet(set []string) bool {
	for _, perm := range set {
		if p.Permissions[perm] {
			return true
		}
	}
	return false
}
