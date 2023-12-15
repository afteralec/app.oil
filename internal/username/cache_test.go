package username

import (
	"testing"

	"github.com/stretchr/testify/require"
	"petrichormud.com/app/internal/shared"
)

func TestKey(t *testing.T) {
	key := Key(69)
	require.Equal(t, key, "un:69")
}

func TestCache(t *testing.T) {
	i := shared.SetupInterfaces()
	defer i.Close()

	var pid int64 = 69

	err := Cache(i.Redis, pid, "nice")
	if err != nil {
		t.Fatal(err)
	}

	u, err := Get(&i, pid)
	if err != nil {
		t.Fatal(err)
	}

	require.Equal(t, u, "nice")
}
