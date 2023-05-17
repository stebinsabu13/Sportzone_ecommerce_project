package interfaces

import (
	"context"

	"github.com/stebinsabu13/ecommerce-api/pkg/domain"
	"github.com/stebinsabu13/ecommerce-api/pkg/utils"
)

type OtpUseCase interface {
	TwilioSendOTP(context.Context, string) (string, error)
	TwilioVerifyOTP(context.Context, utils.Otpverify) (domain.OtpSession, error)
}
