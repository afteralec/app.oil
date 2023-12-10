package password

import (
	"context"
	"fmt"
	"os"

	"github.com/google/uuid"
	redis "github.com/redis/go-redis/v9"
	resend "github.com/resendlabs/resend-go"
	"petrichormud.com/app/internal/shared"
)

const ThirtyMinutesInNanoseconds = 30 * 60 * 1000 * 1000 * 1000

func SetupRecovery(i *shared.Interfaces, pid int64, email string) error {
	id := uuid.NewString()
	key := RecoveryKey(id)

	err := Cache(i.Redis, key, pid)
	if err != nil {
		return err
	}

	if os.Getenv("DISABLE_RESEND") == "true" {
		return nil
	}
	_, err = SendRecoveryEmail(i, id, email)
	if err != nil {
		return err
	}

	return nil
}

func RecoveryKey(id string) string {
	return fmt.Sprintf("%s:%s", shared.RecoverPasswordTokenKey, id)
}

func SendRecoveryEmail(i *shared.Interfaces, key string, email string) (resend.SendEmailResponse, error) {
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
	return i.Resend.Emails.Send(params)
}

func Cache(r *redis.Client, key string, pid int64) error {
	return r.Set(context.Background(), key, pid, ThirtyMinutesInNanoseconds).Err()
}
