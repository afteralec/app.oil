package permissions

import (
	"testing"

	"github.com/stretchr/testify/require"
	"petrichormud.com/app/internal/queries"
)

func TestMakeParams(t *testing.T) {
	var pid int64 = 69

	expected := []queries.CreatePlayerPermissionsParams{
		{
			Pid:        pid,
			Permission: Login,
		},
		{
			Pid:        pid,
			Permission: AddEmail,
		},
	}

	p := []string{Login, AddEmail}

	params := MakeParams(p, pid)
	require.Equal(t, expected, params)
}
