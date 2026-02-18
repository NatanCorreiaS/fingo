package main

import (
	"log"
	"natan/fingo/controller"
	"net/http"
)

func main() {
	mux := http.NewServeMux()

	mux.HandleFunc("GET /users/{id}", controller.GetUserByIDHandler)
	mux.HandleFunc("GET /users", controller.GetAllUsers)
	mux.HandleFunc("POST /users", controller.CreateUserHandler)
	mux.HandleFunc("PATCH /users/{id}", controller.UpdateUserByID)
	mux.HandleFunc("DELETE /users/{id}", controller.DeleteUserByID)

	log.Printf("Server listening on PORT 8000...")
	err := http.ListenAndServe(":8000", mux)
	if err != nil {
		log.Fatal(err)
	}
}
