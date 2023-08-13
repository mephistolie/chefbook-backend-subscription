package subscription

import (
	"github.com/mephistolie/chefbook-backend-subscription/internal/repository/grpc"
	"github.com/mephistolie/chefbook-backend-subscription/internal/service/dependencies/repository"
	"github.com/mephistolie/chefbook-backend-subscription/internal/service/helpers/google"
	"github.com/mephistolie/chefbook-backend-subscription/internal/service/mail"
	googleRest "github.com/mephistolie/chefbook-backend-subscription/pkg/subscription/google/rest"
)

var ()

type Service struct {
	repo            repository.Subscription
	grpc            *grpc.Repository
	mail            *mail.Service
	googleApi       googleRest.Api
	googleSubMapper *google.SubscriptionMapper
}

func NewService(
	repo repository.Subscription,
	grpc *grpc.Repository,
	mail *mail.Service,
	googleApi googleRest.Api,
	googleSubMapper *google.SubscriptionMapper,
) *Service {
	return &Service{
		repo:            repo,
		grpc:            grpc,
		mail:            mail,
		googleApi:       googleApi,
		googleSubMapper: googleSubMapper,
	}
}
