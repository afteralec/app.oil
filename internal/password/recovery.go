package password

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/google/uuid"
	redis "github.com/redis/go-redis/v9"

	"petrichormud.com/app/internal/constant"
	"petrichormud.com/app/internal/interfaces"
	pb "petrichormud.com/app/internal/proto/sending"
)

const ThirtyMinutesInNanoseconds = 30 * 60 * 1000 * 1000 * 1000

func SetupRecovery(i *interfaces.Shared, pid int64, email string) error {
	id := uuid.NewString()
	key := RecoveryKey(id)

	err := Cache(i.Redis, key, pid)
	if err != nil {
		return err
	}

	if os.Getenv("DISABLE_SENDING_STONE") == "true" {
		return nil
	}

	sender := pb.NewSenderClient(i.ClientConn)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	base := os.Getenv("BASE_URL")
	url := fmt.Sprintf("%s/reset/password?t=%s", base, key)
	_, err = sender.SendPasswordRecovery(ctx, &pb.SendPasswordRecoveryRequest{
		Email: email,
		Link:  url,
	})
	if err != nil {
		return err
	}

	return nil
}

func RecoveryKey(key string) string {
	return fmt.Sprintf("%s:%s", constant.RecoverPasswordTokenKey, key)
}

func Cache(r *redis.Client, key string, pid int64) error {
	return r.Set(context.Background(), key, pid, ThirtyMinutesInNanoseconds).Err()
}

func SetupRecoverySuccess(i *interfaces.Shared, email string) (string, error) {
	id := uuid.NewString()
	key := RecoverySuccessKey(id)

	err := CacheEmail(i.Redis, key, email)
	if err != nil {
		return "", err
	}

	return id, nil
}

func CacheEmail(r *redis.Client, key, email string) error {
	return r.Set(context.Background(), key, email, ThirtyMinutesInNanoseconds).Err()
}

func RecoverySuccessKey(key string) string {
	return fmt.Sprintf("%s:%s", constant.RecoverPasswordSuccessTokenKey, key)
}
