package request

import (
	"testing"

	"github.com/stretchr/testify/require"
)

// TODO: MOAR

func TestIsCommentValidSuccess(t *testing.T) {
	require.True(t, IsCommentValid("This name is fantastic."))
}
