package username

import (
	"context"
	"fmt"
	"os"

	resend "github.com/resendlabs/resend-go"
	"petrichormud.com/app/internal/shared"
)

func RecoverUsername(i *shared.Interfaces, pid int64, email string) error {
	username, err := i.Queries.GetPlayerUsernameById(context.Background(), pid)
	if err != nil {
		return err
	}

	if os.Getenv("DISABLE_RESEND") == "true" {
		return nil
	}

	_, err = SendRecoverUsernameEmail(username, email)
	if err != nil {
		return err
	}

	return nil
}

func SendRecoverUsernameEmail(username string, email string) (resend.SendEmailResponse, error) {
	// TODO: Extract this so we aren't building a new client on each request
	client := resend.NewClient(os.Getenv("RESEND_API_KEY"))
	params := &resend.SendEmailRequest{
		To:      []string{email},
		From:    "verify@petrichormud.com",
		Html:    fmt.Sprintf("You received this email as part of recovering your Username. Your username is: %s", username),
		Subject: fmt.Sprintf("[PetrichorMUD] Username Recovery"),
		ReplyTo: "support@petrichormud.com",
	}
	return client.Emails.Send(params)
}
