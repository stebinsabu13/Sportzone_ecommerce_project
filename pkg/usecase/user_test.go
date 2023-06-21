package usecase

import (
	"context"
	"errors"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stebinsabu13/ecommerce-api/pkg/domain"
	"github.com/stebinsabu13/ecommerce-api/pkg/repository/mockRepo"
	"github.com/stretchr/testify/assert"
)

func TestFindbyEmail(t *testing.T) {
	ctrl := gomock.NewController(t)
	userRepo := mockRepo.NewMockUserRepository(ctrl)
	userUseCase := NewUserUseCase(userRepo)
	tests := []struct {
		name           string
		input          string
		expectedOutput domain.User
		buildStub      func(userRepo mockRepo.MockUserRepository)
		expectederr    error
	}{
		{
			name:  "User exsists",
			input: "stebinsabu369@gmail.com",
			expectedOutput: domain.User{
				FirstName: "Stebin",
				LastName:  "Sabu",
				Email:     "stebinsabu369@gmail.com",
				Password:  "Stebin@333",
				MobileNum: "9947650091",
				Block:     false,
				Verified:  true,
			},
			buildStub: func(userRepo mockRepo.MockUserRepository) {
				userRepo.EXPECT().FindbyEmail(
					gomock.Any(), "stebinsabu369@gmail.com").Times(1).Return(domain.User{
					FirstName: "Stebin",
					LastName:  "Sabu",
					Email:     "stebinsabu369@gmail.com",
					Password:  "Stebin@333",
					MobileNum: "9947650091",
					Block:     false,
					Verified:  true,
				}, nil)
			},
			expectederr: nil,
		},
		{
			name:           "not exsist user",
			input:          "notexsist@gmail.com",
			expectedOutput: domain.User{},
			buildStub: func(userRepo mockRepo.MockUserRepository) {
				userRepo.EXPECT().FindbyEmail(
					gomock.Any(), "notexsist@gmail.com").Times(1).Return(
					domain.User{},
					errors.New("not exsist user"),
				)
			},
			expectederr: errors.New("not exsist user"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.buildStub(*userRepo)
			actualOutput, actualErr := userUseCase.FindbyEmail(context.TODO(), tt.input)
			assert.Equal(t, tt.expectedOutput, actualOutput)
			assert.Equal(t, tt.expectederr, actualErr)
		})
	}
}
