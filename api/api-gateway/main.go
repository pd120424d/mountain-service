package main

import (
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
)

func main() {
	http.HandleFunc("/activity/", func(w http.ResponseWriter, r *http.Request) {
		proxy(w, r, "http://activity:8081")
	})
	http.HandleFunc("/employee/", func(w http.ResponseWriter, r *http.Request) {
		proxy(w, r, "http://employee:8082")
	})
	http.HandleFunc("/urgency/", func(w http.ResponseWriter, r *http.Request) {
		proxy(w, r, "http://urgency:8083")
	})

	log.Println("API Gateway is running on port 8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func proxy(w http.ResponseWriter, r *http.Request, target string) {
	url, _ := url.Parse(target)
	proxy := httputil.NewSingleHostReverseProxy(url)
	proxy.ServeHTTP(w, r)
}
