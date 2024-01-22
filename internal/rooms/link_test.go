package rooms

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
	"petrichormud.com/app/internal/shared"
	"petrichormud.com/app/internal/tests"
)

func TestLinkNorth(t *testing.T) {
	i := shared.SetupInterfaces()
	defer i.Close()

	ridOne := tests.CreateTestRoom(t, &i, tests.TestRoom)
	defer tests.DeleteTestRoom(t, &i, ridOne)
	ridTwo := tests.CreateTestRoom(t, &i, tests.TestRoom)
	defer tests.DeleteTestRoom(t, &i, ridTwo)

	tx, err := i.Database.Begin()
	if err != nil {
		t.Fatal(err)
	}
	qtx := i.Queries.WithTx(tx)

	Link(LinkParams{
		Queries:   qtx,
		ID:        ridOne,
		To:        ridTwo,
		Direction: DirectionNorth,
		TwoWay:    true,
	})

	roomOne, err := qtx.GetRoom(context.Background(), ridOne)
	if err != nil {
		t.Fatal(err)
	}

	roomTwo, err := qtx.GetRoom(context.Background(), ridTwo)
	if err != nil {
		t.Fatal(err)
	}

	if err := tx.Commit(); err != nil {
		t.Fatal(err)
	}

	require.Equal(t, roomOne.North, ridTwo)
	require.Equal(t, roomTwo.South, ridOne)
}

func TestLinkNortheast(t *testing.T) {
	i := shared.SetupInterfaces()
	defer i.Close()

	ridOne := tests.CreateTestRoom(t, &i, tests.TestRoom)
	defer tests.DeleteTestRoom(t, &i, ridOne)
	ridTwo := tests.CreateTestRoom(t, &i, tests.TestRoom)
	defer tests.DeleteTestRoom(t, &i, ridTwo)

	tx, err := i.Database.Begin()
	if err != nil {
		t.Fatal(err)
	}
	qtx := i.Queries.WithTx(tx)

	Link(LinkParams{
		Queries:   qtx,
		ID:        ridOne,
		To:        ridTwo,
		Direction: DirectionNortheast,
		TwoWay:    true,
	})

	roomOne, err := qtx.GetRoom(context.Background(), ridOne)
	if err != nil {
		t.Fatal(err)
	}

	roomTwo, err := qtx.GetRoom(context.Background(), ridTwo)
	if err != nil {
		t.Fatal(err)
	}

	if err := tx.Commit(); err != nil {
		t.Fatal(err)
	}

	require.Equal(t, roomOne.Northeast, ridTwo)
	require.Equal(t, roomTwo.Southwest, ridOne)
}

func TestLinkEast(t *testing.T) {
	i := shared.SetupInterfaces()
	defer i.Close()

	ridOne := tests.CreateTestRoom(t, &i, tests.TestRoom)
	defer tests.DeleteTestRoom(t, &i, ridOne)
	ridTwo := tests.CreateTestRoom(t, &i, tests.TestRoom)
	defer tests.DeleteTestRoom(t, &i, ridTwo)

	tx, err := i.Database.Begin()
	if err != nil {
		t.Fatal(err)
	}
	qtx := i.Queries.WithTx(tx)

	Link(LinkParams{
		Queries:   qtx,
		ID:        ridOne,
		To:        ridTwo,
		Direction: DirectionEast,
		TwoWay:    true,
	})

	roomOne, err := qtx.GetRoom(context.Background(), ridOne)
	if err != nil {
		t.Fatal(err)
	}

	roomTwo, err := qtx.GetRoom(context.Background(), ridTwo)
	if err != nil {
		t.Fatal(err)
	}

	if err := tx.Commit(); err != nil {
		t.Fatal(err)
	}

	require.Equal(t, roomOne.East, ridTwo)
	require.Equal(t, roomTwo.West, ridOne)
}

func TestLinkSoutheast(t *testing.T) {
	i := shared.SetupInterfaces()
	defer i.Close()

	ridOne := tests.CreateTestRoom(t, &i, tests.TestRoom)
	defer tests.DeleteTestRoom(t, &i, ridOne)
	ridTwo := tests.CreateTestRoom(t, &i, tests.TestRoom)
	defer tests.DeleteTestRoom(t, &i, ridTwo)

	tx, err := i.Database.Begin()
	if err != nil {
		t.Fatal(err)
	}
	qtx := i.Queries.WithTx(tx)

	Link(LinkParams{
		Queries:   qtx,
		ID:        ridOne,
		To:        ridTwo,
		Direction: DirectionSoutheast,
		TwoWay:    true,
	})

	roomOne, err := qtx.GetRoom(context.Background(), ridOne)
	if err != nil {
		t.Fatal(err)
	}

	roomTwo, err := qtx.GetRoom(context.Background(), ridTwo)
	if err != nil {
		t.Fatal(err)
	}

	if err := tx.Commit(); err != nil {
		t.Fatal(err)
	}

	require.Equal(t, roomOne.Southeast, ridTwo)
	require.Equal(t, roomTwo.Northwest, ridOne)
}

func TestLinkSouth(t *testing.T) {
	i := shared.SetupInterfaces()
	defer i.Close()

	ridOne := tests.CreateTestRoom(t, &i, tests.TestRoom)
	defer tests.DeleteTestRoom(t, &i, ridOne)
	ridTwo := tests.CreateTestRoom(t, &i, tests.TestRoom)
	defer tests.DeleteTestRoom(t, &i, ridTwo)

	tx, err := i.Database.Begin()
	if err != nil {
		t.Fatal(err)
	}
	qtx := i.Queries.WithTx(tx)

	Link(LinkParams{
		Queries:   qtx,
		ID:        ridOne,
		To:        ridTwo,
		Direction: DirectionSouth,
		TwoWay:    true,
	})

	roomOne, err := qtx.GetRoom(context.Background(), ridOne)
	if err != nil {
		t.Fatal(err)
	}

	roomTwo, err := qtx.GetRoom(context.Background(), ridTwo)
	if err != nil {
		t.Fatal(err)
	}

	if err := tx.Commit(); err != nil {
		t.Fatal(err)
	}

	require.Equal(t, roomOne.South, ridTwo)
	require.Equal(t, roomTwo.North, ridOne)
}

func TestLinkSouthwest(t *testing.T) {
	i := shared.SetupInterfaces()
	defer i.Close()

	ridOne := tests.CreateTestRoom(t, &i, tests.TestRoom)
	defer tests.DeleteTestRoom(t, &i, ridOne)
	ridTwo := tests.CreateTestRoom(t, &i, tests.TestRoom)
	defer tests.DeleteTestRoom(t, &i, ridTwo)

	tx, err := i.Database.Begin()
	if err != nil {
		t.Fatal(err)
	}
	qtx := i.Queries.WithTx(tx)

	Link(LinkParams{
		Queries:   qtx,
		ID:        ridOne,
		To:        ridTwo,
		Direction: DirectionSouthwest,
		TwoWay:    true,
	})

	roomOne, err := qtx.GetRoom(context.Background(), ridOne)
	if err != nil {
		t.Fatal(err)
	}

	roomTwo, err := qtx.GetRoom(context.Background(), ridTwo)
	if err != nil {
		t.Fatal(err)
	}

	if err := tx.Commit(); err != nil {
		t.Fatal(err)
	}

	require.Equal(t, roomOne.Southwest, ridTwo)
	require.Equal(t, roomTwo.Northeast, ridOne)
}

func TestLinkWest(t *testing.T) {
	i := shared.SetupInterfaces()
	defer i.Close()

	ridOne := tests.CreateTestRoom(t, &i, tests.TestRoom)
	defer tests.DeleteTestRoom(t, &i, ridOne)
	ridTwo := tests.CreateTestRoom(t, &i, tests.TestRoom)
	defer tests.DeleteTestRoom(t, &i, ridTwo)

	tx, err := i.Database.Begin()
	if err != nil {
		t.Fatal(err)
	}
	qtx := i.Queries.WithTx(tx)

	Link(LinkParams{
		Queries:   qtx,
		ID:        ridOne,
		To:        ridTwo,
		Direction: DirectionWest,
		TwoWay:    true,
	})

	roomOne, err := qtx.GetRoom(context.Background(), ridOne)
	if err != nil {
		t.Fatal(err)
	}

	roomTwo, err := qtx.GetRoom(context.Background(), ridTwo)
	if err != nil {
		t.Fatal(err)
	}

	if err := tx.Commit(); err != nil {
		t.Fatal(err)
	}

	require.Equal(t, roomOne.West, ridTwo)
	require.Equal(t, roomTwo.East, ridOne)
}

func TestLinkNorthwest(t *testing.T) {
	i := shared.SetupInterfaces()
	defer i.Close()

	ridOne := tests.CreateTestRoom(t, &i, tests.TestRoom)
	defer tests.DeleteTestRoom(t, &i, ridOne)
	ridTwo := tests.CreateTestRoom(t, &i, tests.TestRoom)
	defer tests.DeleteTestRoom(t, &i, ridTwo)

	tx, err := i.Database.Begin()
	if err != nil {
		t.Fatal(err)
	}
	qtx := i.Queries.WithTx(tx)

	Link(LinkParams{
		Queries:   qtx,
		ID:        ridOne,
		To:        ridTwo,
		Direction: DirectionNorthwest,
		TwoWay:    true,
	})

	roomOne, err := qtx.GetRoom(context.Background(), ridOne)
	if err != nil {
		t.Fatal(err)
	}

	roomTwo, err := qtx.GetRoom(context.Background(), ridTwo)
	if err != nil {
		t.Fatal(err)
	}

	if err := tx.Commit(); err != nil {
		t.Fatal(err)
	}

	require.Equal(t, roomOne.Northwest, ridTwo)
	require.Equal(t, roomTwo.Southeast, ridOne)
}
