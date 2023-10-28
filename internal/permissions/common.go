package permissions

import "petrichormud.com/app/internal/queries"

const (
	Login      = "Login"
	ViewPlayer = "ViewPlayer"
)

func DefaultSet() [1]string {
	return [1]string{Login}
}

func AdminSet() [2]string {
	return [2]string{Login, ViewPlayer}
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
