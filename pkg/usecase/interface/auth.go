package interfaces

import (
	"context"

	"github.com/stebinsabu13/ecommerce-api/pkg/utils"
)

type AuthUseCase interface {
	GoogleLogin(context.Context, utils.BodySignUpuser) (uint, error)
}
