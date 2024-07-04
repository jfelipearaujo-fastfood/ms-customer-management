package delete_account

import (
	"context"

	"github.com/go-playground/validator/v10"
	"github.com/jfelipearaujo-org/ms-customer-management/internal/shared/custom_error"
)

type DeleteAccountRequest struct {
	Id string `param:"id" validate:"required,uuid4"`
}

func (r *DeleteAccountRequest) Validate() error {
	validator := validator.New()

	if err := validator.Struct(r); err != nil {
		return custom_error.ErrRequestNotValid
	}

	return nil
}

type Service interface {
	Delete(ctx context.Context, request DeleteAccountRequest) error
}
