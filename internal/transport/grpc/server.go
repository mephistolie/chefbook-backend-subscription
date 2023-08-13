package grpc

import (
	api "github.com/mephistolie/chefbook-backend-subscription/api/proto/implementation/v1"
	"github.com/mephistolie/chefbook-backend-subscription/internal/transport/dependencies/service"
)

type SubscriptionServer struct {
	api.SubscriptionServiceServer
	service service.Subscription
}

func NewServer(service service.Subscription) *SubscriptionServer {
	return &SubscriptionServer{service: service}
}
