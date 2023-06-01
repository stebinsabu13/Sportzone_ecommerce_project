package interfaces

import (
	"context"

	"github.com/stebinsabu13/ecommerce-api/pkg/domain"
	"github.com/stebinsabu13/ecommerce-api/pkg/utils"
)

type UserRepository interface {
	// FindAll(ctx context.Context) ([]domain.User, error)
	FindbyEmail(ctx context.Context, email string) (domain.User, error)
	FindbyUserID(context.Context, uint) (domain.User, error)
	UpdateVerify(string, string) error
	SignUpUser(ctx context.Context, user domain.User) (string, error)
	FindbyEmailorMobilenum(ctx context.Context, body utils.OtpLogin) (domain.User, error)
	ShowDetails(ctx context.Context, id uint) (utils.ResponseUsers, error)
	ShowAddress(ctx context.Context, id uint) ([]utils.Address, error)
	AddAddress(context.Context, domain.Address) error
	EditProfile(context.Context, utils.EditProfileReq, uint) error
	ChangePassword(context.Context, string, string) error
	ListAllCategories(ctx context.Context) ([]utils.ResponseCategory, error)
	ViewWallet(uint) ([]utils.ResWallet, int, error)
}
