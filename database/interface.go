package database

import (
	"watchdog/ent"
	"watchdog/ent/service"
)

// Config holds database configuration
type Config struct {
	Host     string
	Port     int
	Username string
	Password string
	Database string
}

// ServiceRecord represents a service record in the database
type ServiceRecord = ent.Service

// ServiceDB defines the interface for database operations on services
// This interface is implemented by EntClient
type ServiceDB interface {
	// Connection management
	Close() error
	HealthCheck() error

	// Service operations
	CreateService(service ServiceRecord) (int64, error)
	GetService(serviceID int64) (*ServiceRecord, error)
	ListServices() ([]ServiceRecord, error)
	UpdateService(serviceID int64, newStatus string, name string, serviceType service.Type, endpoint string) error
	DeleteService(serviceID int64) error

	// Health logging
	LogHealthCheck(status string, serviceCount int) error
}

// Ensure EntClient implementation satisfies the ServiceDB interface
var _ ServiceDB = (*EntClient)(nil)
