package request

import (
	"context"
	"errors"

	"petrichormud.com/app/internal/query"
)

// TODO: ReviewDialog needs consolidated and cleaned up here
type ReviewDialogData struct {
	Path     string
	Variable string
}

type NewParams struct {
	Type string
	PID  int64
}

func New(q *query.Queries, p NewParams) (int64, error) {
	if p.PID == 0 {
		// TODO: Discrete error
		return 0, errors.New("invalid PID")
	}

	if !IsTypeValid(p.Type) {
		return 0, ErrInvalidType
	}

	fields, ok := FieldsByType[p.Type]
	if !ok {
		return 0, ErrInvalidType
	}

	result, err := q.CreateRequest(context.Background(), query.CreateRequestParams{
		PID:  p.PID,
		Type: p.Type,
	})
	if err != nil {
		return 0, err
	}
	rid, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	for _, field := range fields.List() {
		if err := q.CreateRequestField(context.Background(), query.CreateRequestFieldParams{
			RID:  rid,
			Type: field.Type,
			// TODO: Store the Default Field Status in its own const
			Status: FieldStatusNotReviewed,
			Value:  "",
		}); err != nil {
			return 0, err
		}
	}

	return rid, nil
}
