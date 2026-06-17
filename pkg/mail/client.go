package mail

import (
	"context"
	"strings"

	"github.com/a-h/templ"
	resend "github.com/resend/resend-go/v2"
)

type Client struct {
	rc        *resend.Client
	fromEmail string
}

func NewClient(apiKey, fromEmail string) *Client {
	return &Client{
		rc:        resend.NewClient(apiKey),
		fromEmail: fromEmail,
	}
}

func (c *Client) Send(ctx context.Context, to, subject string, component templ.Component) error {
	var sb strings.Builder
	if err := component.Render(ctx, &sb); err != nil {
		return err
	}

	_, err := c.rc.Emails.SendWithContext(ctx, &resend.SendEmailRequest{
		From:    c.fromEmail,
		To:      []string{to},
		Subject: subject,
		Html:    sb.String(),
	})

	return err
}
