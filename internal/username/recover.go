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

	u, err := i.Queries.GetPlayerUsernameById(context.Background(), e.PID)
	if err != nil {
		return "", err
	}

	if os.Getenv("DISABLE_RESEND") == "true" {
		return id, nil
	}

	_, err = SendRecoverUsernameEmail(i, u, e.Address)
	if err != nil {
		return "", err
	}

	return id, nil
}

func SendRecoverUsernameEmail(i *shared.Interfaces, username string, email string) (resend.SendEmailResponse, error) {
	params := &resend.SendEmailRequest{
		To:   []string{email},
		From: "verify@petrichormud.com",
		// TODO: If the user didn't request this, link to the section of the profile for changing your password
		Html:    fmt.Sprintf("You received this email as part of recovering your Username. Your username is: %s", username),
		Subject: "[PetrichorMUD] Username Recovery",
		ReplyTo: "support@petrichormud.com",
	}
	return i.Resend.Emails.Send(params)
}
