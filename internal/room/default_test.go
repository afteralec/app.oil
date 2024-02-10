package room

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestDefault(t *testing.T) {
	require.True(t, IsTitleValid(DefaultTitle))
	require.True(t, IsDescriptionValid(DefaultDescription))
	require.True(t, IsSizeValid(DefaultSize))
}
