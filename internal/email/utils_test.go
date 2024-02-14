package email

import (
	"testing"

	"github.com/stretchr/testify/require"

	"petrichormud.com/app/internal/query"
)

func TestVerified(t *testing.T) {
	u := query.Email{ID: 1, PID: 69, Address: "test@test.com", Verified: false}
	v := query.Email{ID: 2, PID: 69, Address: "testagain@test.com", Verified: true}
	emails := []query.Email{u, v}
	expected := []query.Email{v}
	require.Equal(t, expected, Verified(emails))
}
