package model

import (
	"natan/fingo/utils"
)

// Goal represents a financial goal of the user
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

// GoalUpdate is used for partial updates of Goal, where all fields are optional
type GoalUpdate struct {
	Name     *string      `json:"name,omitempty"`
	Desc     *string      `json:"description,omitempty"`
	Price    *utils.Money `json:"price,omitempty"`
	Pros     *string      `json:"pros,omitempty"`
	Cons     *string      `json:"cons,omitempty"`
	Deadline *string      `json:"deadline,omitempty"`
}
