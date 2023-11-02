package email

import (
	"petrichormud.com/app/internal/queries"
)

func Verified(emails []queries.Email) []queries.Email {
	verifiedEmails := []queries.Email{}
	for i := range emails {
		if emails[i].Verified {
			verifiedEmails = append(verifiedEmails, emails[i])
		}
	}
	return verifiedEmails
}
