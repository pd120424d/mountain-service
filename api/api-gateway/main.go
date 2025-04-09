package main

import (
	"net/http"
	"net/http/httputil"
	"net/url"

	"github.com/pd120424d/mountain-service/api/shared/utils"
)

func main() {
	log, err := utils.NewLogger("api-gateway")
	if err != nil {
		log.Fatalf("failed to create logger: %v", err)
	}
	defer func(log utils.Logger) {
		err := log.Sync()
		if err != nil {
			log.Fatalf("failed to sync logger: %v", err)
		}
	}(log)

	http.HandleFunc("/activity/", func(w http.ResponseWriter, r *http.Request) {
		proxy(w, r, "http://activity:8081")
	})
	http.HandleFunc("/employee/", func(w http.ResponseWriter, r *http.Request) {
		proxy(w, r, "http://employee:8082")
	})
	http.HandleFunc("/urgency/", func(w http.ResponseWriter, r *http.Request) {
		proxy(w, r, "http://urgency:8083")
	})

	log.Info("API Gateway is running on port 8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatalf("failed to start API Gateway: %v", err)
	}
}

func proxy(w http.ResponseWriter, r *http.Request, target string) {
	url, _ := url.Parse(target)
	proxy := httputil.NewSingleHostReverseProxy(url)
	proxy.ServeHTTP(w, r)
}
