package model

import "natan/fingo/utils"

type User struct {
	ID             int64       `json:"id"`
	UserName       string      `json:"user_name"`
	CurrentAmount  utils.Money `json:"current_amount"`
	MonthlyInputs  utils.Money `json:"monthly_inputs"`
	MonthlyOutputs utils.Money `json:"monthly_outputs"`
}
