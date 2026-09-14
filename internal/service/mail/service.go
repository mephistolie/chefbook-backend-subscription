package mail

import (
	"context"
	"time"

	"github.com/mephistolie/chefbook-backend-common/mail"
	"github.com/mephistolie/chefbook-backend-subscription/assets"
	"github.com/mephistolie/chefbook-backend-subscription/internal/config"
	"github.com/mephistolie/chefbook-backend-subscription/internal/logging"
)

type subscriptionPlanChangeMailValues struct {
	SubscriptionPlan string
}

type Service struct {
	sender       mail.Sender
	IsStub       bool
	IsDevEnv     bool
	sendAttempts int
}

func NewService(cfg *config.Config) (*Service, error) {
	var mailSender mail.Sender = mail.NewStubSender()
	var err error = nil
	if len(*cfg.Smtp.Host) > 0 {
		if mailSender, err = mail.NewSmtpSender(
			*cfg.Smtp.Email,
			*cfg.Smtp.Password,
			*cfg.Smtp.Host,
			*cfg.Smtp.Port,
			30*time.Second,
			*cfg.Smtp.Username,
		); err != nil {
			return nil, err
		}
	}
	return &Service{
		sender:       mailSender,
		IsStub:       len(*cfg.Smtp.Host) == 0,
		IsDevEnv:     *cfg.Environment == config.EnvDev,
		sendAttempts: *cfg.Smtp.SendAttempts,
	}, nil
}

func (s *Service) SendEncryptedVaultDeletionMail(ctx context.Context, email, plan string) {
	(logging.Events{}).SubscriptionMailSending(ctx, plan)
	payload := mail.Payload{
		To:      email,
		Subject: "ChefBook Subscription Plan Change",
	}
	mailValues := subscriptionPlanChangeMailValues{
		SubscriptionPlan: plan,
	}
	if err := payload.SetHtmlBody(assets.SubscriptionPlanChangeMailTmplFilePath, mailValues); err != nil {
		(logging.Events{}).SubscriptionMailRenderFailed(ctx, plan, err)
	}
	s.sendMessage(ctx, plan, payload)
}

func (s *Service) sendMessage(ctx context.Context, plan string, payload mail.Payload) {
	if s.IsDevEnv {
		payload.Body = "DEV\n" + payload.Body
	}
	if err := s.sender.Send(payload, s.sendAttempts); err != nil {
		(logging.Events{}).SubscriptionMailSendFailed(ctx, plan, err)
	}
}
