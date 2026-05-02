package postgres

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/mephistolie/chefbook-backend-common/log"
	"github.com/mephistolie/chefbook-backend-common/responses/fail"
)

func (r *Repository) GetUserIdByGooglePurchaseToken(ctx context.Context, purchaseToken string) (*uuid.UUID, error) {
	var userId uuid.UUID

	query := fmt.Sprintf(`
		SELECT user_id
		FROM %s
		WHERE purchase_token=$1
	`, googleTable)

	rows, err := r.db.QueryContext(ctx, query, purchaseToken)
	if err != nil {
		log.Warnf("unable to get profile for google subscription purchase token %s: %s", purchaseToken, err)
		return nil, fail.GrpcUnknown
	}
	defer rows.Close()

	for rows.Next() {
		if err := rows.Scan(&userId); err == nil {
			return &userId, nil
		}
	}
	if err = rows.Err(); err != nil {
		log.Warnf("unable to iterate profile for google subscription purchase token %s: %s", purchaseToken, err)
		return nil, fail.GrpcUnknown
	}

	return nil, nil
}
