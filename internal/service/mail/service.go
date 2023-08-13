package mail

import (
	"github.com/mephistolie/chefbook-backend-common/log"
	"github.com/mephistolie/chefbook-backend-common/mail"
	"github.com/mephistolie/chefbook-backend-subscription/assets"
	"github.com/mephistolie/chefbook-backend-subscription/internal/config"
	"time"
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

func (s *Service) SendEncryptedVaultDeletionMail(email, plan string) {
	log.Info("sending subscription plan change mail to ", email)
	payload := mail.Payload{
		To:      email,
		Subject: "ChefBook Subscription Plan Change",
	}
	mailValues := subscriptionPlanChangeMailValues{
		SubscriptionPlan: plan,
	}
	if err := payload.SetHtmlBody(assets.SubscriptionPlanChangeMailTmplFilePath, mailValues); err != nil {
		log.Error("failed to set HTML Body for mail: ", err)
	}
	s.sendMessage(payload)
}

func (s *Service) sendMessage(payload mail.Payload) {
	if s.IsDevEnv {
		payload.Body = "DEV\n" + payload.Body
	}
	_ = s.sender.Send(payload, s.sendAttempts)
}
