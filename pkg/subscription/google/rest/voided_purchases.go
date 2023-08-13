package rest

const (
	PurchaseKindSubscription = "subscriptionPurchase"
	PurchaseKindProduct      = "productPurchase"
)

type VoidedPurchases struct {
	VoidedPurchases []VoidedPurchase `json:"voidedPurchases"`
}

type VoidedPurchase struct {
	Kind               string `json:"kind"`
	PurchaseToken      string `json:"purchaseToken"`
	PurchaseTimeMillis string `json:"purchaseTimeMillis"`
	VoidedTimeMillis   string `json:"voidedTimeMillis"`
	OrderId            string `json:"orderId"`
	VoidedSource       int    `json:"voidedSource"`
	VoidedReason       int    `json:"voidedReason"`
}

//func (c *Client) GetVoidedSubscriptions() (*SubscriptionPurchase, error) {
//	url := "https://androidpublisher.googleapis.com/androidpublisher/v3/applications/" + c.packageName +
//		"/purchases/voidedpurchases"
//
//	res, err := c.client.Get(url)
//	if err != nil {
//		return nil, err
//	}
//	if res.StatusCode != 200 {
//		return nil, errors.New("error status code")
//	}
//
//	bodyBytes, err := io.ReadAll(res.Body)
//	if err != nil {
//		return nil, err
//	}
//
//	var body VoidedPurchases
//	if err = json.Unmarshal(bodyBytes, &body); err != nil {
//		return nil, err
//	}
//
//	return body, nil
//}
