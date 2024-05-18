package change

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestSanitizeTextWithoutInvalidChars(t *testing.T) {
	text := "Sorry, you can't use \"guy\" as a primary noun."
	require.Equal(t, text, SanitizeText(text))
}

func TestSanitizeTextWithInvalidChars(t *testing.T) {
	text := "test1234@%%"
	expected := "test"
	require.Equal(t, expected, SanitizeText(text))
}

func TestIsTextValidWithValid(t *testing.T) {
	text := "Sorry, you can't use \"guy\" as a primary noun."
	require.True(t, IsTextValid(text))
}

func TestIsTextValidWithInvalid(t *testing.T) {
	text := "test1234@%%"
	require.False(t, IsTextValid(text))
}
