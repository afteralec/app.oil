package rooms

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"

	"petrichormud.com/app/internal/shared"
	"petrichormud.com/app/internal/tests"
)

func TestIsExitTwoWayFalseUnlinked(t *testing.T) {
	i := shared.SetupInterfaces()
	defer i.Close()

	ridOne := tests.CreateTestRoom(t, &i, tests.TestRoom)
	ridTwo := tests.CreateTestRoom(t, &i, tests.TestRoom)

	room, err := i.Queries.GetRoom(context.Background(), ridOne)
	if err != nil {
		t.Fatal(err)
	}

	exitRoom, err := i.Queries.GetRoom(context.Background(), ridTwo)
	if err != nil {
		t.Fatal(err)
	}

	require.False(t, IsExitTwoWay(&room, &exitRoom, DirectionNorth))
}

func TestIsExitTwoWayFalseOneWay(t *testing.T) {
	i := shared.SetupInterfaces()
	defer i.Close()

	ridOne := tests.CreateTestRoom(t, &i, tests.TestRoom)
	ridTwo := tests.CreateTestRoom(t, &i, tests.TestRoom)

	if err := Link(LinkParams{
		Queries:   i.Queries,
		ID:        ridOne,
		To:        ridTwo,
		Direction: DirectionNorth,
		TwoWay:    false,
	}); err != nil {
		t.Fatal(err)
	}

	room, err := i.Queries.GetRoom(context.Background(), ridOne)
	if err != nil {
		t.Fatal(err)
	}

	exitRoom, err := i.Queries.GetRoom(context.Background(), ridTwo)
	if err != nil {
		t.Fatal(err)
	}

	require.False(t, IsExitTwoWay(&room, &exitRoom, DirectionNorth))
}

func TestIsExitTwoWayFalseOneWayOpposite(t *testing.T) {
	i := shared.SetupInterfaces()
	defer i.Close()

	ridOne := tests.CreateTestRoom(t, &i, tests.TestRoom)
	ridTwo := tests.CreateTestRoom(t, &i, tests.TestRoom)

	if err := Link(LinkParams{
		Queries:   i.Queries,
		ID:        ridTwo,
		To:        ridOne,
		Direction: DirectionSouth,
		TwoWay:    false,
	}); err != nil {
		t.Fatal(err)
	}

	room, err := i.Queries.GetRoom(context.Background(), ridOne)
	if err != nil {
		t.Fatal(err)
	}

	exitRoom, err := i.Queries.GetRoom(context.Background(), ridTwo)
	if err != nil {
		t.Fatal(err)
	}

	require.False(t, IsExitTwoWay(&room, &exitRoom, DirectionNorth))
}

func TestIsExitTwoWayTrue(t *testing.T) {
	i := shared.SetupInterfaces()
	defer i.Close()

	ridOne := tests.CreateTestRoom(t, &i, tests.TestRoom)
	ridTwo := tests.CreateTestRoom(t, &i, tests.TestRoom)

	if err := Link(LinkParams{
		Queries:   i.Queries,
		ID:        ridOne,
		To:        ridTwo,
		Direction: DirectionNorth,
		TwoWay:    true,
	}); err != nil {
		t.Fatal(err)
	}

	room, err := i.Queries.GetRoom(context.Background(), ridOne)
	if err != nil {
		t.Fatal(err)
	}

	exitRoom, err := i.Queries.GetRoom(context.Background(), ridTwo)
	if err != nil {
		t.Fatal(err)
	}

	require.True(t, IsExitTwoWay(&room, &exitRoom, DirectionNorth))
}
