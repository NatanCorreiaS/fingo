package model

import (
	"natan/fingo/utils"
)

type Goal struct {
	ID        int64       `json:"id"`
	Name      string      `json:"name"`
	Desc      string      `json:"description,omitempty"`
	Price     utils.Money `json:"price"`
	Pros      string      `json:"pros,omitempty"`
	Cons      string      `json:"cons,omitempty"`
	UserID    int64       `json:"user_id"`
	CreatedAt string      `json:"created_at,omitempty"`
	Deadline  string      `json:"deadline"`
}
