package interfaces

import (
	"context"

	"github.com/stebinsabu13/ecommerce-api/pkg/domain"
)

type OtpRepository interface {
	SaveOtp(context.Context, domain.OtpSession) error
	RetrieveSession(context.Context, string) (domain.OtpSession, error)
}
