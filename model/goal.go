package model

import (
	"natan/fingo/utils"
	"time"
)

type Goal struct {
	ID       int         `json:"id"`
	Name     string      `json:"name"`
	Desc     string      `json:"description,omitempty"`
	Price    utils.Money `json:"price"`
	Pros     *string     `json:"pros,omitempty"`
	Cons     *string     `json:"cons,omitempty"`
	UserID   int         `json:"user_id"`
	Deadline time.Time   `json:"deadline"`
}
