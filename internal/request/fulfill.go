package request

import (
	"petrichormud.com/app/internal/query"
	"petrichormud.com/app/internal/request/definition"
)

type Fulfiller interface {
	For() string
	Fulfill(q *query.Queries, req *query.Request) error
}

var FulfillersByType map[string]Fulfiller = map[string]Fulfiller{
	TypeCharacterApplication: &definition.FulfillerCharacterApplication,
}

func Fulfill(q *query.Queries, req *query.Request, p *query.Player) error {
	fulfiller, ok := FulfillersByType[req.Type]
	if !ok {
		return ErrNoDefinition
	}

	if err := fulfiller.Fulfill(q, req); err != nil {
		return err
	}

	if err := UpdateStatus(q, UpdateStatusParams{
		PID:    req.PID,
		Status: StatusFulfilled,
	}); err != nil {
		return err
	}

	return nil
}
