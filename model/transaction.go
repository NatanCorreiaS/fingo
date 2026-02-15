package model

import (
	"natan/fingo/utils"
	"time"
)

type Transaction struct {
	ID        int         `json:"id"`
	Desc      *string     `json:"description"`
	Amount    utils.Money `json:"amount"`
	IsDebt    bool        `json:"is_debt"`
	CreatedAt time.Time   `json:"created_at"`
	UserID    int         `json:"user_id"`
}
