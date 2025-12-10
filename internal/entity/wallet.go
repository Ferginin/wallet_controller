package entity

import "github.com/google/uuid"

type Wallet struct {
	ID      uuid.UUID `json:"id"`
	Balance int       `json:"balance"`
}
