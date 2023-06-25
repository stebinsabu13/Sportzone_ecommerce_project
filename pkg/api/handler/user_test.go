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

// func TestLoginHandler(t *testing.T) {
// 	ctrl := gomock.NewController(t)
// 	userUseCase := mockUseCase.NewMockUserUseCase(ctrl)
// 	userHandler := NewUserHandler(userUseCase, nil, nil)
// 	tests := []struct {
// 		name           string
// 		body           utils.BodyLogin
// 		buildStub      func(userUseCase mockUseCase.MockUserUseCase)
// 		expectedOutput utils.ResponseUsers
// 		expectedCode   int
// 		expectederr    error
// 	}{
// 		{
// 			name: "valid user",
// 			body: utils.BodyLogin{
// 				Email:    "stebinsabu369@gmail.com",
// 				Password: "Stebin@333",
// 			},
// 			buildStub: func(userUseCase mockUseCase.MockUserUseCase) {
// 				userUseCase.EXPECT().FindbyEmail(
// 					gomock.Any(), "stebinsabu369@gmail.com").Times(1).Return(
// 					utils.ResponseUsers{
// 						ID:        2,
// 						FirstName: "Stebin",
// 						LastName:  "Sabu",
// 						Email:     "stebinsabu369@gmail.com",
// 						Password:  "Stebin@333",
// 						MobileNum: "9947650091",
// 						Block:     false,
// 					}, nil,
// 				)
// 			},
// 			expectedOutput: utils.ResponseUsers{
// 				ID:        2,
// 				FirstName: "Stebin",
// 				LastName:  "Sabu",
// 				Email:     "stebinsabu369@gmail.com",
// 				Password:  "Stebin@333",
// 				MobileNum: "9947650091",
// 				Block:     false,
// 			},
// 			expectedCode: 200,
// 			expectederr:  nil,
// 		},
// 		{
// 			name: "Invalid user",
// 			body: utils.BodyLogin{
// 				Email:    "randomusr@gmail.com",
// 				Password: "random123",
// 			},
// 			buildStub: func(userUseCase mockUseCase.MockUserUseCase) {
// 				userUseCase.EXPECT().FindbyEmail(
// 					gomock.Any(), "randomusr@gmail.com").Times(1).Return(
// 					utils.ResponseUsers{}, errors.New("invalid user"),
// 				)
// 			},
// 			expectedOutput: utils.ResponseUsers{},
// 			expectedCode:   401,
// 			expectederr:    errors.New("invalid user"),
// 		},
// 	}
// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			tt.buildStub(*userUseCase)
// 			gin.SetMode(gin.TestMode)
// 			engine := gin.Default()

// 			// initialize a response recorder for capturing http  response
// 			recorder := httptest.NewRecorder()

// 			//url string for the endpoint
// 			url := "/user/login"

// 			// create a new route for testing
// 			engine.POST(url, userHandler.LoginHandler)

// 			// body is a slice of bytes. It is used for Marshaling data to json and passing to the request body
// 			var body []byte

// 			// marshaling user data in the test case
// 			body, err := json.Marshal(tt.body)

// 			// validate no error occurred while marshaling data to json
// 			assert.NoError(t, err)

// 			// NewRequest returns a new incoming server Request, which we can pass to a http.Handler for testing
// 			req := httptest.NewRequest(http.MethodPost, url, bytes.NewBuffer(body)) // check what is buffer

// 			// req is a pointer to http.Request . With httptest.NewRequest we are mentioning the http method, endpoint and body
// 			engine.ServeHTTP(recorder, req)

// 			// actual will hold the actual reponse
// 			var actual utils.ResponseUsers

// 			// unmarshalling json data to response.Response format
// 			err = json.Unmarshal(recorder.Body.Bytes(), &actual)

// 			// validating no error occurred while unmarshalling json to response.Response struct
// 			assert.NoError(t, err)

// 			// validate expected status code and received status code are same
// 			assert.Equal(t, tt.expectedCode, recorder.Code)

// 			// validate expected response message and received response are same
// 			assert.Equal(t, tt.expectedOutput, actual)

// 		})
// 	}
// }

func TestLoginHandler(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUserUseCase := mockUseCase.NewMockUserUseCase(ctrl)
	handler := &UserHandler{
		userUseCase: mockUserUseCase,
	}
	gin.SetMode(gin.TestMode)
	// Create a new Gin router and set the LoginHandler as the route handler
	router := gin.Default()
	router.POST("/login", handler.LoginHandler)

	// Test case 1: Successful login
	t.Run("Successful login", func(t *testing.T) {
		// Mock the userUseCase behavior
		mockUser := utils.ResponseUsers{
			ID:        1,
			FirstName: "test",
			LastName:  "example",
			MobileNum: "9947650091",
			Block:     false,
			Email:     "test@example.com",
			Password:  "hashedpassword",
			Verified:  true,
		}
		mockUserUseCase.EXPECT().FindbyEmail(gomock.Any(), "test@example.com").Return(mockUser, nil)

		// Create a JSON request body
		requestBody := utils.BodyLogin{
			Email:    "test@example.com",
			Password: "password123",
		}
		bodyJSON := gin.H{
			"email":    requestBody.Email,
			"password": requestBody.Password,
		}

		// Perform a POST request to the /login endpoint
		w := performRequest(router, "POST", "/user/login", bodyJSON)

		// Assertions
		assert.Equal(t, http.StatusOK, w.Code)
		assert.Equal(t, gin.H{"Success": mockUser}, responseJSON(w))
	})

	// Test case 2: Invalid password
	t.Run("Invalid password", func(t *testing.T) {
		// Mock the userUseCase behavior
		mockUser := utils.ResponseUsers{
			ID:        1,
			FirstName: "test",
			LastName:  "example",
			MobileNum: "9947650091",
			Block:     false,
			Email:     "test@example.com",
			Password:  "hashedpassword",
		}
		mockUserUseCase.EXPECT().FindbyEmail(gomock.Any(), "test@example.com").Return(mockUser, errors.New("invalid password"))

		// Create a JSON request body with an invalid password
		requestBody := utils.BodyLogin{
			Email:    "test@example.com",
			Password: "wrongpassword",
		}
		bodyJSON := gin.H{
			"email":    requestBody.Email,
			"password": requestBody.Password,
		}

		// Perform a POST request to the /login endpoint
		w := performRequest(router, "POST", "/user/login", bodyJSON)

		// Assertions
		assert.Equal(t, http.StatusUnauthorized, w.Code)
		assert.Equal(t, gin.H{"error": "invalid password"}, responseJSON(w))
	})

	// More test cases can be added as needed

}

// Helper function to perform a request to the Gin router
func performRequest(router http.Handler, method, path string, body interface{}) *httptest.ResponseRecorder {
	b, _ := json.Marshal(body)
	req, _ := http.NewRequest(method, path, bytes.NewReader(b))
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	return w
}

// Helper function to extract the JSON response body from the response recorder
func responseJSON(w *httptest.ResponseRecorder) gin.H {
	var response gin.H
	_ = json.Unmarshal(w.Body.Bytes(), &response)
	return response
}
