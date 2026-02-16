package model

import (
	"natan/fingo/utils"
	"time"
)

type Transaction struct {
	ID        int64         `json:"id"`
	Desc      *string     `json:"description"`
	Amount    utils.Money `json:"amount"`
	IsDebt    bool        `json:"is_debt"`
	CreatedAt time.Time   `json:"created_at"`
	UserID    int64         `json:"user_id"`
}
