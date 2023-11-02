package username

import (
	"os"
	"testing"

	redis "github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/require"
)

func TestKey(t *testing.T) {
	key := Key(69)
	require.Equal(t, key, "un:69")
}

func TestCache(t *testing.T) {
	r := redis.NewClient(&redis.Options{
		Addr:     os.Getenv("REDIS_ADDR"),
		Password: "",
		DB:       0,
		Protocol: 3,
	})
	defer r.Close()

	var pid int64 = 69

	err := Cache(r, pid, "nice")
	if err != nil {
		t.Fatal(err)
	}

	u, err := Get(r, pid)
	if err != nil {
		t.Fatal(err)
	}

	require.Equal(t, u, "nice")
}
