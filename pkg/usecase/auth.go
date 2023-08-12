package usecase

import (
	"context"
	"fmt"
	"strconv"

	interfaces "github.com/stebinsabu13/ecommerce-api/pkg/repository/interface"
	"github.com/stebinsabu13/ecommerce-api/pkg/support"
	services "github.com/stebinsabu13/ecommerce-api/pkg/usecase/interface"
	"github.com/stebinsabu13/ecommerce-api/pkg/utils"
)

type authUseCase struct {
	userRepo interfaces.UserRepository
}

func NewAuthUseCase(repo interfaces.UserRepository) services.AuthUseCase {
	return &authUseCase{
		userRepo: repo,
	}
}

func (c *authUseCase) GoogleLogin(ctx context.Context, user utils.BodySignUpuser) (uint, error) {

	existUser, err := c.userRepo.FindbyEmail(ctx, user.Email)
	if err != nil {
		return existUser.ID, fmt.Errorf("failed to get user details with given email \nerror:%v", err.Error())
	}

	if existUser.ID != 0 {
		return existUser.ID, nil
	}

	user.ReferalCode = support.ReferalCodeGenerator()
	userID, err := c.userRepo.SignUpUser(ctx, user)
	userid, _ := strconv.Atoi(userID)
	if err != nil {
		return uint(userid), fmt.Errorf("failed to save user details \nerror:%v", err.Error())
	}

	return uint(userid), nil
}
