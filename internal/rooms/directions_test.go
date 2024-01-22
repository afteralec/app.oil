package rooms

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestIsDirectionValid(t *testing.T) {
	require.True(t, IsDirectionValid(DirectionNorth))
	require.True(t, IsDirectionValid(DirectionNortheast))
	require.True(t, IsDirectionValid(DirectionEast))
	require.True(t, IsDirectionValid(DirectionSoutheast))
	require.True(t, IsDirectionValid(DirectionSouth))
	require.True(t, IsDirectionValid(DirectionSouthwest))
	require.True(t, IsDirectionValid(DirectionWest))
	require.True(t, IsDirectionValid(DirectionNorthwest))
}

func TestIsDirectionValidInvalid(t *testing.T) {
	require.False(t, IsDirectionValid("weast"))
}
