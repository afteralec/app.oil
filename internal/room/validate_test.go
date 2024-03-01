package room

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestIsDescriptionValid(t *testing.T) {
	validDescription := "Dark, oiled wood encloses this cozy office, each panel polished to an immaculate sheen. In stark contrast, the floor is a pale, sanded expanse of knotted hardwood, with brightly-colored rugs waiting to soften footsteps. A sweeping vista sprawls beyond the floor-to-ceiling windows, its misty landscape dotted with jagged peaks."

	require.True(t, IsDescriptionValid(validDescription))
}
