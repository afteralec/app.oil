package username

import (
	"context"
	"os"
	"time"

	"petrichormud.com/app/internal/proto/sending"
	"petrichormud.com/app/internal/query"
	"petrichormud.com/app/internal/service"
)

func Recover(i *service.Interfaces, e query.Email) (string, error) {
	id, err := CacheRecoverySuccessEmail(i.Redis, e.Address)
	if err != nil {
		return "", err
	}

	u, err := i.Queries.GetPlayerUsernameById(context.Background(), e.PID)
	if err != nil {
		return "", err
	}

	if os.Getenv("DISABLE_SENDING_STONE") == "true" {
		return id, nil
	}

	sender := sending.NewSenderClient(i.ClientConn)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	_, err = sender.SendUsernameRecovery(ctx, &sending.SendUsernameRecoveryRequest{
		Email:    e.Address,
		Username: u,
	})
	if err != nil {
		return "", err
	}

	return id, nil
}
