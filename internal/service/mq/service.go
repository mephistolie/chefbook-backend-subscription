package mq

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/google/uuid"
	auth "github.com/mephistolie/chefbook-backend-auth/api/mq"
	"github.com/mephistolie/chefbook-backend-common/firebase"
	"github.com/mephistolie/chefbook-backend-common/mq/model"
	"github.com/mephistolie/chefbook-backend-subscription/internal/logging"
	"github.com/mephistolie/chefbook-backend-subscription/internal/service/dependencies/repository"
)

type Service struct {
	repo     repository.MQ
	firebase *firebase.Client
}

func NewService(
	repo repository.MQ,
	firebase *firebase.Client,
) *Service {
	return &Service{
		repo:     repo,
		firebase: firebase,
	}
}

func (s *Service) HandleMessage(msg model.MessageData) error {
	ctx := context.Background()

	logging.Events{}.MQMessageProcessing(ctx, msg.Id.String(), msg.Type)
	switch msg.Type {
	case auth.MsgTypeProfileFirebaseImport:
		return s.handleFirebaseImportMsg(ctx, msg.Id, msg.Body)
	case auth.MsgTypeProfileDeleted:
		return s.handleProfileDeletedMsg(ctx, msg.Id, msg.Body)
	default:
		logging.Events{}.MQMessageTypeUnsupported(ctx, msg.Id.String(), msg.Type)
		return errors.New("not implemented")
	}
}

func (s *Service) handleFirebaseImportMsg(ctx context.Context, messageId uuid.UUID, data []byte) error {
	var body auth.MsgBodyProfileFirebaseImport
	if err := json.Unmarshal(data, &body); err != nil {
		return err
	}

	userId, err := uuid.Parse(body.UserId)
	if err != nil {
		return err
	}

	profile, err := s.firebase.GetProfile(ctx, body.FirebaseId)
	if err != nil {
		logging.Events{}.FirebaseProfileLoadFailed(ctx, body.UserId)
		return err
	}

	if profile.IsPremium {
		return s.repo.ImportPremiumVersion(ctx, userId, messageId)
	}

	return nil
}

func (s *Service) handleProfileDeletedMsg(ctx context.Context, messageId uuid.UUID, data []byte) error {
	var body auth.MsgBodyProfileDeleted
	if err := json.Unmarshal(data, &body); err != nil {
		return err
	}

	userId, err := uuid.Parse(body.UserId)
	if err != nil {
		return err
	}

	logging.Events{}.ProfileDeletedMessageProcessing(ctx, messageId.String(), body.UserId)
	return s.repo.DeleteProfile(ctx, userId, messageId)
}
