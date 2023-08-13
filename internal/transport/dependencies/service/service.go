package service

import (
	"github.com/google/uuid"
	"github.com/mephistolie/chefbook-backend-common/firebase"
	"github.com/mephistolie/chefbook-backend-common/log"
	mq "github.com/mephistolie/chefbook-backend-common/mq/dependencies"
	"github.com/mephistolie/chefbook-backend-subscription/internal/config"
	"github.com/mephistolie/chefbook-backend-subscription/internal/entity"
	"github.com/mephistolie/chefbook-backend-subscription/internal/repository/grpc"
	"github.com/mephistolie/chefbook-backend-subscription/internal/repository/postgres"
	"github.com/mephistolie/chefbook-backend-subscription/internal/service/helpers/google"
	"github.com/mephistolie/chefbook-backend-subscription/internal/service/mail"
	mqInbox "github.com/mephistolie/chefbook-backend-subscription/internal/service/mq"
	"github.com/mephistolie/chefbook-backend-subscription/internal/service/subscription"
	googleRest "github.com/mephistolie/chefbook-backend-subscription/pkg/subscription/google/rest"
)

type Service struct {
	Subscription Subscription
	MQ           mq.Inbox
}

type Subscription interface {
	GetProfileSubscriptions(userId uuid.UUID) []entity.Subscription
	GetProfileCurrentSubscription(userId uuid.UUID) entity.Subscription

	ConfirmGoogleSubscription(userId uuid.UUID, googleSubId, purchaseToken string) error
}

func New(
	repo *postgres.Repository,
	grpcRepository *grpc.Repository,
	cfg *config.Config,
) (*Service, error) {

	mailService, err := mail.NewService(cfg)
	if err != nil {
		return nil, err
	}

	googleSubMapper := google.NewSubscriptionMapper(cfg.Google)

	var firebaseClient *firebase.Client = nil
	if len(*cfg.Firebase.Credentials) > 0 {
		firebaseClient, err = firebase.NewClient([]byte(*cfg.Firebase.Credentials), "")
		if err != nil {
			return nil, err
		}
		log.Info("Firebase client initialized")
	}

	var googleClient *googleRest.Client = nil
	if len(*cfg.Google.JsonKey) > 0 {
		googleClient, err = googleRest.NewClient(*cfg.Google.PackageName, []byte(*cfg.Google.JsonKey))
		if err != nil {
			return nil, err
		}
		log.Info("Google client initialized")
	}

	return &Service{
		Subscription: subscription.NewService(repo, grpcRepository, mailService, googleClient, googleSubMapper),
		MQ:           mqInbox.NewService(repo, firebaseClient),
	}, nil
}
