package main

import (
	"natan/fingo/controller"
	"net/http"
)

// Route groups an HTTP method, path pattern, and handler function.
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

// registerRoutes registers a slice of routes on the given ServeMux.
func registerRoutes(mux *http.ServeMux, routes []Route) {
	for _, route := range routes {
		mux.HandleFunc(route.Method+" "+route.Path, route.Handler)
	}
}

// RouterMux builds and returns the application ServeMux with all routes registered.
func RouterMux() *http.ServeMux {
	mux := http.NewServeMux()
	registerRoutes(mux, UserRoutes)
	registerRoutes(mux, TransactionRoutes)
	registerRoutes(mux, GoalRoutes)
	return mux
}
