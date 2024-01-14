package email

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/google/uuid"
	redis "github.com/redis/go-redis/v9"

	pb "petrichormud.com/app/internal/proto/sending"
	"petrichormud.com/app/internal/shared"
)

const ThirtyMinutesInNanoseconds = 30 * 60 * 1000 * 1000 * 1000

func SendVerificationEmail(i *shared.Interfaces, id int64, email string) error {
	token := uuid.NewString()
	key := VerificationKey(token)
	if err := Cache(i.Redis, key, id); err != nil {
		return err
	}

	if os.Getenv("DISABLE_RESEND") == "true" {
		return nil
	}

	base := os.Getenv("BASE_URL")
	url := fmt.Sprintf("%s/verify?t=%s", base, token)

	sender := pb.NewSenderClient(i.ClientConn)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	_, err := sender.SendEmailVerification(ctx, &pb.SendEmailVerificationRequest{
		Email: email,
		Link:  url,
	})
	if err != nil {
		return err
	}

	return nil
}

func VerificationKey(id string) string {
	return fmt.Sprintf("%s:%s", shared.VerifyEmailTokenKey, id)
}

func Cache(r *redis.Client, key string, id int64) error {
	return r.Set(context.Background(), key, id, ThirtyMinutesInNanoseconds).Err()
}
