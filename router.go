package main

import (
	"natan/fingo/controller"
	"net/http"
)

type Route struct{
	Method string
	Path string
	Handler http.HandlerFunc
}

var UserRoutes = []Route{
	{"GET", "/users/{id}", controller.GetUserByIDHandler},
    {"GET", "/users", controller.GetAllUsers},
    {"POST", "/users", controller.CreateUserHandler},
    {"PATCH", "/users/{id}", controller.UpdateUserByID},
    {"DELETE", "/users/{id}", controller.DeleteUserByID},
}

func registerRoutes(mux *http.ServeMux, routes []Route){
	for _, route := range(routes){
		mux.HandleFunc(route.Method+" "+route.Path, route.Handler)
	}
}

func RouterMux() *http.ServeMux {
	mux := http.NewServeMux()
	registerRoutes(mux, UserRoutes)
	return mux
}
