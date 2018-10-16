package main

import (
	"log"
	"net/http"
)

func main() {
	// Simple static webserver:
	log.Printf("Visit http://localhost:1314 for doc.")
	err := http.ListenAndServe("localhost:1314", http.FileServer(http.Dir(".")))
	if err != nil{
		log.Fatal(err)
	}
}
