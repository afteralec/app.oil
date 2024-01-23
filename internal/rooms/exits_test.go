package rooms

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"

	"petrichormud.com/app/internal/shared"
	"petrichormud.com/app/internal/tests"
)

func TestLoadExitRooms(t *testing.T) {
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

	exitRooms, err := LoadExitRooms(i.Queries, &room)
	if err != nil {
		t.Fatal(err)
	}

	exitRoom, ok := exitRooms[DirectionNorth]
	require.True(t, ok)

	require.Equal(t, exitRoom.ID, ridTwo)
}
