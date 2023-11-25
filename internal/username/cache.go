package username

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	redis "github.com/redis/go-redis/v9"
)

const (
	ThirtyTwoHoursInNanoseconds = 32 * 60 * 60 * 1000 * 1000 * 1000
	FiveMinutesInNanoSeconds    = 5 * 60 * 1000 * 1000 * 1000
)

func Get(r *redis.Client, pid int64) (string, error) {
	key := Key(pid)
	username, err := r.Get(context.Background(), key).Result()
	if err != nil {
		return "", err
	}
	err = r.Expire(context.Background(), key, ThirtyTwoHoursInNanoseconds).Err()
	if err != nil {
		return "", err
	}

	return username, nil
}

func Cache(r *redis.Client, pid int64, username string) error {
	err := r.Set(context.Background(), Key(pid), username, ThirtyTwoHoursInNanoseconds).Err()
	if err != nil {
		return err
	}
	return nil
}

func CacheRecoverySuccessEmail(r *redis.Client, address string) (string, error) {
	id := uuid.NewString()
	key := RecoverySuccessKey(id)
	err := r.Set(context.Background(), key, address, FiveMinutesInNanoSeconds).Err()
	if err != nil {
		return "", err
	}
	return id, nil
}

func Key(pid int64) string {
	return fmt.Sprintf("un:%d", pid)
}

func RecoverySuccessKey(id string) string {
	return fmt.Sprintf("rus:%s", id)
}
