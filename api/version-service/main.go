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
	Version    = "dev"
	GitSHA     = "unknown"
	FullGitSHA = "unknown"
	GitTag     = ""
	startTime  = time.Now()
)

func main() {
	log, err := utils.NewLogger("version-service")
	if err != nil {
		panic("failed to create new logger: " + err.Error())
	}

	log.Info("Starting Version service on :8090...") // Updated to trigger deployment

	r := gin.Default()

	r.Use(log.RequestLogger())

	r.GET("/api/v1/version", versionHandler)

	r.GET("/api/v1/health", func(c *gin.Context) {
		log.WithContext(c.Request.Context()).Info("Health endpoint hit")
		c.JSON(200, gin.H{"message": "Service is healthy", "service": "version"})
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
	log, _ := utils.NewLogger("version-service")
	log = log.WithContext(c.Request.Context())

	// Prefer provided short SHA; if it looks long/unknown, derive from full SHA if available
	shortSHA := GitSHA
	if (shortSHA == "" || shortSHA == "unknown") && len(FullGitSHA) >= 8 {
		shortSHA = FullGitSHA[:8]
	}
	if len(shortSHA) > 8 {
		shortSHA = shortSHA[:8]
	}

	response := gin.H{
		"version":    Version,
		"gitSHA":     shortSHA,
		"gitFullSHA": FullGitSHA,
		"gitTag":     GitTag,
		"uptime":     time.Since(startTime).String(),
	}
	c.JSON(http.StatusOK, response)
}
