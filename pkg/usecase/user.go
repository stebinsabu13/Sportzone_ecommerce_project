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
func (c *userUseCase) FindbyUserID(ctx context.Context, id uint) (domain.User, error) {
	return c.userRepo.FindbyUserID(ctx, id)
}
func (c *userUseCase) SignUpUser(ctx context.Context, user domain.User) (string, error) {
	mobile_num, err := c.userRepo.SignUpUser(ctx, user)
	return mobile_num, err
}

func (c *userUseCase) FindbyEmailorMobilenum(ctx context.Context, body utils.OtpLogin) (domain.User, error) {
	user, err := c.userRepo.FindbyEmailorMobilenum(ctx, body)
	return user, err
}

func (c *userUseCase) ShowDetails(ctx context.Context, id uint) (utils.ResponseUsers, error) {
	user, err := c.userRepo.ShowDetails(ctx, id)
	return user, err
}

func (c *userUseCase) ShowAddress(ctx context.Context, id uint) ([]utils.Address, error) {
	address, err := c.userRepo.ShowAddress(ctx, id)
	return address, err
}

func (c *userUseCase) UpdateVerify(ctx context.Context, number string) error {
	err := c.userRepo.UpdateVerify(ctx, number)
	return err
}

func (c *userUseCase) AddAddress(ctx context.Context, address domain.Address) error {
	err := c.userRepo.AddAddress(ctx, address)
	return err
}

func (c *userUseCase) EditProfile(ctx context.Context, profile utils.EditProfileReq, id uint) error {
	return c.userRepo.EditProfile(ctx, profile, id)
}
func (c *userUseCase) ChangePassword(ctx context.Context, newpassword string, mobile_num string) error {
	return c.userRepo.ChangePassword(ctx, newpassword, mobile_num)
}

func (c *userUseCase) ListAllCategories(ctx context.Context) ([]utils.ResponseCategory, error) {
	categories, err := c.userRepo.ListAllCategories(ctx)
	return categories, err
}
