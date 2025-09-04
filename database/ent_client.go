package database

import (
	"context"
	"fmt"
	"log"
	"time"

	"entgo.io/ent/dialect"
	_ "github.com/go-sql-driver/mysql"

	"watchdog/ent"
	"watchdog/ent/service"
)

// EntClient wraps the generated Ent client to provide the same interface as the original DB
type EntClient struct {
	client *ent.Client
}

// Helper function to convert ent.Service to ServiceRecord
func entToServiceRecord(entService *ent.Service) *ServiceRecord {
	serviceRecord := ServiceRecord(*entService)
	return &serviceRecord
}

// Helper function to convert ServiceRecord to ent.Service
func serviceRecordToEnt(serviceRecord ServiceRecord) *ent.Service {
	entService := ent.Service(serviceRecord)
	return &entService
}

// NewEntClient creates a new database connection using the generated Ent client
func NewEntClient(config Config) (*EntClient, error) {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?parseTime=true&loc=UTC",
		config.Username,
		config.Password,
		config.Host,
		config.Port,
		config.Database,
	)

	// Create Ent client with MySQL driver
	client, err := ent.Open(dialect.MySQL, dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to open database connection: %w", err)
	}

	log.Printf("Connected to MySQL database using Ent at %s:%d", config.Host, config.Port)

	return &EntClient{client: client}, nil
}

// AutoMigrate runs automatic schema migration
func (db *EntClient) AutoMigrate(ctx context.Context) error {
	if err := db.client.Schema.Create(ctx); err != nil {
		return fmt.Errorf("failed creating schema resources: %w", err)
	}
	log.Println("Database schema migration completed successfully")
	return nil
}

// Close closes the database connection
func (db *EntClient) Close() error {
	return db.client.Close()
}

// HealthCheck checks if the database connection is healthy
func (db *EntClient) HealthCheck() error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Try a simple query to test the connection
	_, err := db.client.Service.Query().Count(ctx)
	return err
}

// CreateService creates a new service using the generated Ent client
func (db *EntClient) CreateService(serviceRecord ServiceRecord) (int64, error) {
	ctx := context.Background()

	// Convert ServiceRecord to ent.Service for creation
	entService := serviceRecordToEnt(serviceRecord)

	created, err := db.client.Service.Create().
		SetName(entService.Name).
		SetEndpoint(entService.Endpoint).
		SetType(entService.Type).
		SetStatus(entService.Status).
		SetLastHeartbeat(entService.LastHeartbeat).
		Save(ctx)

	if err != nil {
		return 0, fmt.Errorf("failed to create service: %w", err)
	}

	log.Printf("Service %s created with ID %d", created.Name, created.ID)
	return created.ID, nil
}

// GetService retrieves a service by ID using the generated Ent client
func (db *EntClient) GetService(serviceID int64) (*ServiceRecord, error) {
	ctx := context.Background()

	entService, err := db.client.Service.Get(ctx, serviceID)
	if err != nil {
		if ent.IsNotFound(err) {
			return nil, fmt.Errorf("service not found")
		}
		return nil, fmt.Errorf("failed to get service: %w", err)
	}

	return entToServiceRecord(entService), nil
}

// ListServices lists all services using the generated Ent client
func (db *EntClient) ListServices() ([]ServiceRecord, error) {
	ctx := context.Background()

	entServices, err := db.client.Service.Query().
		Order(ent.Desc(service.FieldCreatedAt)).
		All(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to list services: %w", err)
	}

	// Convert ent.Service to ServiceRecord (which is an alias)
	services := make([]ServiceRecord, len(entServices))
	for i, entService := range entServices {
		services[i] = *entToServiceRecord(entService)
	}

	return services, nil
}

// UpdateService updates a service using the generated Ent client
func (db *EntClient) UpdateService(serviceID int64, newStatus string, name string, serviceType service.Type, endpoint string) error {
	ctx := context.Background()

	// Get current service first
	currentService, err := db.GetService(serviceID)
	if err != nil {
		return fmt.Errorf("service not found: %w", err)
	}

	// Use current values for empty parameters to avoid validation errors
	updateName := name
	updateEndpoint := endpoint
	updateServiceType := serviceType

	if updateName == "" {
		updateName = currentService.Name
	}
	if updateEndpoint == "" {
		updateEndpoint = currentService.Endpoint
	}
	// Use zero value check for enum type
	if updateServiceType == "" {
		updateServiceType = currentService.Type
	}

	// Direct type usage - no string conversion needed
	_, err = db.client.Service.UpdateOneID(serviceID).
		SetStatus(newStatus).
		SetName(updateName).
		SetEndpoint(updateEndpoint).
		SetType(updateServiceType).
		SetLastHeartbeat(time.Now()).
		Save(ctx)

	if err != nil {
		return fmt.Errorf("failed to update service: %w", err)
	}

	// Log the status change
	log.Printf("Service %d updated: status=%s, name=%s, type=%s, endpoint=%s",
		serviceID, newStatus, updateName, string(updateServiceType), updateEndpoint)

	return nil
}

// DeleteService deletes a service using the generated Ent client
func (db *EntClient) DeleteService(serviceID int64) error {
	ctx := context.Background()

	err := db.client.Service.DeleteOneID(serviceID).Exec(ctx)
	if err != nil {
		if ent.IsNotFound(err) {
			return fmt.Errorf("service not found")
		}
		return fmt.Errorf("failed to delete service: %w", err)
	}

	return nil
}

// LogHealthCheck logs a health check
func (db *EntClient) LogHealthCheck(status string, serviceCount int) error {
	log.Printf("Health check: %s, services: %d", status, serviceCount)
	return nil
}
