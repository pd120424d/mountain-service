package main

import (
	"log"
	"net/http"

	"github.com/pd120424d/mountain-service/api/urgency/internal/handler"
)

func main() {
	http.HandleFunc("/hello", handler.HelloHandler)
	log.Println("Service1 is running on port 8081")
	log.Fatal(http.ListenAndServe(":8081", nil))
}
