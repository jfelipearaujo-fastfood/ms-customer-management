package delete_account

import (
	"context"

	"github.com/jfelipearaujo-org/ms-customer-management/internal/repository/customer"
)

type service struct {
	repository customer.Repository
}

func NewService(repository customer.Repository) Service {
	return &service{
		repository: repository,
	}
}

func (s *service) Delete(ctx context.Context, request DeleteAccountRequest) error {
	if err := request.Validate(); err != nil {
		return err
	}

	return s.repository.Delete(ctx, request.Id)
}
