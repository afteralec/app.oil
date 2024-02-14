package username

import (
	"context"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
	"github.com/google/uuid"
	redis "github.com/redis/go-redis/v9"

	"petrichormud.com/app/internal/constant"
	"petrichormud.com/app/internal/interfaces"
)

const (
	ThirtyTwoHoursInNanoseconds = 32 * 60 * 60 * 1000 * 1000 * 1000
	FiveMinutesInNanoSeconds    = 5 * 60 * 1000 * 1000 * 1000
)

func Get(i *interfaces.Shared, pid int64) (string, error) {
	key := Key(pid)
	username, err := i.Redis.Get(context.Background(), key).Result()
	if err != nil {
		if err == redis.Nil {
			u, err := i.Queries.GetPlayerUsername(context.Background(), pid)
			if err != nil {
				return "", err
			}
			if err = Cache(i.Redis, pid, u); err != nil {
				return "", err
			}
			return u, nil
		}
		return "", err
	}
	err = i.Redis.Expire(context.Background(), key, ThirtyTwoHoursInNanoseconds).Err()
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
	return fmt.Sprintf("%s:%d", constant.UsernameTokenKey, pid)
}

func RecoverySuccessKey(id string) string {
	return fmt.Sprintf("%s:%s", constant.UsernameRecoverySuccessTokenKey, id)
}
