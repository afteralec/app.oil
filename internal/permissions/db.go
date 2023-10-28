package permissions

import (
	"context"
	"fmt"
	"strings"

	redis "github.com/redis/go-redis/v9"

	"petrichormud.com/app/internal/queries"
)

const TwoHoursInSeconds = 120 * 60

func List(q *queries.Queries, r *redis.Client, pid int64) ([]string, error) {
	key := Key(pid)
	exists, err := r.Exists(context.Background(), key).Result()
	if err != nil {
		return nil, err
	}

	if exists == 1 {
		perms, err := r.SMembers(context.Background(), key).Result()
		if err != nil {
			return nil, err
		}
		return perms, nil
	} else {
		records, err := q.ListPlayerPermissions(context.Background(), pid)
		if err != nil {
			return nil, err
		}

		var perms []string
		for i := 0; i < len(records); i++ {
			record := records[i]
			perms = append(perms, record.Permission)
		}

		r.SAdd(context.Background(), key, strings.Join(perms, " "), TwoHoursInSeconds)
		return perms, nil
	}
}

func Key(pid int64) string {
	return fmt.Sprintf("perm:%v", pid)
}
