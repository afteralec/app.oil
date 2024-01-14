package username

import (
	"context"
	"os"
	"time"

	pb "petrichormud.com/app/internal/proto/sending"
	"petrichormud.com/app/internal/queries"
	"petrichormud.com/app/internal/shared"
)

func Recover(i *shared.Interfaces, e queries.Email) (string, error) {
	id, err := CacheRecoverySuccessEmail(i.Redis, e.Address)
	if err != nil {
		return "", err
	}

	if os.Getenv("DISABLE_RESEND") == "true" {
		return id, nil
	}

	u, err := i.Queries.GetPlayerUsernameById(context.Background(), e.PID)
	if err != nil {
		return "", err
	}

	sender := pb.NewSenderClient(i.ClientConn)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	_, err = sender.SendUsernameRecovery(ctx, &pb.SendUsernameRecoveryRequest{
		Email:    e.Address,
		Username: u,
	})
	if err != nil {
		return "", err
	}

	return id, nil
}
