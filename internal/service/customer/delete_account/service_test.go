package delete_account_test

import (
	"context"
	"testing"

	"github.com/jfelipearaujo-org/ms-customer-management/internal/repository/customer"
	"github.com/jfelipearaujo-org/ms-customer-management/internal/service/customer/delete_account"
	"github.com/stretchr/testify/assert"
)

func TestService_Delete(t *testing.T) {
	t.Run("Should delete a customer", func(t *testing.T) {
		// Arrange
		ctx := context.Background()

		repo := customer.NewMockRepository(t)

		repo.On("Delete", ctx, "733f1ba6-1f62-4495-bf33-6f181fdf1030").
			Return(nil)

		service := delete_account.NewService(repo)

		// Act
		err := service.Delete(ctx, delete_account.DeleteAccountRequest{
			Id: "733f1ba6-1f62-4495-bf33-6f181fdf1030",
		})

		// Assert
		assert.NoError(t, err)
		repo.AssertExpectations(t)
	})

	t.Run("Should return an error", func(t *testing.T) {
		// Arrange
		ctx := context.Background()

		repo := customer.NewMockRepository(t)

		service := delete_account.NewService(repo)

		// Act
		err := service.Delete(ctx, delete_account.DeleteAccountRequest{
			Id: "id",
		})

		// Assert
		assert.Error(t, err)
		repo.AssertExpectations(t)
	})
}
