package delete_account

import (
	"net/http"

	"github.com/jfelipearaujo-org/ms-customer-management/internal/service/customer/delete_account"
	"github.com/jfelipearaujo-org/ms-customer-management/internal/shared/custom_error"
	"github.com/labstack/echo/v4"
)

type Handler struct {
	service delete_account.Service
}

func NewHandler(service delete_account.Service) *Handler {
	return &Handler{
		service: service,
	}
}

func (h *Handler) Handle(ctx echo.Context) error {
	var request delete_account.DeleteAccountRequest

	if err := ctx.Bind(&request); err != nil {
		return custom_error.NewHttpAppError(http.StatusBadRequest, "invalid request", err)
	}

	context := ctx.Request().Context()

	if err := h.service.Delete(context, request); err != nil {
		if custom_error.IsBusinessErr(err) {
			return custom_error.NewHttpAppErrorFromBusinessError(err)
		}

		return custom_error.NewHttpAppError(http.StatusInternalServerError, "internal error deleting customer", err)
	}

	return ctx.NoContent(http.StatusNoContent)
}
