package handler

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
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
		buildStub      func(userUseCase mockUseCase.MockUserUseCase)
		expectedOutput utils.ResponseUsers
		expectedCode   int
		expectederr    error
	}{
		{
			name: "binding error",
			body: utils.BodyLogin{},
			buildStub: func(userUseCase mockUseCase.MockUserUseCase) {

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
			buildStub: func(userUseCase mockUseCase.MockUserUseCase) {
				userUseCase.EXPECT().FindbyEmail(
					gomock.Any(), "randomusr@gmail.com").Times(1).Return(
					utils.ResponseUsers{}, errors.New("invalid user"),
				)
			},
			expectedOutput: utils.ResponseUsers{},
			expectedCode:   401,
			expectederr:    errors.New("invalid user"),
		},
		{
			name: "valid user",
			body: utils.BodyLogin{
				Email:    "stebinsabu369@gmail.com",
				Password: "Stebin@333",
			},
			buildStub: func(userUseCase mockUseCase.MockUserUseCase) {
				userUseCase.EXPECT().FindbyEmail(
					gomock.Any(), "stebinsabu369@gmail.com").Times(1).Return(
					utils.ResponseUsers{
						ID:        2,
						FirstName: "Stebin",
						LastName:  "Sabu",
						Email:     "stebinsabu369@gmail.com",
						Password:  "Stebin@333",
						MobileNum: "9947650091",
						Block:     false,
						Verified:  true,
					}, nil,
				)
			},
			expectedOutput: utils.ResponseUsers{
				ID:        2,
				FirstName: "Stebin",
				LastName:  "Sabu",
				Email:     "stebinsabu369@gmail.com",
				Password:  "Stebin@333",
				MobileNum: "9947650091",
				Block:     false,
				Verified:  true,
			},
			expectedCode: 200,
			expectederr:  nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.buildStub(*userUseCase)
			gin.SetMode(gin.TestMode)
			engine := gin.Default()

			recorder := httptest.NewRecorder()

			// var body []byte

			// marshaling user data in the test case
			bodyjson := gin.H{
				"email":    tt.body.Email,
				"password": tt.body.Password,
			}
			b, err := json.Marshal(bodyjson)

			assert.NoError(t, err)
			url := "/user/login"

			req := httptest.NewRequest(http.MethodPost, url, bytes.NewBuffer(b))
			engine.POST(url, userHandler.LoginHandler)

			engine.ServeHTTP(recorder, req)

			var actual utils.ResponseUsers

			err = json.Unmarshal(recorder.Body.Bytes(), &actual)

			assert.NoError(t, err)

			assert.Equal(t, tt.expectedCode, recorder.Code)
			if tt.expectedCode == http.StatusOK {
				assert.Equal(t, gin.H{"Success": tt.expectedOutput}, actual)

			} else {
				assert.Equal(t, gin.H{"error": tt.expectedOutput}, actual)
			}
			// assert.Equal(t,tt.expectederr,)
		})
	}
}
