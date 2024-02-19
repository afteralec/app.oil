package player

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestAllHardCodedPermissionsHaveValidNames(t *testing.T) {
	for _, permission := range AllPermissions {
		require.True(t, IsValidPermissionName(permission.Name))
	}
}

func TestIsValidNameReturnsFalseForInvalidName(t *testing.T) {
	require.False(t, IsValidPermissionName("not-a-permission"))
}
