package request

import (
	"context"
	"encoding/json"
	"errors"

	"petrichormud.com/app/internal/queries"
)

// TODO: Find a way to get each content extractor into the definition of the request type
func GetContent(qtx *queries.Queries, req *queries.Request) (map[string]string, error) {
	var b []byte
	m := map[string]string{}

	switch req.Type {
	case TypeCharacterApplication:
		app, err := qtx.GetCharacterApplicationContentForRequest(context.Background(), req.ID)
		if err != nil {
			return m, err
		}

		b, err = json.Marshal(app)
		if err != nil {
			return m, err
		}
	default:
		return m, errors.New("invalid type")
	}

	if err := json.Unmarshal(b, &m); err != nil {
		return map[string]string{}, err
	}

	return m, nil
}

func GetNextIncompleteField(t string, content map[string]string) string {
	fields := FieldNamesByType[t]
	for _, field := range fields {
		value, ok := content[field]
		if !ok {
			continue
		}
		if len(value) == 0 {
			return field
		}
	}
	return ""
}
