package permission

import "petrichormud.com/app/internal/queries"

func IsValid(perm string) bool {
	_, ok := AllPlayerByName[perm]
	return ok
}

func filterInvalidPlayerPermissions(perms []queries.PlayerPermission) []queries.PlayerPermission {
	result := []queries.PlayerPermission{}
	for _, perm := range perms {
		// TODO: Turn this into a map check or a  list
		if perm.Permission == PlayerGrantAllPermissionsName {
			result = append(result, perm)
		}
		if perm.Permission == PlayerReviewCharacterApplicationsName {
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
