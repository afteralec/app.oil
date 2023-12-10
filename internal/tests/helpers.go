package tests

import (
	"context"
	"testing"

	"petrichormud.com/app/internal/queries"
	"petrichormud.com/app/internal/shared"
)

func ListEmailsForPlayer(t *testing.T, i *shared.Interfaces, username string) []queries.Email {
	p, err := i.Queries.GetPlayerByUsername(context.Background(), username)
	if err != nil {
		t.Fatal(err)
	}
	emails, err := i.Queries.ListEmails(context.Background(), p.ID)
	if err != nil {
		t.Fatal(err)
	}
	return emails
}
