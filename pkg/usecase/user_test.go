package usecase

import (
	"context"
	"errors"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stebinsabu13/ecommerce-api/pkg/repository/mockRepo"
	"github.com/stebinsabu13/ecommerce-api/pkg/utils"
	"github.com/stretchr/testify/assert"
)

func TestFindbyEmail(t *testing.T) {
	ctrl := gomock.NewController(t)
	userRepo := mockRepo.NewMockUserRepository(ctrl)
	userUseCase := NewUserUseCase(userRepo)
	tests := []struct {
		name           string
		input          string
		expectedOutput utils.ResponseUsers
		buildStub      func(userRepo mockRepo.MockUserRepository)
		expectederr    error
	}{
		{
			name:  "User exsists",
			input: "stebinsabu369@gmail.com",
			expectedOutput: utils.ResponseUsers{
				ID:        1,
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
					gomock.Any(), "stebinsabu369@gmail.com").Times(1).Return(utils.ResponseUsers{
					ID:        1,
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
			expectedOutput: utils.ResponseUsers{},
			buildStub: func(userRepo mockRepo.MockUserRepository) {
				userRepo.EXPECT().FindbyEmail(
					gomock.Any(), "notexsist@gmail.com").Times(1).Return(
					utils.ResponseUsers{},
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

func TestSignUpUser(t *testing.T) {
	ctrl := gomock.NewController(t)
	userRepo := mockRepo.NewMockUserRepository(ctrl)
	userUseCase := NewUserUseCase(userRepo)
	tests := []struct {
		name           string
		input          utils.BodySignUpuser
		expectedOutput string
		buildStub      func(userRepo mockRepo.MockUserRepository)
		expectederr    error
	}{
		{
			name: "user already exsists",
			input: utils.BodySignUpuser{
				FirstName:   "already",
				LastName:    "exsists",
				Email:       "exsistinguser@gmail.com",
				MobileNum:   "9947650091",
				Password:    "exsists@333",
				ReferalCode: "jd34f",
			},
			expectedOutput: "",
			buildStub: func(userRepo mockRepo.MockUserRepository) {
				userRepo.EXPECT().FindbyEmail(
					gomock.Any(), "exsistinguser@gmail.com").Times(1).Return(
					utils.ResponseUsers{
						FirstName:   "already",
						LastName:    "exsists",
						Email:       "exsistinguser@gmail.com",
						MobileNum:   "9947650091",
						Password:    "exsists@333",
						ReferalCode: "jd34f",
					}, nil,
				)
			},
			expectederr: errors.New("user already exsists"),
		},
		{
			name: "not exsisting user",
			input: utils.BodySignUpuser{
				FirstName:   "not",
				LastName:    "exsists",
				Email:       "notexsistinguser@gmail.com",
				MobileNum:   "8376423610",
				Password:    "notexsists@333",
				ReferalCode: "jidj3j",
			},
			expectedOutput: "8376423610",
			buildStub: func(userRepo mockRepo.MockUserRepository) {
				userRepo.EXPECT().FindbyEmail(
					gomock.Any(), "notexsistinguser@gmail.com").Times(1).Return(
					utils.ResponseUsers{}, errors.New("invalid email"),
				)
				userRepo.EXPECT().SignUpUser(
					gomock.Any(), gomock.Eq(utils.BodySignUpuser{
						FirstName:   "not",
						LastName:    "exsists",
						Email:       "notexsistinguser@gmail.com",
						MobileNum:   "8376423610",
						Password:    gomock.Any().String(),
						ReferalCode: gomock.Any().String(),
					})).Times(1).Return(
					"8376423610", nil,
				)
			},
			expectederr: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.buildStub(*userRepo)
			actualOutput, actualErr := userUseCase.SignUpUser(context.TODO(), tt.input)
			assert.Equal(t, tt.expectedOutput, actualOutput)
			assert.Equal(t, tt.expectederr, actualErr)
		})
	}
}
