package mail

import (
	"context"
	"strings"

	"github.com/a-h/templ"
	"github.com/juevigrace/diva-server/internal/models"
	resend "github.com/resend/resend-go/v2"
)

type Client struct {
	apiKey    string
	fromEmail string
}

func NewClient(apiKey, fromEmail string) *Client {
	return &Client{
		apiKey:    apiKey,
		fromEmail: fromEmail,
	}
}

func (c *Client) Send(ctx context.Context, to, subject string, component templ.Component) error {
	client := resend.NewClient(c.apiKey)

	var sb strings.Builder
	if err := component.Render(ctx, &sb); err != nil {
		return err
	}

	_, err := client.Emails.SendWithContext(ctx, &resend.SendEmailRequest{
		From:    c.fromEmail,
		To:      []string{to},
		Subject: subject,
		Html:    sb.String(),
	})

	return err
}

func (c *Client) SendVerificationEmail(ctx context.Context, to string, verification *models.UserVerification) error {
	return c.Send(ctx, to, "Email Verification", VerificationEmail(verification))
}
