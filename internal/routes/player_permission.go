package routes

import "fmt"

const PlayerPermissions string = "/player/permissions"

func PlayerPermissionsSearchPath() string {
	return fmt.Sprintf("%s/search", PlayerPermissions)
}

func PlayerPermissionsDetailPath(u string) string {
	return fmt.Sprintf("%s/%s", PlayerPermissions, u)
}

func PlayerPermissionsPath(id string) string {
	return fmt.Sprintf("%s/%s", PlayerPermissions, id)
}

func PlayerPermissionsTogglePath(id, tag string) string {
	return fmt.Sprintf("%s/%s/%s", PlayerPermissions, id, tag)
}
