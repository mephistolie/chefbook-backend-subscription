package rest

import (
	"encoding/json"
	"errors"
	"io"
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

func (c *Client) GetSubscriptionInfo(subscriptionId string, purchaseToken string) (*SubscriptionPurchase, error) {
	url := "https://androidpublisher.googleapis.com/androidpublisher/v3/applications/" + c.packageName +
		"/purchases/subscriptions/" + subscriptionId +
		"/tokens/" + purchaseToken

	res, err := c.client.Get(url)
	if err != nil {
		return nil, err
	}
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

func (c *Client) AcknowledgeSubscriptionInfo(subscriptionId string, purchaseToken string) error {
	url := "https://androidpublisher.googleapis.com/androidpublisher/v3/applications/" + c.packageName +
		"/purchases/subscriptions/" + subscriptionId +
		"/tokens/" + purchaseToken + ":acknowledge"

	res, err := c.client.Post(url, "application/json", strings.NewReader(""))
	if err != nil {
		return err
	}
	if res.StatusCode != 200 {
		return errors.New("error status code")
	}

	return nil
}

func (c *Client) CancelSubscriptionInfo(subscriptionId string, purchaseToken string) error {
	url := "https://androidpublisher.googleapis.com/androidpublisher/v3/applications/" + c.packageName +
		"/purchases/subscriptions/" + subscriptionId +
		"/tokens/" + purchaseToken + ":cancel"

	res, err := c.client.Post(url, "application/json", strings.NewReader(""))
	if err != nil {
		return err
	}
	if res.StatusCode != 200 {
		return errors.New("error status code")
	}

	return nil
}
