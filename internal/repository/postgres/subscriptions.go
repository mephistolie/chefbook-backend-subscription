package postgres

import (
	"fmt"
	"github.com/google/uuid"
	"github.com/mephistolie/chefbook-backend-common/log"
	"github.com/mephistolie/chefbook-backend-common/responses/fail"
	"github.com/mephistolie/chefbook-backend-common/subscription"
	"github.com/mephistolie/chefbook-backend-subscription/internal/entity"
	"sort"
	"time"
)

func (r *Repository) GetProfileSubscriptions(userId uuid.UUID) []entity.Subscription {
	var subscriptions []entity.Subscription

	query := fmt.Sprintf(`
		SELECT plan, source, expiration_timestamp, auto_renew
		FROM %s
		WHERE user_id=$1 AND expiration_timestamp > $2 AND start_timestamp <= $3
	`, subscriptionsTable)

	rows, err := r.db.Query(query, userId, time.Now().Add(24*time.Hour), time.Now())
	if err != nil {
		log.Errorf("unable to get profile %s subscriptions: %s", userId, err)
		return []entity.Subscription{}
	}

	for rows.Next() {
		sub := entity.Subscription{}
		if err = rows.Scan(&sub.Plan, &sub.Source, &sub.Expiration, &sub.AutoRenew); err != nil {
			log.Errorf("unable to parse profile %s subscription: %s", userId, err)
			continue
		}
		subscriptions = append(subscriptions, sub)
	}

	sort.Slice(subscriptions, func(i, j int) bool {
		return subscriptions[i].Plan == subscription.PlanPremium && subscriptions[j].Plan == subscription.PlanMaximum
	})

	return subscriptions
}

func (r *Repository) ClaimProfileSubscription(input entity.SubscriptionInput) error {
	_, err := r.getExistingSubscription(input.UserId, input.Plan, input.Source)
	if err != nil {
		return r.createSubscription(input)
	}
	return r.UpdateProfileSubscription(input)
}

func (r *Repository) getExistingSubscription(userId uuid.UUID, plan string, source string) (*time.Time, error) {
	var expirationTimestamp *time.Time

	query := fmt.Sprintf(`
		SELECT expiration_timestamp
		FROM %s
		WHERE user_id=$1 AND plan=$2 AND source=$3
	`, subscriptionsTable)

	row := r.db.QueryRow(query, userId, plan, source)
	err := row.Scan(&expirationTimestamp)

	return expirationTimestamp, err
}

func (r *Repository) createSubscription(input entity.SubscriptionInput) error {
	query := fmt.Sprintf(`
		INSERT INTO %s (user_id, plan, source, start_timestamp, expiration_timestamp, auto_renew)
		VALUES ($1, $2, $3, $4, $5, $6)
	`, subscriptionsTable)

	if _, err := r.db.Exec(query, input.UserId, input.Plan, input.Source, input.Start, input.Expiration, input.AutoRenew); err != nil {
		if isUniqueViolationError(err) {
			return nil
		}
		log.Errorf("unable to add profile %s subscription: %s", input.UserId, err)
		return fail.GrpcUnknown
	}

	return nil
}

func (r *Repository) UpdateProfileSubscription(input entity.SubscriptionInput) error {
	query := fmt.Sprintf(`
		UPDATE %s
		SET expiration_timestamp=$4, auto_renew=$5
		WHERE user_id=$1 AND plan=$2 AND source=$3 AND expiration_timestamp<$4
	`, subscriptionsTable)

	if _, err := r.db.Exec(query, input.UserId, input.Plan, input.Source, input.Expiration, input.AutoRenew); err != nil {
		log.Errorf("unable to update user %s subscription: %s", input.UserId, input.Plan)
		return fail.GrpcUnknown
	}

	return nil
}

func (r *Repository) SetProfileSubscriptionAutoRenewStatus(input entity.SubscriptionInput) error {
	query := fmt.Sprintf(`
		UPDATE %s
		SET auto_renew=$4
		WHERE user_id=$1 AND plan=$2 AND source=$3
	`, subscriptionsTable)

	if _, err := r.db.Exec(query, input.UserId, input.Plan, input.Source, input.AutoRenew); err != nil {
		log.Errorf("unable to update user %s subscription auto renew status: %s", input.UserId, input.Plan)
		return fail.GrpcUnknown
	}

	return nil
}

func (r *Repository) EndProfileSubscription(userId uuid.UUID, plan string, source string) error {
	query := fmt.Sprintf(`
		UPDATE %s
		SET expiration_timestamp=$4, auto_renew=false
		WHERE user_id=$1 AND plan=$2 AND source=$3
	`, subscriptionsTable)

	if _, err := r.db.Exec(query, userId, plan, source, time.Now()); err != nil {
		log.Errorf("unable to end profile %s subscription: %s", userId, plan)
		return fail.GrpcUnknown
	}

	return nil
}
