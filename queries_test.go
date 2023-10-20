package main

import (
	"context"
	"database/sql"
	"os"
	"testing"

	"ariga.io/atlas-go-sdk/atlasexec"
	_ "github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
	"github.com/stretchr/testify/require"

	"oil/app/odb"
)

func setupDatabase(t *testing.T) {
	db, err := sql.Open("mysql", "root:pass@/test")
	if err != nil {
		t.Fatal(err)
	}
	db.Exec("DROP DATABASE test;")
	db.Exec("CREATE DATABASE test;")
	db.Close()
}

func migrateDatabase(t *testing.T) {
	wd, err := os.Getwd()
	if err != nil {
		t.Fatalf("pwd failed with err: %v", err)
	}
	client, err := atlasexec.NewClient(wd, "atlas")
	if err != nil {
		t.Fatalf("Failed to create client with err: %v", err)
	}
	p := atlasexec.ApplyParams{}
	p.DirURL = "file://migrations"
	p.URL = "mysql://root:pass@127.0.0.1:3306/test"
	report, err := client.Apply(context.Background(), &p)
	if err != nil {
		t.Fatalf("Apply failed with error: %v", err)
	}
	_ = report
}

func TestPlayers(t *testing.T) {
	err := godotenv.Load()
	if err != nil {
		t.Fatal(err)
	}
	if os.Getenv("CLEAN_DATABASE") == "true" {
		setupDatabase(t)
		migrateDatabase(t)
	}
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
