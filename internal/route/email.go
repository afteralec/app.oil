package route

import "fmt"

const (
	Email       = "/player/email"
	VerifyEmail = "/verify"
)

func NewEmailPath() string {
	return fmt.Sprintf("%s/new", Email)
}

func EmailPath(id string) string {
	return fmt.Sprintf("%s/%s", Email, id)
}

func ResendEmailVerificationPath(id string) string {
	return fmt.Sprintf("%s/resend", EmailPath(id))
}

func VerifyEmailWithToken(t string) string {
	return fmt.Sprintf("%s?t=%s", VerifyEmail, t)
}
