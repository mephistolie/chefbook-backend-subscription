package rest

import (
	"context"
	_ "golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	_ "golang.org/x/oauth2/google"
	"net/http"
)

type Api interface {
	GetSubscriptionInfo(subscriptionId string, purchaseToken string) (*SubscriptionPurchase, error)
	AcknowledgeSubscriptionInfo(subscriptionId string, purchaseToken string) error
	CancelSubscriptionInfo(subscriptionId string, purchaseToken string) error
}

type Client struct {
	client      *http.Client
	packageName string
}

func NewClient(packageName string, jsonKey []byte) (*Client, error) {
	jwtCfg, err := google.JWTConfigFromJSON(
		jsonKey,
		"https://www.googleapis.com/auth/androidpublisher",
	)
	if err != nil {
		return nil, err
	}

	client := jwtCfg.Client(context.Background())

	return &Client{
		client:      client,
		packageName: packageName,
	}, nil
}
