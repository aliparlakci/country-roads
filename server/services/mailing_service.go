package services

import (
	"context"
	"github.com/mailgun/mailgun-go/v4"
)

type MailgunStore struct {
	Client *mailgun.MailgunImpl
}

type MailingService interface {
	MailSender
}

type MailSender interface {
	SendMail(ctx context.Context, subject, body, to string) error
}

func (m *MailgunStore) SendMail(ctx context.Context, subject, body, to string) error {
	from := "tuzlapool <mailgun@tuzlapool.xyz>"
	message := m.Client.NewMessage(from, subject, body, to)
	_, _, err := m.Client.Send(ctx, message)
	return err
}

func NewMailgunStore(client *mailgun.MailgunImpl) *MailgunStore {
	client.SetAPIBase(mailgun.APIBaseEU)
	return &MailgunStore{Client: client}
}