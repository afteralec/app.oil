package email

import (
	"testing"

	"github.com/stretchr/testify/require"
	"petrichormud.com/app/internal/queries"
)

func TestVerified(t *testing.T) {
	u := queries.Email{ID: 1, PID: 69, Address: "test@test.com", Verified: false}
	v := queries.Email{ID: 2, PID: 69, Address: "testagain@test.com", Verified: true}
	emails := []queries.Email{u, v}
	expected := []queries.Email{v}
	require.Equal(t, expected, Verified(emails))
}
