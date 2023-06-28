package handler

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/stebinsabu13/ecommerce-api/pkg/support"
	"github.com/stebinsabu13/ecommerce-api/pkg/usecase/mockUseCase"
	"github.com/stebinsabu13/ecommerce-api/pkg/utils"
	"github.com/stretchr/testify/assert"
)

func TestLoginHandler(t *testing.T) {
	ctrl := gomock.NewController(t)
	userUseCase := mockUseCase.NewMockUserUseCase(ctrl)
	userHandler := NewUserHandler(userUseCase, nil, nil)
	tests := []struct {
		name           string
		body           utils.BodyLogin
		buildStub      func(userUseCase mockUseCase.MockUserUseCase, body utils.BodyLogin)
		expectedOutput utils.ResponseUsers
		expectedCode   int
		expectederr    error
	}{
		{
			name: "binding error",
			body: utils.BodyLogin{},
			buildStub: func(userUseCase mockUseCase.MockUserUseCase, body utils.BodyLogin) {

			},
			expectedOutput: utils.ResponseUsers{},
			expectedCode:   http.StatusBadRequest,
			expectederr:    errors.New("failed to bind the required fields"),
		},
		{
			name: "Invalid user",
			body: utils.BodyLogin{
				Email:    "randomusr@gmail.com",
				Password: "random123",
			},
			buildStub: func(userUseCase mockUseCase.MockUserUseCase, body utils.BodyLogin) {
				userUseCase.EXPECT().FindbyEmail(
					gomock.Any(), body.Email).Times(1).Return(
					utils.ResponseUsers{}, errors.New("invalid user"),
				)
			},
			expectedOutput: utils.ResponseUsers{},
			expectedCode:   http.StatusUnauthorized,
			expectederr:    errors.New("invalid user"),
		},
		{
			name: "valid user and invalid password",
			body: utils.BodyLogin{
				Email:    "stebinsabu369@gmail.com",
				Password: "suh@1334",
			},
			buildStub: func(userUseCase mockUseCase.MockUserUseCase, body utils.BodyLogin) {
				hash, _ := support.HashPassword(body.Password)
				userUseCase.EXPECT().FindbyEmail(
					gomock.Any(), body.Email).Times(1).Return(
					utils.ResponseUsers{
						ID:          1,
						FirstName:   "Stebin",
						LastName:    "Sabu",
						Email:       "stebinsabu369@gmail.com",
						MobileNum:   "9947650091",
						Password:    hash,
						Block:       false,
						Verified:    true,
						ReferalCode: "jsjil",
					}, errors.New("invalid password"),
				)
			},
			expectedOutput: utils.ResponseUsers{},
			expectedCode:   http.StatusUnauthorized,
			expectederr:    errors.New("invalid password"),
		},
		// {
		// 	name: "Valid user and correct password",
		// 	body: utils.BodyLogin{
		// 		Email:    "stebinsabu369@gmail.com",
		// 		Password: "Stebin@333",
		// 	},
		// 	buildStub: func(userUseCase mockUseCase.MockUserUseCase, body utils.BodyLogin) {
		// 		hash, _ := support.HashPassword(body.Password)
		// 		userUseCase.EXPECT().FindbyEmail(
		// 			gomock.Any(), body.Email).Times(1).Return(
		// 			utils.ResponseUsers{
		// 				ID:          1,
		// 				FirstName:   "Stebin",
		// 				LastName:    "Sabu",
		// 				Email:       "stebinsabu369@gmail.com",
		// 				MobileNum:   "9947650091",
		// 				Password:    hash,
		// 				Block:       false,
		// 				Verified:    true,
		// 				ReferalCode: "jsjil",
		// 			}, nil,
		// 		)
		// 	},
		// 	expectedOutput: utils.ResponseUsers{
		// 		ID:          1,
		// 		FirstName:   "Stebin",
		// 		LastName:    "Sabu",
		// 		Email:       "stebinsabu369@gmail.com",
		// 		MobileNum:   "9947650091",
		// 		Password:    gomock.Any().String(),
		// 		Block:       false,
		// 		Verified:    true,
		// 		ReferalCode: "jsjil",
		// 	},
		// 	expectedCode: http.StatusOK,
		// 	expectederr:  nil,
		// },
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.buildStub(*userUseCase, tt.body)
			gin.SetMode(gin.TestMode)
			engine := gin.Default()

			recorder := httptest.NewRecorder()

			var body []byte
			var err error

			// marshaling user data in the test case
			// bodyjson := gin.H{
			// 	"email":    tt.body.Email,
			// 	"password": tt.body.Password,
			// }
			body, err = json.Marshal(tt.body)

			assert.NoError(t, err)
			url := "/user/login"

			req := httptest.NewRequest(http.MethodPost, url, bytes.NewBuffer(body))
			engine.POST(url, userHandler.LoginHandler)

			engine.ServeHTTP(recorder, req)

			var actual utils.ResponseUsers

			err = json.Unmarshal(recorder.Body.Bytes(), &actual)

			assert.NoError(t, err)
			fmt.Println("The output is", actual)

			assert.Equal(t, tt.expectedCode, recorder.Code)

			// validating expected data and received are same. If not test will fail
			if !reflect.DeepEqual(tt.expectedOutput, actual) {
				t.Errorf("got %v, but want %v", actual, tt.expectedOutput)
			}
		})
	}
}
