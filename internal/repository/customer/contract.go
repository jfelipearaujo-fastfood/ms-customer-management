package customer

import (
	"context"

	"github.com/jfelipearaujo-org/ms-customer-management/internal/entity"
)

type Repository interface {
	Get(ctx context.Context, id string) (entity.Customer, error)
	Delete(ctx context.Context, id string) error
}
