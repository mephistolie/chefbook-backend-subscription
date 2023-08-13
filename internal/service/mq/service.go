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
	log.Infof("processing message %s with type %s", msg.Id, msg.Type)
	switch msg.Type {
	case auth.MsgTypeProfileFirebaseImport:
		return s.handleFirebaseImportMsg(msg.Id, msg.Body)
	case auth.MsgTypeProfileDeleted:
		return s.handleProfileDeletedMsg(msg.Id, msg.Body)
	default:
		log.Warnf("got unsupported message type %s for message %s", msg.Type, msg.Id)
		return errors.New("not implemented")
	}
}

func (s *Service) handleFirebaseImportMsg(messageId uuid.UUID, data []byte) error {
	var body auth.MsgBodyProfileFirebaseImport
	if err := json.Unmarshal(data, &body); err != nil {
		return err
	}

	userId, err := uuid.Parse(body.UserId)
	if err != nil {
		return err
	}

	profile, err := s.firebase.GetProfile(context.Background(), body.FirebaseId)
	if err != nil {
		log.Warnf("unable to get firebase profile for user %s: %s", body.UserId, err)
		return err
	}

	if profile.IsPremium {
		return s.repo.ImportPremiumVersion(userId, messageId)
	}

	return nil
}

func (s *Service) handleProfileDeletedMsg(messageId uuid.UUID, data []byte) error {
	var body auth.MsgBodyProfileDeleted
	if err := json.Unmarshal(data, &body); err != nil {
		return err
	}

	userId, err := uuid.Parse(body.UserId)
	if err != nil {
		return err
	}

	log.Infof("deleting user %s...", body.UserId)
	return s.repo.DeleteProfile(userId, messageId)
}
