package player

import (
	"testing"

	"github.com/stretchr/testify/require"

	"petrichormud.com/app/internal/query"
)

func TestAllHardCodedPermissionsHaveValidNames(t *testing.T) {
	for _, permission := range AllPermissions {
		require.True(t, IsValidPermissionName(permission.Name))
	}
}

func TestIsValidNameReturnsFalseForInvalidName(t *testing.T) {
	require.False(t, IsValidPermissionName("not-a-permission"))
}

func NewPermissionsBuildsCorrectPermissions(t *testing.T) {
	pid := int64(1)
	permissionrecords := []query.PlayerPermission{
		{
			PID:  pid,
			Name: PermissionGrantAll.Name,
		},
	}

	permissions := NewPermissions(pid, permissionrecords)

	require.True(t, permissions.HasPermission(PermissionGrantAll.Name))
	require.False(t, permissions.HasPermission(PermissionRevokeAll.Name))
}
