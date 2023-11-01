package email

import (
	"context"
	"fmt"
	"os"

	"github.com/google/uuid"
	redis "github.com/redis/go-redis/v9"
	resend "github.com/resendlabs/resend-go"
)

const ThirtyMinutesInNanoseconds = 30 * 60 * 1000 * 1000 * 1000

func Verify(r *redis.Client, id int64, email string) error {
	key := uuid.NewString()

	err := Cache(r, key, id)
	if err != nil {
		return err
	}
	_, err = SendEmail(key, email)
	if err != nil {
		return err
	}

	return nil
}

func SendEmail(key string, email string) (resend.SendEmailResponse, error) {
	client := resend.NewClient(os.Getenv("RESEND_API_KEY"))
	base := os.Getenv("BASE_URL")
	url := fmt.Sprintf("%s/verify?t=%s", base, key)
	params := &resend.SendEmailRequest{
		To:      []string{email},
		From:    "verify@petrichormud.com",
		Html:    fmt.Sprintf("Welcome to PetrichorMUD! Please <a href=%q>click here</a> to verify your email address.", url),
		Subject: fmt.Sprintf("[PetrichorMUD] Verify %s", email),
		ReplyTo: "support@petrichormud.com",
	}
	return client.Emails.Send(params)
}

func Cache(r *redis.Client, key string, id int64) error {
	return r.Set(context.Background(), key, id, ThirtyMinutesInNanoseconds).Err()
}
