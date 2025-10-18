package config

import "time"

// Package config contains service configuration constants and utilities

const (
	EmployeeServiceName = "employee-service"
	UrgencyServiceName  = "urgency-service"
	ActivityServiceName = "activity-service"

	EmployeeDBName = "employee_service_db"
	UrgencyDBName  = "urgency_service_db"
	ActivityDBName = "activity_service_db"

	EmployeeServicePort = "8082"
	UrgencyServicePort  = "8083"
	ActivityServicePort = "8084"
)

const (
	DB_USER     = "DB_USER"
	DB_PASSWORD = "DB_PASSWORD"
)

const (
	REDIS_ADDR = "REDIS_ADDR"
	REDIS_DB   = "REDIS_DB"
)

const (
	DefaultListTimeout  = 300 * time.Millisecond
	CursorListTimeout   = 1 * time.Second
	PostgresListTimeout = 3 * time.Second
	CountTimeout        = 2 * time.Second

	// Fresh-read (RYW) propagation via header and default window duration
	FreshWindowHeader  = "X-Fresh-Until"
	DefaultFreshWindow = 2 * time.Second
)
