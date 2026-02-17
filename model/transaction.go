package model

import (
	"natan/fingo/utils"
)

// Transaction represents a financial transaction of the user
type Transaction struct {
	ID        int64       `json:"id"`
	Desc      string      `json:"description,omitempty"`
	Amount    utils.Money `json:"amount"`
	IsDebt    bool        `json:"is_debt"`
	CreatedAt string      `json:"created_at,omitempty"`
	UserID    int64       `json:"user_id"`
}

// TransactionUpdate is used for partial updates of Transaction, where all fields are optional
type TransactionUpdate struct {
	Desc   *string      `json:"description,omitempty"`
	Amount *utils.Money `json:"amount,omitempty"`
	IsDebt *bool        `json:"is_debt,omitempty"`
}
