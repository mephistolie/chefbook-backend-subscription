package rest

import (
	"context"
	_ "golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	_ "golang.org/x/oauth2/google"
	"net/http"
)

type Api interface {
	GetSubscriptionInfo(ctx context.Context, subscriptionId string, purchaseToken string) (*SubscriptionPurchase, error)
	AcknowledgeSubscriptionInfo(ctx context.Context, subscriptionId string, purchaseToken string) error
	CancelSubscriptionInfo(ctx context.Context, subscriptionId string, purchaseToken string) error
}

type Client struct {
	client      *http.Client
	packageName string
}

func NewClient(ctx context.Context, packageName string, jsonKey []byte) (*Client, error) {
	jwtCfg, err := google.JWTConfigFromJSON(
		jsonKey,
		"https://www.googleapis.com/auth/androidpublisher",
	)
	if err != nil {
		return nil, err
	}

	client := jwtCfg.Client(ctx)

	return &Client{
		client:      client,
		packageName: packageName,
	}, nil
}
