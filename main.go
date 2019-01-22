package main

import (
	"log"
	"net/http"
)

func main() {

	router := NewRouter()
	log.Printf("[INFO] Server listening on PORT 6767")
	server := http.ListenAndServe(":6767", router)
	log.Fatal(server)
}
