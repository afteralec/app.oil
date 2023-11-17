package password

import (
	"context"
	"fmt"
	"os"

	"github.com/google/uuid"
	redis "github.com/redis/go-redis/v9"
	resend "github.com/resendlabs/resend-go"
)

const ThirtyMinutesInNanoseconds = 30 * 60 * 1000 * 1000 * 1000

func SetupRecovery(r *redis.Client, pid int64, email string) error {
	id := uuid.NewString()
	key := fmt.Sprintf("rp:%s", id)

	err := Cache(r, key, pid)
	if err != nil {
		return err
	}

	if os.Getenv("DISABLE_RESEND") == "true" {
		return nil
	}
	_, err = SendRecoveryEmail(key, email)
	if err != nil {
		return err
	}

	return nil
}

func SendRecoveryEmail(key string, email string) (resend.SendEmailResponse, error) {
	// TODO: Extract this so we aren't building a new client on each request
	client := resend.NewClient(os.Getenv("RESEND_API_KEY"))
	base := os.Getenv("BASE_URL")
	url := fmt.Sprintf("%s/reset/password?t=%s", base, key)
	params := &resend.SendEmailRequest{
		To:   []string{email},
		From: "recover@petrichormud.com",
		// TODO: Set up a doc page for what people should do if they didn't request a password recovery
		// TODO: Link to it here
		Html:    fmt.Sprintf("Hello! <a href=%q>Click here</a> to reset your password.", url),
		Subject: "[PetrichorMUD] Password Recovery",
		ReplyTo: "support@petrichormud.com",
	}
	return client.Emails.Send(params)
}

func Cache(r *redis.Client, key string, pid int64) error {
	return r.Set(context.Background(), key, pid, ThirtyMinutesInNanoseconds).Err()
}
