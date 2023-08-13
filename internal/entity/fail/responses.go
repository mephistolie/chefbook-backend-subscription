package fail

import "github.com/mephistolie/chefbook-backend-common/responses/fail"

var (
	GrpcInvalidPaymentService = fail.CreateGrpcClient(fail.TypeInvalidBody, "subscription payment service not supported")
	GrpcInvalidSubscriptionId = fail.CreateGrpcClient(fail.TypeInvalidBody, "subscription ID not supported")
	GrpcSubscriptionInactive  = fail.CreateGrpcClient(fail.TypeInvalidBody, "subscription inactive")
	GrpcSubscriptionExpired   = fail.CreateGrpcClient(fail.TypeInvalidBody, "subscription expired")
)
