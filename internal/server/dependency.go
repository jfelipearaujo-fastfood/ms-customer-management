package server

import "github.com/jfelipearaujo-org/ms-customer-management/internal/provider/time_provider"

type Dependency struct {
	TimeProvider *time_provider.TimeProvider
}
