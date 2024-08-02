package entity

import (
	"time"

	"github.com/google/uuid"
)

type DeletionRequest struct {
	Id         string    `json:"id"`
	CustomerId string    `json:"customer_id"`
	Name       string    `json:"name"`
	Address    string    `json:"address"`
	Phone      string    `json:"phone"`
	Executed   bool      `json:"executed"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

func NewDeleteRequest(customerId, name, address, phone string) DeletionRequest {
	return DeletionRequest{
		Id: uuid.NewString(),

		CustomerId: customerId,
		Name:       name,
		Address:    address,
		Phone:      phone,

		Executed:  false,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
}
