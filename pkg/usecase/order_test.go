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

func TestAddtoOrders(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockCartRepo := mockRepo.NewMockCartRepository(ctrl)
	mockOrderRepo := mockRepo.NewMockOrderRepository(ctrl)
	orderUseCase := NewOrderUseCase(mockOrderRepo, mockCartRepo)

	// Test cases
	testCases := []struct {
		name           string
		buildStub      func(mockcart mockRepo.MockCartRepository, mockorder mockRepo.MockOrderRepository)
		findCartErr    error
		findItemsErr   error
		addToOrdersErr error
		expectedErr    error
	}{
		{
			name: "Success",
			buildStub: func(mockcart mockRepo.MockCartRepository, mockorder mockRepo.MockOrderRepository) {
				mockcart.EXPECT().FindCartById(uint(1)).Times(1).Return(
					domain.Cart{
						ID:         1,
						UserID:     1,
						GrandTotal: 100,
					}, nil,
				)
				mockorder.EXPECT().Findcartitems(uint(1)).Times(1).Return(
					[]utils.ResCartItems{
						{
							CartID:          1,
							ProductDetailID: 1,
							Quantity:        2,
						},
					}, nil,
				)
				mockorder.EXPECT().AddtoOrders(gomock.Any(), gomock.Any()).Times(1).Return(
					nil,
				)
			},
			findCartErr:    nil,
			findItemsErr:   nil,
			addToOrdersErr: nil,
			expectedErr:    nil,
		},
		{
			name: "Error finding cart",
			buildStub: func(mockcart mockRepo.MockCartRepository, mockorder mockRepo.MockOrderRepository) {
				mockcart.EXPECT().FindCartById(uint(1)).Times(1).Return(
					domain.Cart{}, errors.New("error finding cart"),
				)
			},
			findCartErr: errors.New("error finding cart"),
			expectedErr: errors.New("error finding cart"),
		},
		{
			name: "Error finding cart items",
			buildStub: func(mockcart mockRepo.MockCartRepository, mockorder mockRepo.MockOrderRepository) {
				mockcart.EXPECT().FindCartById(uint(1)).Times(1).Return(
					domain.Cart{
						ID:         1,
						UserID:     1,
						GrandTotal: 100,
					}, nil,
				)
				mockorder.EXPECT().Findcartitems(uint(1)).Times(1).Return(
					[]utils.ResCartItems{}, errors.New("error finding cart items"),
				)
			},
			findCartErr:  nil,
			findItemsErr: errors.New("error finding cart items"),
			expectedErr:  errors.New("error finding cart items"),
		},
		{
			name: "Error adding to orders",
			buildStub: func(mockcart mockRepo.MockCartRepository, mockorder mockRepo.MockOrderRepository) {
				mockcart.EXPECT().FindCartById(uint(1)).Times(1).Return(
					domain.Cart{
						ID:         1,
						UserID:     1,
						GrandTotal: 100,
					}, nil,
				)
				mockorder.EXPECT().Findcartitems(uint(1)).Times(1).Return(
					[]utils.ResCartItems{
						{
							CartID:          1,
							ProductDetailID: 1,
							Quantity:        2,
						},
					}, nil,
				)
				mockorder.EXPECT().AddtoOrders(gomock.Any(), gomock.Any()).Times(1).Return(
					errors.New("error adding to orders"),
				)
			},
			findCartErr:    nil,
			findItemsErr:   nil,
			addToOrdersErr: errors.New("error adding to orders"),
			expectedErr:    errors.New("error adding to orders"),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc.buildStub(*mockCartRepo, *mockOrderRepo)

			err := orderUseCase.AddtoOrders(1, 1, 1, nil)
			assert.Equal(t, tc.expectedErr, err)
		})
	}
}
