package model

import (
	"natan/fingo/utils"
)

type Transaction struct {
	ID        int64       `json:"id"`
	Desc      string     `json:"description,omitempty"`
	Amount    utils.Money `json:"amount"`
	IsDebt    bool        `json:"is_debt"`
	CreatedAt string  `json:"created_at,omitempty"`
	UserID    int64       `json:"user_id"`
}
