package request

import (
	"context"
	"testing"

	_ "github.com/go-sql-driver/mysql"
	"github.com/stretchr/testify/require"

	"petrichormud.com/app/internal/query"
	"petrichormud.com/app/internal/service"
)

// TODO: Get this in a shared helper
func DeleteTestRequest(t *testing.T, i *service.Interfaces, rid int64) {
	_, err := i.Database.Exec("DELETE FROM requests WHERE id = ?;", rid)
	if err != nil {
		t.Fatal(err)
	}
}

func TestIsEditablePlayerFieldFalseNotPlayer(t *testing.T) {
	i := service.NewInterfaces()

	rid, err := New(i.Queries, NewParams{
		PID:  3,
		Type: TypeCharacterApplication,
	})
	if err != nil {
		t.Fatal(err)
	}
	defer DeleteTestRequest(t, &i, rid)

	req, err := i.Queries.GetRequest(context.Background(), rid)
	if err != nil {
		t.Fatal(err)
	}

	fd, err := GetFieldDefinition(req.Type, "name")
	if err != nil {
		t.Fatal(err)
	}

	require.False(t, IsEditable(1, &req, fd))
}

func TestIsEditablePlayerFieldFalseReviewer(t *testing.T) {
	i := service.NewInterfaces()

	rid, err := New(i.Queries, NewParams{
		PID:  3,
		Type: TypeCharacterApplication,
	})
	if err != nil {
		t.Fatal(err)
	}
	defer DeleteTestRequest(t, &i, rid)

	if err := i.Queries.UpdateRequestStatus(context.Background(), query.UpdateRequestStatusParams{
		ID:     rid,
		Status: StatusInReview,
	}); err != nil {
		t.Fatal(err)
	}

	if err := i.Queries.UpdateRequestReviewer(context.Background(), query.UpdateRequestReviewerParams{
		ID:   rid,
		RPID: 4,
	}); err != nil {
		t.Fatal(err)
	}

	req, err := i.Queries.GetRequest(context.Background(), rid)
	if err != nil {
		t.Fatal(err)
	}

	fd, err := GetFieldDefinition(req.Type, "name")
	if err != nil {
		t.Fatal(err)
	}

	require.False(t, IsEditable(4, &req, fd))
}

func TestIsEditablePlayerFieldFalseWrongStatus(t *testing.T) {
	i := service.NewInterfaces()

	rid, err := New(i.Queries, NewParams{
		PID:  3,
		Type: TypeCharacterApplication,
	})
	if err != nil {
		t.Fatal(err)
	}
	defer DeleteTestRequest(t, &i, rid)

	if err := i.Queries.UpdateRequestStatus(context.Background(), query.UpdateRequestStatusParams{
		ID:     rid,
		Status: StatusSubmitted,
	}); err != nil {
		t.Fatal(err)
	}

	req, err := i.Queries.GetRequest(context.Background(), rid)
	if err != nil {
		t.Fatal(err)
	}

	fd, err := GetFieldDefinition(req.Type, "name")
	if err != nil {
		t.Fatal(err)
	}

	require.False(t, IsEditable(3, &req, fd))
}

func TestIsEditablePlayerFieldTrue(t *testing.T) {
	i := service.NewInterfaces()

	rid, err := New(i.Queries, NewParams{
		PID:  3,
		Type: TypeCharacterApplication,
	})
	if err != nil {
		t.Fatal(err)
	}
	defer DeleteTestRequest(t, &i, rid)

	req, err := i.Queries.GetRequest(context.Background(), rid)
	if err != nil {
		t.Fatal(err)
	}

	fd, err := GetFieldDefinition(req.Type, "name")
	if err != nil {
		t.Fatal(err)
	}

	require.True(t, IsEditable(3, &req, fd))
}

func TestIsEditableReviewerFieldFalseNotReviewer(t *testing.T) {
	i := service.NewInterfaces()

	rid, err := New(i.Queries, NewParams{
		PID:  3,
		Type: TypeCharacterApplication,
	})
	if err != nil {
		t.Fatal(err)
	}
	defer DeleteTestRequest(t, &i, rid)

	if err := i.Queries.UpdateRequestStatus(context.Background(), query.UpdateRequestStatusParams{
		ID:     rid,
		Status: StatusInReview,
	}); err != nil {
		t.Fatal(err)
	}

	if err := i.Queries.UpdateRequestReviewer(context.Background(), query.UpdateRequestReviewerParams{
		ID:   rid,
		RPID: 4,
	}); err != nil {
		t.Fatal(err)
	}

	req, err := i.Queries.GetRequest(context.Background(), rid)
	if err != nil {
		t.Fatal(err)
	}

	fd, err := GetFieldDefinition(req.Type, "keywords")
	if err != nil {
		t.Fatal(err)
	}

	require.False(t, IsEditable(3, &req, fd))
}

func TestIsEditableReviewerFieldFalseWrongStatus(t *testing.T) {
	i := service.NewInterfaces()

	rid, err := New(i.Queries, NewParams{
		PID:  3,
		Type: TypeCharacterApplication,
	})
	if err != nil {
		t.Fatal(err)
	}
	defer DeleteTestRequest(t, &i, rid)

	if err := i.Queries.UpdateRequestReviewer(context.Background(), query.UpdateRequestReviewerParams{
		ID:   rid,
		RPID: 4,
	}); err != nil {
		t.Fatal(err)
	}

	req, err := i.Queries.GetRequest(context.Background(), rid)
	if err != nil {
		t.Fatal(err)
	}

	fd, err := GetFieldDefinition(req.Type, "keywords")
	if err != nil {
		t.Fatal(err)
	}

	require.False(t, IsEditable(4, &req, fd))
}

func TestIsEditableReviewerFieldTrue(t *testing.T) {
	i := service.NewInterfaces()

	rid, err := New(i.Queries, NewParams{
		PID:  3,
		Type: TypeCharacterApplication,
	})
	if err != nil {
		t.Fatal(err)
	}
	defer DeleteTestRequest(t, &i, rid)

	if err := i.Queries.UpdateRequestStatus(context.Background(), query.UpdateRequestStatusParams{
		ID:     rid,
		Status: StatusInReview,
	}); err != nil {
		t.Fatal(err)
	}

	if err := i.Queries.UpdateRequestReviewer(context.Background(), query.UpdateRequestReviewerParams{
		ID:   rid,
		RPID: 4,
	}); err != nil {
		t.Fatal(err)
	}

	req, err := i.Queries.GetRequest(context.Background(), rid)
	if err != nil {
		t.Fatal(err)
	}

	fd, err := GetFieldDefinition(req.Type, "keywords")
	if err != nil {
		t.Fatal(err)
	}

	require.True(t, IsEditable(4, &req, fd))
}
