package main

import (
	"natan/fingo/controller"
	"net/http"
)

type Route struct {
	Method  string
	Path    string
	Handler http.HandlerFunc
}

var UserRoutes = []Route{
	{"GET", "/users/{id}", controller.GetUserByIDHandler},
	{"GET", "/users", controller.GetAllUsersHandler},
	{"POST", "/users", controller.CreateUserHandler},
	{"PATCH", "/users/{id}", controller.UpdateUserByIDHandler},
	{"DELETE", "/users/{id}", controller.DeleteUserByIDHandler},
}

var TransactionRoutes = []Route{
	{"GET", "/transactions/{id}", controller.GetTransactionByIDHandler},
	{"GET", "/transactions", controller.GetAllTransactionsHandler},
	{"POST", "/transactions", controller.CreateTransactionHandler},
	{"PATCH", "/transactions/{id}", controller.UpdateTransactionByIDHandler},
	{"DELETE", "/transactions/{id}", controller.DeleteTransactionByIDHandler},
}

var GoalRoutes = []Route{
	{"GET", "/goals/{id}", controller.GetGoalByIDHandler},
	{"GET", "/goals", controller.GetAllGoalsHandler},
	{"POST", "/goals", controller.CreateGoalHandler},
	{"PATCH", "/goals/{id}", controller.UpdateGoalByIDHandler},
	{"DELETE", "/goals/{id}", controller.DeleteGoalByIDHandler},
}

func registerRoutes(mux *http.ServeMux, routes []Route) {
	for _, route := range routes {
		mux.HandleFunc(route.Method+" "+route.Path, route.Handler)
	}
}

func RouterMux() *http.ServeMux {
	mux := http.NewServeMux()
	registerRoutes(mux, UserRoutes)
	registerRoutes(mux, TransactionRoutes)
	registerRoutes(mux, GoalRoutes)
	return mux
}
