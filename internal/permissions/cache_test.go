package permissions

import (
	"database/sql"
	"os"
	"testing"

	"github.com/gofiber/fiber/v2/middleware/session"
	redis "github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/require"
	"petrichormud.com/app/internal/shared"
)

func TestCache(t *testing.T) {
	var pid int64 = 69
	expected := []string{Login, AddEmail}

	// TODO: Peel this out into a setup function
	// TODO: Put this under test?
	db, err := sql.Open("mysql", os.Getenv("DATABASE_URL"))
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	r := redis.NewClient(&redis.Options{
		Addr:     os.Getenv("REDIS_ADDR"),
		Password: "",
		DB:       0,
		Protocol: 3,
	})
	defer r.Close()
	s := session.New()
	i := shared.InterfacesBuilder().Database(db).Redis(r).Sessions(s).Build()

	err = Cache(r, pid, expected)
	if err != nil {
		t.Fatal(err)
	}

	perms, err := List(&i, pid)
	if err != nil {
		t.Fatal(err)
	}

	require.Equal(t, expected, perms)
}
