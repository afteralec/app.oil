package permissions

import (
	"context"
	"fmt"

	redis "github.com/redis/go-redis/v9"
	"petrichormud.com/app/internal/shared"
)

const TwoHoursInNanoseconds = 2 * 60 * 60 * 1000 * 1000 * 1000

func List(i *shared.Interfaces, pid int64) ([]string, error) {
	key := Key(pid)
	exists, err := i.Redis.Exists(context.Background(), key).Result()
	if err != nil {
		return nil, err
	}

	if exists == 1 {
		perms, err := i.Redis.SMembers(context.Background(), key).Result()
		if err != nil {
			return nil, err
		}
		return perms, nil
	} else {
		records, err := i.Queries.ListPlayerPermissions(context.Background(), pid)
		if err != nil {
			return nil, err
		}

		var perms []string
		for i := 0; i < len(records); i++ {
			record := records[i]
			perms = append(perms, record.Permission)
		}

		err = Cache(i.Redis, key, perms)
		if err != nil {
			return nil, err
		}
		return perms, nil
	}
}

func Cache(r *redis.Client, key string, perms []string) error {
	err := r.SAdd(context.Background(), key, perms).Err()
	if err != nil {
		return err
	}
	return r.Expire(context.Background(), key, TwoHoursInNanoseconds).Err()
}

func Key(pid int64) string {
	return fmt.Sprintf("perm:%v", pid)
}
