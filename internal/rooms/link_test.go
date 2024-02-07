package rooms

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"

	"petrichormud.com/app/internal/interfaces"
	"petrichormud.com/app/internal/tests"
)

func TestLinkTwoWay(t *testing.T) {
	i := interfaces.SetupShared()
	defer i.Close()

	ridOne := tests.CreateTestRoom(t, &i, tests.TestRoom)
	defer tests.DeleteTestRoom(t, &i, ridOne)
	ridTwo := tests.CreateTestRoom(t, &i, tests.TestRoom)
	defer tests.DeleteTestRoom(t, &i, ridTwo)

	for _, dir := range DirectionsList {
		tx, err := i.Database.Begin()
		if err != nil {
			t.Fatal(err)
		}
		defer tx.Rollback()
		qtx := i.Queries.WithTx(tx)

		Link(LinkParams{
			Queries:   qtx,
			ID:        ridOne,
			To:        ridTwo,
			Direction: dir,
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

		require.Equal(t, ridTwo, ExitID(&roomOne, dir), ridTwo)
		require.Equal(t, ridOne, ExitID(&roomTwo, DirectionOpposite(dir)))
	}
}

func TestLinkOneWay(t *testing.T) {
	i := interfaces.SetupShared()
	defer i.Close()

	ridOne := tests.CreateTestRoom(t, &i, tests.TestRoom)
	defer tests.DeleteTestRoom(t, &i, ridOne)
	ridTwo := tests.CreateTestRoom(t, &i, tests.TestRoom)
	defer tests.DeleteTestRoom(t, &i, ridTwo)

	for _, dir := range DirectionsList {
		tx, err := i.Database.Begin()
		if err != nil {
			t.Fatal(err)
		}
		defer tx.Rollback()
		qtx := i.Queries.WithTx(tx)

		Link(LinkParams{
			Queries:   qtx,
			ID:        ridOne,
			To:        ridTwo,
			Direction: dir,
			TwoWay:    false,
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

		require.Equal(t, ridTwo, ExitID(&roomOne, dir), ridTwo)
		require.Equal(t, int64(0), ExitID(&roomTwo, DirectionOpposite(dir)))
	}
}

func TestUnlink(t *testing.T) {
	i := interfaces.SetupShared()
	defer i.Close()

	ridOne := tests.CreateTestRoom(t, &i, tests.TestRoom)
	defer tests.DeleteTestRoom(t, &i, ridOne)
	ridTwo := tests.CreateTestRoom(t, &i, tests.TestRoom)
	defer tests.DeleteTestRoom(t, &i, ridTwo)

	for _, dir := range DirectionsList {
		tx, err := i.Database.Begin()
		if err != nil {
			t.Fatal(err)
		}
		defer tx.Rollback()
		qtx := i.Queries.WithTx(tx)

		Link(LinkParams{
			Queries:   qtx,
			ID:        ridOne,
			To:        ridTwo,
			Direction: dir,
			TwoWay:    true,
		})

		Unlink(UnlinkParams{
			Queries:   qtx,
			ID:        ridOne,
			Direction: dir,
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

		require.Equal(t, int64(0), ExitID(&roomOne, dir))
		require.Equal(t, ridOne, ExitID(&roomTwo, DirectionOpposite(dir)))
	}
}
