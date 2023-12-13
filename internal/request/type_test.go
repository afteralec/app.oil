package request

import (
	"context"
	"fmt"
	"testing"

	_ "github.com/go-sql-driver/mysql"
	"github.com/stretchr/testify/require"

	"petrichormud.com/app/internal/queries"
	"petrichormud.com/app/internal/shared"
)

func TestIsTypeDBTrue(t *testing.T) {
	i := shared.SetupInterfaces()
	defer i.Close()

	r, err := i.Queries.CreateRequest(context.Background(), queries.CreateRequestParams{
		Type: TypeCharacterApplication,
		PID:  1,
	})
	if err != nil {
		t.Fatal(err)
	}
	rid, err := r.LastInsertId()
	if err != nil {
		t.Fatal(err)
	}

	is, err := IsTypeDB(i.Database, TypeCharacterApplication, rid)
	if err != nil {
		t.Fatal(err)
	}
	require.Equal(t, true, is)

	query := fmt.Sprintf("DELETE FROM requests WHERE id = %d", rid)
	_, err = i.Database.Exec(query)
	if err != nil {
		t.Fatal(err)
	}
}

func TestIsTypeTxTrue(t *testing.T) {
	i := shared.SetupInterfaces()
	defer i.Close()

	tx, err := i.Database.Begin()
	defer tx.Rollback()
	if err != nil {
		t.Fatal(err)
	}
	qtx := i.Queries.WithTx(tx)

	r, err := qtx.CreateRequest(context.Background(), queries.CreateRequestParams{
		Type: TypeCharacterApplication,
		PID:  1,
	})
	if err != nil {
		t.Fatal(err)
	}
	rid, err := r.LastInsertId()
	if err != nil {
		t.Fatal(err)
	}

	is, err := IsTypeTx(tx, TypeCharacterApplication, rid)
	if err != nil {
		t.Fatal(err)
	}
	require.Equal(t, true, is)

	query := fmt.Sprintf("DELETE FROM requests WHERE id = %d", rid)
	_, err = tx.Exec(query)
	if err != nil {
		t.Fatal(err)
	}
	if err = tx.Commit(); err != nil {
		t.Fatal(err)
	}
}
