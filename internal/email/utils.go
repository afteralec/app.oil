package email

import "petrichormud.com/app/internal/query"

func Verified(emails []query.Email) []query.Email {
	verifiedEmails := []query.Email{}
	for i := range emails {
		if emails[i].Verified {
			verifiedEmails = append(verifiedEmails, emails[i])
		}
	}
	return verifiedEmails
}
