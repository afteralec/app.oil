package room

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"

	"petrichormud.com/app/internal/interfaces"
	"petrichormud.com/app/internal/test"
)

func TestLinkTwoWay(t *testing.T) {
	i := interfaces.SetupShared()
	defer i.Close()

	ridOne := test.CreateTestRoom(t, &i, test.TestRoom)
	defer test.DeleteTestRoom(t, &i, ridOne)
	ridTwo := test.CreateTestRoom(t, &i, test.TestRoom)
	defer test.DeleteTestRoom(t, &i, ridTwo)

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

		rmOne, err := qtx.GetRoom(context.Background(), ridOne)
		if err != nil {
			t.Fatal(err)
		}

		rmTwo, err := qtx.GetRoom(context.Background(), ridTwo)
		if err != nil {
			t.Fatal(err)
		}

		if err := tx.Commit(); err != nil {
			t.Fatal(err)
		}

		require.Equal(t, ridTwo, ExitID(&rmOne, dir), ridTwo)
		require.Equal(t, ridOne, ExitID(&rmTwo, DirectionOpposite(dir)))
	}
}

func TestLinkOneWay(t *testing.T) {
	i := interfaces.SetupShared()
	defer i.Close()

	ridOne := test.CreateTestRoom(t, &i, test.TestRoom)
	defer test.DeleteTestRoom(t, &i, ridOne)
	ridTwo := test.CreateTestRoom(t, &i, test.TestRoom)
	defer test.DeleteTestRoom(t, &i, ridTwo)

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

		rmOne, err := qtx.GetRoom(context.Background(), ridOne)
		if err != nil {
			t.Fatal(err)
		}

		rmTwo, err := qtx.GetRoom(context.Background(), ridTwo)
		if err != nil {
			t.Fatal(err)
		}

		if err := tx.Commit(); err != nil {
			t.Fatal(err)
		}

		require.Equal(t, ridTwo, ExitID(&rmOne, dir), ridTwo)
		require.Equal(t, int64(0), ExitID(&rmTwo, DirectionOpposite(dir)))
	}
}

func TestUnlink(t *testing.T) {
	i := interfaces.SetupShared()
	defer i.Close()

	ridOne := test.CreateTestRoom(t, &i, test.TestRoom)
	defer test.DeleteTestRoom(t, &i, ridOne)
	ridTwo := test.CreateTestRoom(t, &i, test.TestRoom)
	defer test.DeleteTestRoom(t, &i, ridTwo)

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

		rmOne, err := qtx.GetRoom(context.Background(), ridOne)
		if err != nil {
			t.Fatal(err)
		}

		rmTwo, err := qtx.GetRoom(context.Background(), ridTwo)
		if err != nil {
			t.Fatal(err)
		}

		if err := tx.Commit(); err != nil {
			t.Fatal(err)
		}

		require.Equal(t, int64(0), ExitID(&rmOne, dir))
		require.Equal(t, ridOne, ExitID(&rmTwo, DirectionOpposite(dir)))
	}
}
