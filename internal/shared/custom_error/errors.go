package custom_error

import "net/http"

var (
	ErrRequestNotValid BusinessError = New(http.StatusUnprocessableEntity, "validation error", "request not valid, please check the fields")
)
