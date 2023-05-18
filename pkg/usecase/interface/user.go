package interfaces

import (
	"context"

	domain "github.com/stebinsabu13/ecommerce-api/pkg/domain"
	"github.com/stebinsabu13/ecommerce-api/pkg/utils"
)

type UserUseCase interface {
	// FindAll(ctx context.Context) ([]domain.User, error)
	FindbyEmail(ctx context.Context, email string) (domain.User, error)
	FindbyUserID(context.Context, uint) (domain.User, error)
	UpdateVerify(context.Context, string) error
	SignUpUser(ctx context.Context, user domain.User) (string, error)
	FindbyEmailorMobilenum(ctx context.Context, body utils.OtpLogin) (domain.User, error)
	ShowDetails(ctx context.Context, id int) (utils.ResponseUsers, error)
	ShowAddress(ctx context.Context, id int) ([]utils.Address, error)
	AddAddress(context.Context, domain.Address) error
	EditProfile(context.Context, utils.EditProfileReq, uint) error
}
