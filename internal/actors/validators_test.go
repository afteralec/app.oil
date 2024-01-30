package actors

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestIsShortDescriptionValid(t *testing.T) {
	require.True(t, IsShortDescriptionValid(DefaultImageShortDescription))
}

func TestIsDescriptionValid(t *testing.T) {
	require.True(t, IsDescriptionValid(DefaultImageDescription))
}
