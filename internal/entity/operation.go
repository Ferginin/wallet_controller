package entity

import "github.com/google/uuid"

type OperationRequest struct {
	WalletID      uuid.UUID `json:"wallet_id"`
	OperationType string    `json:"operation_type"`
	Amount        int       `json:"amount"`
}
