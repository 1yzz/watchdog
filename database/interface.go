package database

import "watchdog/ent"

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
	UpdateServiceStatus(serviceID int64, newStatus string) error
	DeleteService(serviceID int64) error

	// Health logging
	LogHealthCheck(status string, serviceCount int) error
}

// Ensure EntClient implementation satisfies the interface
var _ ServiceDB = (*EntClient)(nil)
