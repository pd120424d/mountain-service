package main

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/pd120424d/mountain-service/api/shared/utils"

	"github.com/gin-gonic/gin"
)

var (
	Version   = "dev"
	GitSHA    = "unknown"
	startTime = time.Now()
)

func main() {
	log, err := utils.NewLogger("version-service")
	if err != nil {
		panic("failed to create new logger: " + err.Error())
	}

	log.Info("Starting Version service on :8090...")

	r := gin.Default()

	r.Use(log.RequestLogger())

	r.GET("/api/v1/version", versionHandler)

	r.GET("/ping", func(c *gin.Context) {
		log.Info("Ping route hit")
		c.JSON(http.StatusOK, gin.H{"message": "pong"})
	})

	srv := &http.Server{
		Addr:    fmt.Sprintf(":%v", "8090"),
		Handler: r.Handler(),
	}

	// Run server in a goroutine
	go func() {
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("failed to start Version service: %v", err)
		}
	}()

	// Wait for interrupt signal to gracefully shut down the server
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Info("Shutting down Version service...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("Version service forced to shutdown: %v", err)
	}

	log.Info("Version service exiting")

}

func versionHandler(c *gin.Context) {
	response := map[string]string{
		"version": Version,
		"gitSHA":  GitSHA,
		"uptime":  time.Since(startTime).String(),
	}
	c.JSON(http.StatusOK, response)
}
