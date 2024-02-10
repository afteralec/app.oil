package room

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"

	"petrichormud.com/app/internal/interfaces"
	"petrichormud.com/app/internal/tests"
)

func TestIsExitTwoWayFalseUnlinked(t *testing.T) {
	i := interfaces.SetupShared()
	defer i.Close()

	ridOne := tests.CreateTestRoom(t, &i, tests.TestRoom)
	ridTwo := tests.CreateTestRoom(t, &i, tests.TestRoom)

	rm, err := i.Queries.GetRoom(context.Background(), ridOne)
	if err != nil {
		t.Fatal(err)
	}

	exitrm, err := i.Queries.GetRoom(context.Background(), ridTwo)
	if err != nil {
		t.Fatal(err)
	}

	require.False(t, IsExitTwoWay(&rm, &exitrm, DirectionNorth))
}

func TestIsExitTwoWayFalseOneWay(t *testing.T) {
	i := interfaces.SetupShared()
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

	rm, err := i.Queries.GetRoom(context.Background(), ridOne)
	if err != nil {
		t.Fatal(err)
	}

	exitrm, err := i.Queries.GetRoom(context.Background(), ridTwo)
	if err != nil {
		t.Fatal(err)
	}

	require.False(t, IsExitTwoWay(&rm, &exitrm, DirectionNorth))
}

func TestIsExitTwoWayFalseOneWayOpposite(t *testing.T) {
	i := interfaces.SetupShared()
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

	rm, err := i.Queries.GetRoom(context.Background(), ridOne)
	if err != nil {
		t.Fatal(err)
	}

	exitrm, err := i.Queries.GetRoom(context.Background(), ridTwo)
	if err != nil {
		t.Fatal(err)
	}

	require.False(t, IsExitTwoWay(&rm, &exitrm, DirectionNorth))
}

func TestIsExitTwoWayTrue(t *testing.T) {
	i := interfaces.SetupShared()
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

	rm, err := i.Queries.GetRoom(context.Background(), ridOne)
	if err != nil {
		t.Fatal(err)
	}

	exitrm, err := i.Queries.GetRoom(context.Background(), ridTwo)
	if err != nil {
		t.Fatal(err)
	}

	require.True(t, IsExitTwoWay(&rm, &exitrm, DirectionNorth))
}
