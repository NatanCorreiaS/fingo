package main

import (
	"log"
	"natan/fingo/dbsqlite"
	"net/http"
)

func main() {
	if err := dbsqlite.CheckAndCreate(); err != nil{
		log.Println(err)
	}
	log.Printf("Server listening on PORT 8000...")
	mux := RouterMux()
	err := http.ListenAndServe(":8000", mux)
	if err != nil {
		log.Fatal(err)
	}
}
