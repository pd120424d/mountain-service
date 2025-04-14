package main

import (
	"encoding/json"
	"net/http"
	"time"
)

var (
	Version   = "dev"
	GitSHA    = "unknown"
	startTime = time.Now()
)

func main() {
	http.HandleFunc("api/v1/version", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		uptime := time.Since(startTime).String()

		json.NewEncoder(w).Encode(map[string]string{
			"version": Version,
			"gitSha":  GitSHA,
			"uptime":  uptime,
		})
	})

	if err := http.ListenAndServe(":8090", nil); err != nil {
		panic("failed to start version-service: " + err.Error())
	}
}
