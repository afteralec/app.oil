package integrationtest

import (
	"context"
	"database/sql"
	"os"
	"testing"

	_ "github.com/go-sql-driver/mysql"
	"github.com/stretchr/testify/require"

	"petrichormud.com/app/internal/queries"
)

func TestPlayers(t *testing.T) {
	if os.Getenv("RUN_INTEGRATION_TESTS") != "true" {
		t.Skip("Skipping integration test; set RUN_INTEGRATION_TESTS to \"true\" to run.")
	}
	ctx := context.Background()
	username := "alec"
	pw := "test-pw-hash"

	db, err := sql.Open("mysql", os.Getenv("DATABASE_URL"))
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	q := queries.New(db)

	result, err := q.CreatePlayer(ctx, queries.CreatePlayerParams{
		Username: username,
		PwHash:   pw,
	})
	if err != nil {
		t.Fatal(err)
	}

	playerId, err := result.LastInsertId()
	if err != nil {
		t.Fatal(err)
	}

	player, err := q.GetPlayer(ctx, playerId)
	if err != nil {
		t.Fatal(err)
	}

	require.Equal(t, username, player.Username)
	require.Equal(t, pw, player.PwHash)

	player, err = q.GetPlayerByUsername(ctx, username)
	if err != nil {
		t.Fatal(err)
	}

	require.Equal(t, playerId, player.ID)
	require.Equal(t, pw, player.PwHash)
}
