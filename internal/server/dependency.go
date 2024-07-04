package server

import (
	"github.com/jfelipearaujo-org/ms-customer-management/internal/provider/time_provider"

	customer_repository "github.com/jfelipearaujo-org/ms-customer-management/internal/repository/customer"
	customer_delete_account_svc "github.com/jfelipearaujo-org/ms-customer-management/internal/service/customer/delete_account"
)

type Dependency struct {
	TimeProvider *time_provider.TimeProvider

	CustomerRepository customer_repository.Repository
	CustomerService    customer_delete_account_svc.Service
}
