package postgres

import (
	"database/sql"
	"fmt"
	"github.com/google/uuid"
	"github.com/mephistolie/chefbook-backend-common/log"
	"github.com/mephistolie/chefbook-backend-common/responses/fail"
	"github.com/mephistolie/chefbook-backend-common/subscription"
	"github.com/mephistolie/chefbook-backend-subscription/internal/entity"
	"time"
)

func (r *Repository) GetExpiringSubscriptions() []entity.ExpiringSubscription {
	var subscriptions []entity.ExpiringSubscription

	query := fmt.Sprintf(`
		SELECT user_id, plan, source
		FROM %s
		WHERE expiration_timestamp > $1 AND expiration_timestamp <= $1 AND auto_renew=true
	`, subscriptionsTable)

	rows, err := r.db.Query(query, time.Now().Add(-2*24*time.Hour), time.Now().Add(6*time.Hour))
	if err != nil {
		log.Errorf("unable to get expiring subscriptions: %s", err)
		return []entity.ExpiringSubscription{}
	}

	for rows.Next() {
		sub := entity.ExpiringSubscription{}
		if err = rows.Scan(&sub.UserId, &sub.Plan, &sub.Source); err != nil {
			log.Errorf("unable to parse expiring subscription: %s", err)
			continue
		}
		subscriptions = append(subscriptions, sub)
	}

	return subscriptions
}

func (r *Repository) ImportPremiumVersion(userId, messageId uuid.UUID) error {
	tx, err := r.handleMessageIdempotently(messageId)
	if err != nil {
		if isUniqueViolationError(err) {
			return nil
		} else {
			return fail.GrpcUnknown
		}
	}

	query := fmt.Sprintf(`
		INSERT INTO %s (user_id, plan, source)
		VALUES ($1, $2, $3)
	`, subscriptionsTable)

	if _, err = tx.Exec(query, userId, subscription.PlanPremium, entity.SourceFirebase); err != nil {
		log.Warnf("unable to import premium app version for profile %s: %s", userId, err)
		return errorWithTransactionRollback(tx, fail.GrpcUnknown)
	}

	return commitTransaction(tx)
}

func (r *Repository) DeleteProfile(userId, messageId uuid.UUID) error {
	tx, err := r.handleMessageIdempotently(messageId)
	if err != nil {
		if isUniqueViolationError(err) {
			return nil
		} else {
			return fail.GrpcUnknown
		}
	}

	query := fmt.Sprintf(`
		DELETE FROM %s
		WHERE user_id=$1
	`, subscriptionsTable)

	if _, err = tx.Exec(query, userId); err != nil {
		log.Warnf("unable to delete profile %s subscriptions: %s", userId, err)
		return errorWithTransactionRollback(tx, fail.GrpcUnknown)
	}

	return commitTransaction(tx)
}

func (r *Repository) handleMessageIdempotently(messageId uuid.UUID) (*sql.Tx, error) {
	tx, err := r.startTransaction()
	if err != nil {
		return nil, err
	}

	addMessageQuery := fmt.Sprintf(`
		INSERT INTO %s (message_id)
		VALUES ($1)
	`, inboxTable)

	if _, err = tx.Exec(addMessageQuery, messageId); err != nil {
		if !isUniqueViolationError(err) {
			log.Error("unable to add message to inbox: ", err)
		}
		return nil, errorWithTransactionRollback(tx, err)
	}

	deleteOutdatedMessagesQuery := fmt.Sprintf(`
		DELETE FROM %[1]v
		WHERE ctid IN
		(
			SELECT ctid IN
			FROM %[1]v
			ORDER BY timestamp DESC
			OFFSET 1000
		)
	`, inboxTable)

	if _, err = tx.Exec(deleteOutdatedMessagesQuery); err != nil {
		return nil, errorWithTransactionRollback(tx, err)
	}

	return tx, nil
}
