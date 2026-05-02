package rest

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"strings"
)

const (
	PaymentStatePending         = 0
	AcknowledgementStatePending = 0
)

type SubscriptionPurchase struct {
	Kind                 string  `json:"kind"`
	StartTimeMillis      string  `json:"startTimeMillis"`
	ExpiryTimeMillis     string  `json:"expiryTimeMillis"`
	AutoResumeTimeMillis *string `json:"autoResumeTimeMillis"`
	AutoRenewing         bool    `json:"autoRenewing"`
	PaymentState         *int    `json:"paymentState"`
	AcknowledgementState int     `json:"acknowledgementState"`
}

func (c *Client) GetSubscriptionInfo(ctx context.Context, subscriptionId string, purchaseToken string) (*SubscriptionPurchase, error) {
	url := "https://androidpublisher.googleapis.com/androidpublisher/v3/applications/" + c.packageName +
		"/purchases/subscriptions/" + subscriptionId +
		"/tokens/" + purchaseToken

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}
	res, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	if res.StatusCode != 200 {
		return nil, errors.New("error status code")
	}

	bodyBytes, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	var body SubscriptionPurchase
	if err = json.Unmarshal(bodyBytes, &body); err != nil {
		return nil, err
	}

	return &body, nil
}

func (c *Client) AcknowledgeSubscriptionInfo(ctx context.Context, subscriptionId string, purchaseToken string) error {
	url := "https://androidpublisher.googleapis.com/androidpublisher/v3/applications/" + c.packageName +
		"/purchases/subscriptions/" + subscriptionId +
		"/tokens/" + purchaseToken + ":acknowledge"

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, strings.NewReader(""))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	res, err := c.client.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()
	if res.StatusCode != 200 {
		return errors.New("error status code")
	}

	return nil
}

func (c *Client) CancelSubscriptionInfo(ctx context.Context, subscriptionId string, purchaseToken string) error {
	url := "https://androidpublisher.googleapis.com/androidpublisher/v3/applications/" + c.packageName +
		"/purchases/subscriptions/" + subscriptionId +
		"/tokens/" + purchaseToken + ":cancel"

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, strings.NewReader(""))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	res, err := c.client.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()
	if res.StatusCode != 200 {
		return errors.New("error status code")
	}

	return nil
}
