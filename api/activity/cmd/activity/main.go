package main

import (
	"api/activity/internal/handler"
	"log"
	"net/http"
)

func main() {
	http.HandleFunc("/hello", handler.HelloHandler)
	log.Println("Service1 is running on port 8081")
	log.Fatal(http.ListenAndServe(":8081", nil))
}
