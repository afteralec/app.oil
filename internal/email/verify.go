package email

import (
	"context"
	"fmt"
	"os"

	"github.com/google/uuid"
	redis "github.com/redis/go-redis/v9"
	resend "github.com/resendlabs/resend-go"
)

const ThirtyMinutesInSeconds = 30 * 60

func Verify(r *redis.Client, id int64, email string) error {
	key := uuid.NewString()
	Cache(r, key, id)
	_, err := SendEmail(key, email)
	if err != nil {
		return err
	}
	return nil
}

func SendEmail(key string, email string) (resend.SendEmailResponse, error) {
	client := resend.NewClient(os.Getenv("RESEND_API_KEY"))
	url := fmt.Sprintf("https://petrichormud.com/verify?t=%s", key)
	params := &resend.SendEmailRequest{
		To:      []string{email},
		From:    "verify@petrichormud.com",
		Text:    fmt.Sprintf("Welcome to PetrichorMUD! Please <a href=%q>click here</a> to verify your email address.", url),
		Subject: fmt.Sprintf("[PetrichorMUD] Verify %s", email),
		ReplyTo: "supprot@petrichormud.com",
	}
	return client.Emails.Send(params)
}

func Cache(r *redis.Client, key string, id int64) {
	r.Set(context.Background(), key, id, ThirtyMinutesInSeconds)
}
