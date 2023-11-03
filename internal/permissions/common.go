package permissions

import (
	"petrichormud.com/app/internal/queries"
)

const (
	Login           = "Login"
	AddEmail        = "AddEmail"
	ViewPermissions = "ViewPermissions"
)

func DefaultSet() [2]string {
	return [2]string{Login, AddEmail}
}

func AdminSet() [3]string {
	return [3]string{Login, AddEmail, ViewPermissions}
}

func MakeParams(p []string, pid int64) []queries.CreatePlayerPermissionsParams {
	var params []queries.CreatePlayerPermissionsParams
	for i := 0; i < len(p); i++ {
		permission := p[i]
		param := queries.CreatePlayerPermissionsParams{
			Pid:        pid,
			Permission: permission,
		}
		params = append(params, param)
	}
	return params
}
