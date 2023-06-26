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

			// initialize a response recorder for capturing http  response
			recorder := httptest.NewRecorder()

			//url string for the endpoint
			url := "/user/login"

			// create a new route for testing
			engine.POST(url, userHandler.LoginHandler)

			// body is a slice of bytes. It is used for Marshaling data to json and passing to the request body
			// var body []byte

			// marshaling user data in the test case
			bodyjson := gin.H{
				"email":    tt.body.Email,
				"password": tt.body.Password,
			}
			b, err := json.Marshal(bodyjson)
			// validate no error occurred while marshaling data to json
			assert.NoError(t, err)

			// NewRequest returns a new incoming server Request, which we can pass to a http.Handler for testing
			req := httptest.NewRequest(http.MethodPost, url, bytes.NewBuffer(b)) // check what is buffer

			// req is a pointer to http.Request . With httptest.NewRequest we are mentioning the http method, endpoint and body
			engine.ServeHTTP(recorder, req)

			// actual will hold the actual reponse
			var actual utils.ResponseUsers

			// unmarshalling json data to response.Response format
			err = json.Unmarshal(recorder.Body.Bytes(), &actual)

			// validating no error occurred while unmarshalling json to response.Response struct
			assert.NoError(t, err)

			// validate expected status code and received status code are same
			assert.Equal(t, tt.expectedCode, recorder.Code)

			// validate expected response message and received response are same
			assert.Equal(t, gin.H{"Success": tt.expectedOutput}, actual)

		})
	}
}

func TestSignUp(t *testing.T) {

}
