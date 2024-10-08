package main

import (
	"log"
	"net/http"

	"api/urgency/internal/handler"
)

func main() {
	http.HandleFunc("/hello", handler.HelloHandler)
	log.Println("Service1 is running on port 8081")
	log.Fatal(http.ListenAndServe(":8081", nil))
}
