package permission

import (
	"errors"

	"petrichormud.com/app/internal/queries"
)

type PlayerPermissions struct {
	Permissions     map[string]bool
	PermissionsList []string
	PID             int64
}

func MakePlayerPermissions(perms []queries.PlayerPermission) (PlayerPermissions, error) {
	if len(perms) == 0 {
		return PlayerPermissions{}, errors.New("cannot make PlayerPermissions from empty list")
	}
	filtered := filterInvalidPlayerPermissions(perms)
	if len(filtered) == 0 {
		return PlayerPermissions{}, errors.New("cannot make PlayerPermissions from empty list")
	}
	pid := filtered[0].PID
	list := []string{}
	for _, perm := range filtered {
		list = append(list, perm.Permission)
	}
	return PlayerPermissions{
		PID:             pid,
		PermissionsList: list,
		Permissions:     makePermissionMap(filtered),
	}, nil
}

func (p *PlayerPermissions) HasPermissionInSet(set []string) bool {
	for _, perm := range set {
		if p.Permissions[perm] {
			return true
		}
	}
	return false
}
