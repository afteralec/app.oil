package test

import (
	"context"
	"database/sql"
	"os"
	"testing"

	_ "github.com/go-sql-driver/mysql"
	"github.com/stretchr/testify/require"

	"oil/app/odb"
)

func TestPlayers(t *testing.T) {
	ctx := context.Background()
	username := "alec"
	pw := "test-pw-hash"

	db, err := sql.Open("mysql", os.Getenv("DATABASE_URL"))
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	queries := odb.New(db)

	result, err := queries.CreatePlayer(ctx, odb.CreatePlayerParams{
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

	player, err := queries.GetPlayer(ctx, playerId)
	if err != nil {
		t.Fatal(err)
	}

	require.Equal(t, username, player.Username)
	require.Equal(t, pw, player.PwHash)

	player, err = queries.GetPlayerByUsername(ctx, username)
	if err != nil {
		t.Fatal(err)
	}

	require.Equal(t, playerId, player.ID)
	require.Equal(t, pw, player.PwHash)
}
