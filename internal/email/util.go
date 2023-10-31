package email

type PlayerEmail struct {
	Email    string
	Verified bool
	ID       int64
}

func Verified(emails []PlayerEmail) []PlayerEmail {
	verifiedEmails := []PlayerEmail{}
	for i := range emails {
		if emails[i].Verified {
			verifiedEmails = append(verifiedEmails, emails[i])
		}
	}
	return verifiedEmails
}
