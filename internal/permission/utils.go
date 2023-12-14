package permission

import "petrichormud.com/app/internal/queries"

func filterInvalidPlayerPermissions(perms []queries.PlayerPermission) []queries.PlayerPermission {
	result := []queries.PlayerPermission{}
	for _, perm := range perms {
		if perm.Permission == PlayerAssignAllPermissions {
			result = append(result, perm)
		}
	}
	return result
}

func makePermissionMap(perms []queries.PlayerPermission) map[string]bool {
	filtered := filterInvalidPlayerPermissions(perms)
	result := map[string]bool{}
	for _, perm := range filtered {
		result[perm.Permission] = true
	}
	return result
}
