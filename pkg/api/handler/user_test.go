package handler

import (
	"errors"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stebinsabu13/ecommerce-api/pkg/domain"
	"github.com/stebinsabu13/ecommerce-api/pkg/usecase/mockUseCase"
	"github.com/stebinsabu13/ecommerce-api/pkg/utils"
)

func TestLoginHandler(t *testing.T) {
	ctrl := gomock.NewController(t)
	userUseCase := mockUseCase.NewMockUserUseCase(ctrl)
	_ = NewUserHandler(userUseCase, nil, nil)
	_ = []struct {
		name           string
		input          utils.BodyLogin
		buildStub      func(userUseCase mockUseCase.MockUserUseCase)
		expectedOutput domain.User
		expectederr    error
	}{
		{
			name: "valid user",
			input: utils.BodyLogin{
				Email:    "stebinsabu369@gmail.com",
				Password: "Stebin@333",
			},
			buildStub: func(userUseCase mockUseCase.MockUserUseCase) {
				userUseCase.EXPECT().FindbyEmail(
					gomock.Any(), "stebinsabu369@gmail.com").Times(1).Return(
					domain.User{
						FirstName: "Stebin",
						LastName:  "Sabu",
						MobileNum: "9947650091",
						Email:     "stebinsabu369@gmail.com",
						Block:     false,
						Verified:  true,
						Password:  "Stebin@333",
					},
				)
			},
			expectedOutput: domain.User{
				FirstName: "Stebin",
				LastName:  "Sabu",
				Email:     "stebinsabu369@gmail.com",
				Password:  "Stebin@333",
				MobileNum: "9947650091",
				Block:     false,
				Verified:  true,
			},
			expectederr: nil,
		},
		{
			name: "Invalid user",
			input: utils.BodyLogin{
				Email:    "randomusr@gmail.com",
				Password: "random123",
			},
			buildStub: func(userUseCase mockUseCase.MockUserUseCase) {

			},
			expectedOutput: domain.User{},
			expectederr:    errors.New("invalid user"),
		},
	}
}
