package routes

import "fmt"

const PlayerPermissions string = "/player/permissions"

func PlayerPermissionsSearchPath() string {
	return fmt.Sprintf("%s/search", PlayerPermissions)
}

func PlayerPermissionsDetailPath(u string) string {
	return fmt.Sprintf("%s/%s", PlayerPermissions, u)
}
