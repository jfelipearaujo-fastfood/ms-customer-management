package delete_account_test

import (
	"context"
	"testing"

	"github.com/jfelipearaujo-org/ms-customer-management/internal/entity"
	"github.com/jfelipearaujo-org/ms-customer-management/internal/repository/customer"
	"github.com/jfelipearaujo-org/ms-customer-management/internal/repository/delete_request"
	"github.com/jfelipearaujo-org/ms-customer-management/internal/service/customer/delete_account"
	"github.com/jfelipearaujo-org/ms-customer-management/internal/shared/custom_error"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestService_Delete(t *testing.T) {
	t.Run("Should delete a customer", func(t *testing.T) {
		// Arrange
		ctx := context.Background()

		customerRepository := customer.NewMockRepository(t)
		deleteRequestRepository := delete_request.NewMockRepository(t)

		customerRepository.On("Get", ctx, "733f1ba6-1f62-4495-bf33-6f181fdf1030").
			Return(entity.Customer{}, nil)

		deleteRequestRepository.On("GetByCustomerId", ctx, mock.Anything).
			Return(entity.DeletionRequest{}, nil)

		deleteRequestRepository.On("Create", ctx, mock.Anything).
			Return(nil)

		service := delete_account.NewService(customerRepository, deleteRequestRepository)

		// Act
		err := service.Delete(ctx, delete_account.DeleteAccountRequest{
			Id:      "733f1ba6-1f62-4495-bf33-6f181fdf1030",
			Name:    "John Doe",
			Address: "Av. Brasil, 1000",
			Phone:   "1122334455",
		})

		// Assert
		assert.NoError(t, err)
		customerRepository.AssertExpectations(t)
		deleteRequestRepository.AssertExpectations(t)
	})

	t.Run("Should return an error when customer is not found", func(t *testing.T) {
		// Arrange
		ctx := context.Background()

		customerRepository := customer.NewMockRepository(t)
		deleteRequestRepository := delete_request.NewMockRepository(t)

		customerRepository.On("Get", ctx, "733f1ba6-1f62-4495-bf33-6f181fdf1030").
			Return(entity.Customer{}, custom_error.ErrCustomerNotFound)

		service := delete_account.NewService(customerRepository, deleteRequestRepository)

		// Act
		err := service.Delete(ctx, delete_account.DeleteAccountRequest{
			Id:      "733f1ba6-1f62-4495-bf33-6f181fdf1030",
			Name:    "John Doe",
			Address: "Av. Brasil, 1000",
			Phone:   "1122334455",
		})

		// Assert
		assert.Error(t, err)
		customerRepository.AssertExpectations(t)
		deleteRequestRepository.AssertExpectations(t)
	})

	t.Run("Should return an error when delete request already created", func(t *testing.T) {
		// Arrange
		ctx := context.Background()

		customerRepository := customer.NewMockRepository(t)
		deleteRequestRepository := delete_request.NewMockRepository(t)

		customerRepository.On("Get", ctx, "733f1ba6-1f62-4495-bf33-6f181fdf1030").
			Return(entity.Customer{}, nil)

		deleteRequestRepository.On("GetByCustomerId", ctx, mock.Anything).
			Return(entity.DeletionRequest{
				Id: "id",
			}, nil)

		service := delete_account.NewService(customerRepository, deleteRequestRepository)

		// Act
		err := service.Delete(ctx, delete_account.DeleteAccountRequest{
			Id:      "733f1ba6-1f62-4495-bf33-6f181fdf1030",
			Name:    "John Doe",
			Address: "Av. Brasil, 1000",
			Phone:   "1122334455",
		})

		// Assert
		assert.Error(t, err)
		customerRepository.AssertExpectations(t)
		deleteRequestRepository.AssertExpectations(t)
	})

	t.Run("Should return an error when try to create a deletion request", func(t *testing.T) {
		// Arrange
		ctx := context.Background()

		customerRepository := customer.NewMockRepository(t)
		deleteRequestRepository := delete_request.NewMockRepository(t)

		customerRepository.On("Get", ctx, "733f1ba6-1f62-4495-bf33-6f181fdf1030").
			Return(entity.Customer{}, nil)

		deleteRequestRepository.On("GetByCustomerId", ctx, mock.Anything).
			Return(entity.DeletionRequest{}, nil)

		deleteRequestRepository.On("Create", ctx, mock.Anything).
			Return(custom_error.ErrRequestNotValid)

		service := delete_account.NewService(customerRepository, deleteRequestRepository)

		// Act
		err := service.Delete(ctx, delete_account.DeleteAccountRequest{
			Id:      "733f1ba6-1f62-4495-bf33-6f181fdf1030",
			Name:    "John Doe",
			Address: "Av. Brasil, 1000",
			Phone:   "1122334455",
		})

		// Assert
		assert.Error(t, err)
		customerRepository.AssertExpectations(t)
		deleteRequestRepository.AssertExpectations(t)
	})
}
