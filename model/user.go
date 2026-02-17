package model

import "natan/fingo/utils"

// User represents a system user with financial information
type User struct {
	ID             int64       `json:"id"`
	UserName       string      `json:"user_name"`
	CurrentAmount  utils.Money `json:"current_amount"`
	MonthlyInputs  utils.Money `json:"monthly_inputs"`
	MonthlyOutputs utils.Money `json:"monthly_outputs"`
}

// UserUpdate is used for partial updates of User, where all fields are optional
type UserUpdate struct {
	UserName       *string      `json:"user_name,omitempty"`
	CurrentAmount  *utils.Money `json:"current_amount,omitempty"`
	MonthlyInputs  *utils.Money `json:"monthly_inputs,omitempty"`
	MonthlyOutputs *utils.Money `json:"monthly_outputs,omitempty"`
}
