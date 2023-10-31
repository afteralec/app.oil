package email

import (
	"petrichormud.com/app/internal/queries"
)

func Verified(emails []queries.PlayerEmail) []queries.PlayerEmail {
	verifiedEmails := []queries.PlayerEmail{}
	for i := range emails {
		if emails[i].Verified {
			verifiedEmails = append(verifiedEmails, emails[i])
		}
	}
	return verifiedEmails
}
