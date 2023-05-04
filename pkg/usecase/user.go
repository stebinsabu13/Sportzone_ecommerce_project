package usecase

import (
	"context"

	domain "github.com/stebinsabu13/ecommerce-api/pkg/domain"
	interfaces "github.com/stebinsabu13/ecommerce-api/pkg/repository/interface"
	services "github.com/stebinsabu13/ecommerce-api/pkg/usecase/interface"
	"github.com/stebinsabu13/ecommerce-api/pkg/utils"
)

type userUseCase struct {
	userRepo interfaces.UserRepository
}

func NewUserUseCase(repo interfaces.UserRepository) services.UserUseCase {
	return &userUseCase{
		userRepo: repo,
	}
}

//	func (c *userUseCase) FindAll(ctx context.Context) ([]domain.User, error) {
//		users, err := c.userRepo.FindAll(ctx)
//		return users, err
//	}
func (c *userUseCase) FindbyEmail(ctx context.Context, email string) (domain.User, error) {
	user, err := c.userRepo.FindbyEmail(ctx, email)
	return user, err
}

func (c *userUseCase) SignUpUser(ctx context.Context, user domain.User) error {
	err := c.userRepo.SignUpUser(ctx, user)
	return err
}

func (c *userUseCase) FindbyEmailorMobilenum(ctx context.Context, body utils.OtpLogin) (domain.User, error) {
	user, err := c.userRepo.FindbyEmailorMobilenum(ctx, body)
	return user, err
}

func (c *userUseCase) ShowDetails(ctx context.Context, id string) (utils.ResponseUsers, error) {
	user, err := c.userRepo.ShowDetails(ctx, id)
	return user, err
}

func (c *userUseCase) ShowAddress(ctx context.Context, id string) ([]utils.Address, error) {
	address, err := c.userRepo.ShowAddress(ctx, id)
	return address, err
}
