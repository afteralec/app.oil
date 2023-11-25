package username

import (
	"context"
	"fmt"
	"os"

	resend "github.com/resendlabs/resend-go"
	"petrichormud.com/app/internal/queries"
	"petrichormud.com/app/internal/shared"
)

func Recover(i *shared.Interfaces, e queries.Email) (string, error) {
	id, err := CacheRecoverySuccessEmail(i.Redis, e.Address)
	if err != nil {
		return "", err
	}

	username, err := i.Queries.GetPlayerUsernameById(context.Background(), e.Pid)
	if err != nil {
		return "", err
	}

	if os.Getenv("DISABLE_RESEND") == "true" {
		return id, nil
	}

	_, err = SendRecoverUsernameEmail(username, e.Address)
	if err != nil {
		return "", err
	}

	return id, nil
}

func SendRecoverUsernameEmail(username string, email string) (resend.SendEmailResponse, error) {
	// TODO: Extract this so we aren't building a new client on each request
	client := resend.NewClient(os.Getenv("RESEND_API_KEY"))
	params := &resend.SendEmailRequest{
		To:   []string{email},
		From: "verify@petrichormud.com",
		// TODO: Add a doc for what to do if the user didn't request this
		// TODO: Link to that doc here
		Html:    fmt.Sprintf("You received this email as part of recovering your Username. Your username is: %s", username),
		Subject: "[PetrichorMUD] Username Recovery",
		ReplyTo: "support@petrichormud.com",
	}
	return client.Emails.Send(params)
}
