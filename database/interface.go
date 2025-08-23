package database

import "time"

// Config holds database configuration
type Config struct {
	Host     string
	Port     int
	Username string
	Password string
	Database string
}

// ServiceRecord represents a service record in the database
type ServiceRecord struct {
	ID            int64     `db:"id"`
	Name          string    `db:"name"`
	Endpoint      string    `db:"endpoint"`
	Type          string    `db:"type"`
	Status        string    `db:"status"`
	LastHeartbeat time.Time `db:"last_heartbeat"`
	CreatedAt     time.Time `db:"created_at"`
	UpdatedAt     time.Time `db:"updated_at"`
}

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