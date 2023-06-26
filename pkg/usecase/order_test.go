package usecase

import (
	"errors"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stebinsabu13/ecommerce-api/pkg/domain"
	"github.com/stebinsabu13/ecommerce-api/pkg/repository/mockRepo"
	"github.com/stebinsabu13/ecommerce-api/pkg/utils"
	"github.com/stretchr/testify/assert"
)

// func TestAddtoOrders(t *testing.T) {
// 	ctrl := gomock.NewController(t)
// 	orderRepo := mockRepo.NewMockOrderRepository(ctrl)
// 	cartRepo := mockRepo.NewMockCartRepository(ctrl)
// 	orderusecase := NewOrderUseCase(orderRepo, cartRepo)
// 	tests:=[]struct{
// 		name string
// 		addressid uint
// 		paymentid uint
// 		userid uint
// 		couponid *uint
// 		expectedOutput error
// 		buildStub func(cartrepo mockRepo.MockCartRepository)
// 		buildStubForItems func(orderrepo mockRepo.MockOrderRepository)
// 	}{
// 		{
// 			name: "order placed",
// 			addressid: 1,
// 			paymentid: 1,
// 			userid: 1,
// 			couponid: 1,
// 			expectedOutput: nil,
// 			buildStub: func( cartrepo mockRepo.MockCartRepository) {
// 				cartrepo.EXPECT().FindCartById(1).Times(1).Return(
// 					domain.Cart{
// 						ID: 1,
// 						UserID: 1,
// 						GrandTotal: 1000,
// 					},nil,
// 				)
// 			},
// 			bu
// 		}
// 	}
// }

func TestAddtoOrders(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockCartRepo := mockRepo.NewMockCartRepository(ctrl)
	mockOrderRepo := mockRepo.NewMockOrderRepository(ctrl)

	// Test cases
	testCases := []struct {
		name           string
		findCartErr    error
		findItemsErr   error
		addToOrdersErr error
		expectedErr    error
	}{
		{
			name:           "Success",
			findCartErr:    nil,
			findItemsErr:   nil,
			addToOrdersErr: nil,
			expectedErr:    nil,
		},
		{
			name:        "Error finding cart",
			findCartErr: errors.New("error finding cart"),
			expectedErr: errors.New("error finding cart"),
		},
		{
			name:         "Error finding cart items",
			findCartErr:  nil,
			findItemsErr: errors.New("error finding cart items"),
			expectedErr:  errors.New("error finding cart items"),
		},
		{
			name:           "Error adding to orders",
			findCartErr:    nil,
			findItemsErr:   nil,
			addToOrdersErr: errors.New("error adding to orders"),
			expectedErr:    errors.New("error adding to orders"),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Set up mocks and expectations
			mockCartRepo.EXPECT().FindCartById(gomock.Any()).Return(domain.Cart{}, tc.findCartErr)
			mockOrderRepo.EXPECT().Findcartitems(gomock.Any()).Return([]utils.ResCartItems{}, tc.findItemsErr)
			mockOrderRepo.EXPECT().AddtoOrders(gomock.Any(), gomock.Any()).Return(tc.addToOrdersErr)

			// Create the use case instance with mocks
			orderUseCase := NewOrderUseCase(mockOrderRepo, mockCartRepo)

			// Call the function under test
			err := orderUseCase.AddtoOrders(1, 1, 1, nil)

			// Check the error result
			assert.Equal(t, tc.expectedErr, err)
		})
	}
}
