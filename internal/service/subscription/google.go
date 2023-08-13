package subscription

import (
	"context"
	"github.com/google/uuid"
	authApi "github.com/mephistolie/chefbook-backend-auth/api/proto/implementation/v1"
	"github.com/mephistolie/chefbook-backend-common/log"
	"github.com/mephistolie/chefbook-backend-common/responses/fail"
	"github.com/mephistolie/chefbook-backend-subscription/internal/entity"
	subscriptionFail "github.com/mephistolie/chefbook-backend-subscription/internal/entity/fail"
	"github.com/mephistolie/chefbook-backend-subscription/pkg/subscription/google/cloud/pubsub"
	"github.com/mephistolie/chefbook-backend-subscription/pkg/subscription/google/rest"
	"strconv"
	"time"
)

func (s *Service) ConfirmGoogleSubscription(userId uuid.UUID, googleSubId, purchaseToken string) error {
	input, err := s.getGoogleSubscriptionInput(userId, googleSubId, purchaseToken)
	if err != nil {
		return nil
	}

	if err = s.repo.ClaimProfileSubscription(input); err != nil {
		return err
	}

	if err = s.googleApi.AcknowledgeSubscriptionInfo(googleSubId, purchaseToken); err != nil {
		log.Errorf("unable to acknowledge purchase: %s", err)
		return err
	}

	go func() {
		if info, err := s.grpc.Auth.GetAuthInfo(context.Background(), &authApi.GetAuthInfoRequest{Id: userId.String()}); err == nil {
			s.mail.SendEncryptedVaultDeletionMail(info.Email, input.Plan)
		}
	}()

	return nil
}

func (s *Service) HandleSubscriptionEvent(event pubsub.SubscriptionEvent) error {
	switch event.NotificationType {
	case pubsub.NotificationTypeSubscriptionRenewed:
		return s.onSubscriptionRenewed(event)
	case pubsub.NotificationTypeSubscriptionCanceled:
		return s.onSubscriptionAutoRenewStatusChanged(event, false)
	case pubsub.NotificationTypeSubscriptionRestarted:
		return s.onSubscriptionAutoRenewStatusChanged(event, true)
	case pubsub.NotificationTypeSubscriptionRevoked:
		return s.onSubscriptionRevoked(event)
	default:
		return nil
	}
}

func (s *Service) onSubscriptionRenewed(event pubsub.SubscriptionEvent) error {
	userId, err := s.repo.GetUserIdByGooglePurchaseToken(event.PurchaseToken)
	if err != nil {
		return err
	}
	if userId == nil {
		return nil
	}

	input, err := s.getGoogleSubscriptionInput(*userId, event.SubscriptionId, event.PurchaseToken)
	if err != nil {
		if err == subscriptionFail.GrpcSubscriptionInactive || err == subscriptionFail.GrpcSubscriptionExpired {
			return nil
		}
		return err
	}

	return s.repo.UpdateProfileSubscription(input)
}

func (s *Service) onSubscriptionAutoRenewStatusChanged(event pubsub.SubscriptionEvent, autoRenew bool) error {
	userId, err := s.repo.GetUserIdByGooglePurchaseToken(event.PurchaseToken)
	if err != nil {
		return err
	}
	if userId == nil {
		return nil
	}

	plan := s.googleSubMapper.Map(event.SubscriptionId)
	if plan == nil {
		log.Errorf("unable to parse google subscription ID %s in RTDN", event.SubscriptionId)
		return subscriptionFail.GrpcInvalidSubscriptionId
	}

	return s.repo.SetProfileSubscriptionAutoRenewStatus(entity.SubscriptionInput{
		UserId:    *userId,
		Plan:      *plan,
		Source:    entity.SourceGoogle,
		AutoRenew: autoRenew,
	})
}

func (s *Service) onSubscriptionRevoked(event pubsub.SubscriptionEvent) error {
	userId, err := s.repo.GetUserIdByGooglePurchaseToken(event.PurchaseToken)
	if err != nil {
		return err
	}
	if userId == nil {
		return nil
	}

	subscriptionId := s.googleSubMapper.Map(event.SubscriptionId)
	if subscriptionId == nil {
		log.Errorf("unable to parse google subscription ID %s in RTDN", event.SubscriptionId)
		return subscriptionFail.GrpcInvalidSubscriptionId
	}

	return s.repo.EndProfileSubscription(*userId, *subscriptionId, entity.SourceGoogle)
}

func (s *Service) getGoogleSubscriptionInput(userId uuid.UUID, googleSubId, purchaseToken string) (entity.SubscriptionInput, error) {
	if s.googleApi == nil {
		log.Warnf("google subscription is disabled")
		return entity.SubscriptionInput{}, subscriptionFail.GrpcInvalidPaymentService
	}

	plan := s.googleSubMapper.Map(googleSubId)
	if plan == nil {
		return entity.SubscriptionInput{}, subscriptionFail.GrpcInvalidSubscriptionId
	}

	info, err := s.googleApi.GetSubscriptionInfo(googleSubId, purchaseToken)
	if err != nil {
		log.Debugf("unable to validate purchase: %s", err)
		return entity.SubscriptionInput{}, fail.CreateGrpcClient(fail.TypeInvalidBody, "unable to validate purchase")
	}

	input, err := s.validateGoogleSubscription(userId, *plan, *info)
	if err != nil {
		return entity.SubscriptionInput{}, err
	}

	return input, nil
}

func (s *Service) validateGoogleSubscription(userId uuid.UUID, plan string, info rest.SubscriptionPurchase) (entity.SubscriptionInput, error) {
	if info.PaymentState == nil || *info.PaymentState == rest.PaymentStatePending {
		return entity.SubscriptionInput{}, subscriptionFail.GrpcSubscriptionInactive
	}

	strStartTime := info.StartTimeMillis
	if info.AutoResumeTimeMillis != nil && *info.AutoResumeTimeMillis > info.StartTimeMillis {
		strStartTime = *info.AutoResumeTimeMillis
	}
	rawStartTime, err := strconv.ParseInt(strStartTime, 10, 64)
	if err != nil {
		log.Errorf("unable to parse subscription start time: %s", err)
		return entity.SubscriptionInput{}, fail.GrpcUnknown
	}
	startTime := time.UnixMilli(rawStartTime)
	rawExpirationTime, err := strconv.ParseInt(info.ExpiryTimeMillis, 10, 64)
	if err != nil {
		log.Errorf("unable to parse subscription expiration time: %s", err)
		return entity.SubscriptionInput{}, fail.GrpcUnknown
	}
	expirationTime := time.UnixMilli(rawExpirationTime)

	if time.Now().After(expirationTime) {
		return entity.SubscriptionInput{}, subscriptionFail.GrpcSubscriptionExpired
	}

	return entity.SubscriptionInput{
		UserId:     userId,
		Plan:       plan,
		Source:     entity.SourceGoogle,
		AutoRenew:  info.AutoRenewing,
		Start:      startTime,
		Expiration: expirationTime,
	}, nil
}
