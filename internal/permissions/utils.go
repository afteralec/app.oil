package permissions

import "petrichormud.com/app/internal/queries"

func IsValidName(name string) bool {
	_, ok := AllPlayerByName[name]
	return ok
}

func filterInvalidPlayerPermissions(perms []queries.PlayerPermission) []queries.PlayerPermission {
	result := []queries.PlayerPermission{}
	for _, perm := range perms {
		if IsValidName(perm.Permission) {
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
