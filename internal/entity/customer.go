package entity

import "time"

type Customer struct {
	Id          string    `json:"id"`
	DocumentId  string    `json:"document_id"`
	Password    string    `json:"password"`
	IsAnonymous bool      `json:"is_anonymous"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}
