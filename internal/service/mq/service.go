package mq

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/google/uuid"
	auth "github.com/mephistolie/chefbook-backend-auth/api/mq"
	"github.com/mephistolie/chefbook-backend-common/firebase"
	"github.com/mephistolie/chefbook-backend-common/log"
	"github.com/mephistolie/chefbook-backend-common/mq/model"
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

	log.Log(ctx, log.Event{
		Event:     "mq.message.processing",
		Message:   "processing message",
		Component: log.ComponentAMQP,
		MessageID: msg.Id.String(),
		Payload: map[string]any{
			"message_type": msg.Type,
		},
	})
	switch msg.Type {
	case auth.MsgTypeProfileFirebaseImport:
		return s.handleFirebaseImportMsg(ctx, msg.Id, msg.Body)
	case auth.MsgTypeProfileDeleted:
		return s.handleProfileDeletedMsg(ctx, msg.Id, msg.Body)
	default:
		log.LogWarn(ctx, log.Event{
			Event:     "mq.message.unsupported_type",
			Message:   "got unsupported message type",
			Component: log.ComponentAMQP,
			MessageID: msg.Id.String(),
			Payload: map[string]any{
				"message_type": msg.Type,
			},
		})
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
		log.LogWarnError(ctx, log.Event{
			Event:     "firebase.profile.load_failed",
			Message:   "unable to get firebase profile",
			Component: log.ComponentFirebase,
			UserID:    body.UserId,
		}, err)
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

	log.Log(ctx, log.Event{
		Event:     "profile.deleted.message.processing",
		Message:   "processing profile deleted message",
		Component: log.ComponentAMQP,
		MessageID: messageId.String(),
		UserID:    body.UserId,
	})
	return s.repo.DeleteProfile(ctx, userId, messageId)
}
