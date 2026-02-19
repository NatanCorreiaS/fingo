package main

import (
	"log"
	"net/http"
)

func main() {
	log.Printf("Server listening on PORT 8000...")
	mux := RouterMux()
	err := http.ListenAndServe(":8000", mux)
	if err != nil {
		log.Fatal(err)
	}
}
