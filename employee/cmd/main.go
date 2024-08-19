package main

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"mountain-service/employee/config"
	"mountain-service/employee/internal/handlers"
	"mountain-service/employee/internal/models"
	"mountain-service/employee/internal/repositories"
	"mountain-service/shared/utils"

	"github.com/gin-gonic/gin"
)

const hostname = "localhost"

func main() {
	env := os.Getenv("APP_ENV")
	if env == "" {
		env = "staging"
	}

	log := utils.NewLogger()
	defer func(log utils.Logger) {
		err := log.Sync()
		if err != nil {
			log.Fatalf("failed to sync logger: %v", err)
		}
	}(log)

	// Create the employee_service database if it doesn't exist
	config.CreateEmployeeDB(log, hostname)

	db := config.GetEmployeeDB(log, hostname)

	// Auto migrate the models
	err := db.AutoMigrate(&models.Employee{})
	if err != nil {
		log.Fatalf("failed to migrate employee models: %v", err)
	}

	employeeRepo := repositories.NewEmployeeRepository(db)
	employeeHandler := handlers.NewEmployeeHandler(log, employeeRepo)

	r := gin.Default()
	r.POST("/employees", employeeHandler.CreateEmployee)
	r.GET("/employees", employeeHandler.GetAllEmployees)
	r.DELETE("/employees/:id", employeeHandler.DeleteEmployee)

	srv := &http.Server{
		Addr:    fmt.Sprintf(":%v", config.ServerPort),
		Handler: r,
	}

	// Run server in a goroutine
	go func() {
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("listen: %s\n", err)
		}
	}()

	log.Infof("Starting Employee Service on port %s", config.ServerPort)

	// Wait for interrupt signal to gracefully shut down the server
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Info("Shutting down Employee Service...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("Employee Service forced to shutdown: %v", err)
	}

	log.Info("Employee Service exiting")
}

func executeSchema(db *sql.DB) error {
	schema, err := ioutil.ReadFile("employee/schema.sql")
	if err != nil {
		return err
	}
	_, err = db.Exec(string(schema))
	return err
}
