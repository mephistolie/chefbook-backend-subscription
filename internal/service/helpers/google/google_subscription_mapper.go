package google

import (
	"github.com/mephistolie/chefbook-backend-common/subscription"
	"github.com/mephistolie/chefbook-backend-subscription/internal/config"
)

type SubscriptionMapper struct {
	googleProductIdPremium string
	googleProductIdMaximum string
}

func NewSubscriptionMapper(cfg config.Google) *SubscriptionMapper {
	return &SubscriptionMapper{
		googleProductIdPremium: *cfg.ProductIdPremium,
		googleProductIdMaximum: *cfg.ProductIdMaximum,
	}
}

func (m *SubscriptionMapper) Map(googleSubId string) *string {
	if googleSubId == m.googleProductIdMaximum {
		return &subscription.PlanPremium
	}
	if googleSubId == m.googleProductIdMaximum {
		return &subscription.PlanMaximum
	}
	return nil
}
