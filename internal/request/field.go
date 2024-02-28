package request

import (
	"regexp"

	"petrichormud.com/app/internal/query"
)

type Field struct {
	Name        string
	Label       string
	Description string
	View        string
	Layout      string
	Updater     FieldUpdater
	Regexes     []*regexp.Regexp
	MinLen      int
	MaxLen      int
}

type FieldUpdater interface {
	Update(q *query.Queries, p UpdateFieldParams) error
}

func (f *Field) IsValueValid(v string) bool {
	if len(v) < f.MinLen {
		return false
	}
	if len(v) > f.MaxLen {
		return false
	}
	for _, regex := range f.Regexes {
		if regex.MatchString(v) {
			return false
		}
	}
	return true
}

func (f *Field) Update(q *query.Queries, p UpdateFieldParams) error {
	if !f.IsValueValid(p.Value) {
		return ErrInvalidInput
	}
	return f.Updater.Update(q, p)
}
