package delete_request

import (
	"context"

	"github.com/jfelipearaujo-org/ms-customer-management/internal/entity"
)

type Repository interface {
	GetByCustomerId(ctx context.Context, customerId string) (entity.DeletionRequest, error)
	Create(ctx context.Context, request entity.DeletionRequest) error
}
