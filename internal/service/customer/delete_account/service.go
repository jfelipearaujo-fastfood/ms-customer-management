package delete_account

import (
	"context"

	"github.com/jfelipearaujo-org/ms-customer-management/internal/entity"
	"github.com/jfelipearaujo-org/ms-customer-management/internal/repository/customer"
	"github.com/jfelipearaujo-org/ms-customer-management/internal/repository/delete_request"
	"github.com/jfelipearaujo-org/ms-customer-management/internal/shared/custom_error"
)

type service struct {
	customerRepository      customer.Repository
	deleteRequestRepository delete_request.Repository
}

func NewService(
	customerRepository customer.Repository,
	deleteRequestRepository delete_request.Repository,
) Service {
	return &service{
		customerRepository:      customerRepository,
		deleteRequestRepository: deleteRequestRepository,
	}
}

func (s *service) Delete(ctx context.Context, request DeleteAccountRequest) error {
	if err := request.Validate(); err != nil {
		return err
	}

	if _, err := s.customerRepository.Get(ctx, request.Id); err != nil {
		return err
	}

	existingDeleteRequest, err := s.deleteRequestRepository.GetByCustomerId(ctx, request.Id)
	if err != nil && err != custom_error.ErrDeletionRequestNotFound {
		return err
	}

	if existingDeleteRequest.Id != "" {
		return custom_error.ErrDeletionRequestAlreadyCreated
	}

	deleteRequest := entity.NewDeleteRequest(request.Id,
		request.Name,
		request.Address,
		request.Phone)

	return s.deleteRequestRepository.Create(ctx, deleteRequest)
}
