package server

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os/exec"
	"strconv"
	"strings"
	"time"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"watchdog/api"
	"watchdog/database"
	"watchdog/ent/service"
)

type WatchdogServer struct {
	api.UnimplementedWatchdogServiceServer
	db database.ServiceDB
}

func NewWatchdogServer(db database.ServiceDB) *WatchdogServer {
	return &WatchdogServer{
		db: db,
	}
}

func (s *WatchdogServer) GetHealth(ctx context.Context, req *api.HealthRequest) (*api.HealthResponse, error) {
	if err := s.db.HealthCheck(); err != nil {
		log.Printf("Database health check failed: %v", err)
		return &api.HealthResponse{
			Status:  "unhealthy",
			Message: "Database connection failed",
		}, nil
	}

	services, err := s.db.ListServices()
	if err != nil {
		log.Printf("Failed to count services: %v", err)
	}

	err = s.db.LogHealthCheck("healthy", len(services))
	if err != nil {
		log.Printf("Failed to log health check: %v", err)
	}

	return &api.HealthResponse{
		Status:  "healthy",
		Message: fmt.Sprintf("Watchdog service is running with %d registered services", len(services)),
	}, nil
}

func (s *WatchdogServer) RegisterService(ctx context.Context, req *api.RegisterServiceRequest) (*api.RegisterServiceResponse, error) {
	if req.Name == "" {
		return nil, status.Errorf(codes.InvalidArgument, "service name cannot be empty")
	}

	if req.Endpoint == "" {
		return nil, status.Errorf(codes.InvalidArgument, "service endpoint cannot be empty")
	}

	service := database.ServiceRecord{
		Name:          req.Name,
		Endpoint:      req.Endpoint,
		Type:          service.Type(req.Type),
		Status:        "active",
		LastHeartbeat: time.Now(),
	}

	serviceID, err := s.db.CreateService(service)
	if err != nil {
		log.Printf("Failed to create service: %v", err)
		return nil, status.Errorf(codes.Internal, "failed to register service")
	}

	return &api.RegisterServiceResponse{
		ServiceId: fmt.Sprintf("%d", serviceID),
		Message:   fmt.Sprintf("Service %s registered successfully with ID %d", req.Name, serviceID),
	}, nil
}

func (s *WatchdogServer) UnregisterService(ctx context.Context, req *api.UnregisterServiceRequest) (*api.UnregisterServiceResponse, error) {
	if req.ServiceId == "" {
		return nil, status.Errorf(codes.InvalidArgument, "service ID cannot be empty")
	}

	// Convert string ID to int64
	serviceID, err := strconv.ParseInt(req.ServiceId, 10, 64)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid service ID format")
	}

	err = s.db.DeleteService(serviceID)
	if err != nil {
		if err.Error() == "service not found" {
			return nil, status.Errorf(codes.NotFound, "service not found")
		}
		log.Printf("Failed to delete service: %v", err)
		return nil, status.Errorf(codes.Internal, "failed to unregister service")
	}

	return &api.UnregisterServiceResponse{
		Message: "Service unregistered successfully",
	}, nil
}

func (s *WatchdogServer) ListServices(ctx context.Context, req *api.ListServicesRequest) (*api.ListServicesResponse, error) {
	services, err := s.db.ListServices()
	if err != nil {
		log.Printf("Failed to list services: %v", err)
		return nil, status.Errorf(codes.Internal, "failed to list services")
	}

	var apiServices []*api.ServiceInfo
	for _, service := range services {
		apiService := &api.ServiceInfo{
			Id:            fmt.Sprintf("%d", service.ID),
			Name:          service.Name,
			Endpoint:      service.Endpoint,
			Status:        service.Status,
			LastHeartbeat: service.LastHeartbeat.Unix(),
			Type:          stringToServiceType(string(service.Type)),
		}
		apiServices = append(apiServices, apiService)
	}

	return &api.ListServicesResponse{
		Services: apiServices,
	}, nil
}

func (s *WatchdogServer) UpdateService(ctx context.Context, req *api.UpdateServiceRequest) (*api.UpdateServiceResponse, error) {
	if req.ServiceId == "" {
		return nil, status.Errorf(codes.InvalidArgument, "service ID cannot be empty")
	}

	// Convert string ID to int64
	serviceID, err := strconv.ParseInt(req.ServiceId, 10, 64)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid service ID format")
	}

	err = s.db.UpdateService(serviceID, req.Status, req.Name, apiToEntServiceType(req.Type), req.Endpoint)
	if err != nil {
		if err.Error() == "service not found" {
			return nil, status.Errorf(codes.NotFound, "service not found")
		}
		log.Printf("Failed to update service: %v", err)
		return nil, status.Errorf(codes.Internal, "failed to update service")
	}

	return &api.UpdateServiceResponse{
		Message: "Service updated successfully",
	}, nil
}

func stringToServiceType(serviceType string) api.ServiceType {
	switch serviceType {
	case "SERVICE_TYPE_HTTP":
		return api.ServiceType_SERVICE_TYPE_HTTP
	case "SERVICE_TYPE_GRPC":
		return api.ServiceType_SERVICE_TYPE_GRPC
	case "SERVICE_TYPE_DATABASE":
		return api.ServiceType_SERVICE_TYPE_DATABASE
	case "SERVICE_TYPE_CACHE":
		return api.ServiceType_SERVICE_TYPE_CACHE
	case "SERVICE_TYPE_QUEUE":
		return api.ServiceType_SERVICE_TYPE_QUEUE
	case "SERVICE_TYPE_STORAGE":
		return api.ServiceType_SERVICE_TYPE_STORAGE
	case "SERVICE_TYPE_EXTERNAL_API":
		return api.ServiceType_SERVICE_TYPE_EXTERNAL_API
	case "SERVICE_TYPE_MICROSERVICE":
		return api.ServiceType_SERVICE_TYPE_MICROSERVICE
	case "SERVICE_TYPE_OTHER":
		return api.ServiceType_SERVICE_TYPE_OTHER
	case "SERVICE_TYPE_SYSTEMD":
		return api.ServiceType_SERVICE_TYPE_SYSTEMD
	default:
		return api.ServiceType_SERVICE_TYPE_UNSPECIFIED
	}
}

func serviceTypeToString(serviceType api.ServiceType) string {
	switch serviceType {
	case api.ServiceType_SERVICE_TYPE_HTTP:
		return "SERVICE_TYPE_HTTP"
	case api.ServiceType_SERVICE_TYPE_GRPC:
		return "SERVICE_TYPE_GRPC"
	case api.ServiceType_SERVICE_TYPE_DATABASE:
		return "SERVICE_TYPE_DATABASE"
	case api.ServiceType_SERVICE_TYPE_CACHE:
		return "SERVICE_TYPE_CACHE"
	case api.ServiceType_SERVICE_TYPE_QUEUE:
		return "SERVICE_TYPE_QUEUE"
	case api.ServiceType_SERVICE_TYPE_STORAGE:
		return "SERVICE_TYPE_STORAGE"
	case api.ServiceType_SERVICE_TYPE_EXTERNAL_API:
		return "SERVICE_TYPE_EXTERNAL_API"
	case api.ServiceType_SERVICE_TYPE_MICROSERVICE:
		return "SERVICE_TYPE_MICROSERVICE"
	case api.ServiceType_SERVICE_TYPE_OTHER:
		return "SERVICE_TYPE_OTHER"
	case api.ServiceType_SERVICE_TYPE_SYSTEMD:
		return "SERVICE_TYPE_SYSTEMD"
	default:
		return "SERVICE_TYPE_UNSPECIFIED"
	}
}

func apiToEntServiceType(apiType api.ServiceType) service.Type {
	switch apiType {
	case api.ServiceType_SERVICE_TYPE_HTTP:
		return service.TypeSERVICE_TYPE_HTTP
	case api.ServiceType_SERVICE_TYPE_GRPC:
		return service.TypeSERVICE_TYPE_GRPC
	case api.ServiceType_SERVICE_TYPE_DATABASE:
		return service.TypeSERVICE_TYPE_DATABASE
	case api.ServiceType_SERVICE_TYPE_CACHE:
		return service.TypeSERVICE_TYPE_CACHE
	case api.ServiceType_SERVICE_TYPE_QUEUE:
		return service.TypeSERVICE_TYPE_QUEUE
	case api.ServiceType_SERVICE_TYPE_STORAGE:
		return service.TypeSERVICE_TYPE_STORAGE
	case api.ServiceType_SERVICE_TYPE_EXTERNAL_API:
		return service.TypeSERVICE_TYPE_EXTERNAL_API
	case api.ServiceType_SERVICE_TYPE_MICROSERVICE:
		return service.TypeSERVICE_TYPE_MICROSERVICE
	case api.ServiceType_SERVICE_TYPE_OTHER:
		return service.TypeSERVICE_TYPE_OTHER
	case api.ServiceType_SERVICE_TYPE_SYSTEMD:
		return service.TypeSERVICE_TYPE_SYSTEMD
	default:
		return service.TypeSERVICE_TYPE_UNSPECIFIED
	}
}

func (s *WatchdogServer) checkHTTPHealth(endpoint string) (string, error) {
	client := http.Client{
		Timeout: 10 * time.Second,
	}

	resp, err := client.Get(endpoint)
	if err != nil {
		return "unhealthy", err
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "unhealthy", fmt.Errorf("HTTP status code: %d", resp.StatusCode)
	}

	return "healthy", nil
}

func (s *WatchdogServer) checkSystemdHealth(endpoint string) (string, error) {
	out, err := exec.CommandContext(context.Background(), "systemctl", "is-active", endpoint).Output()
	if err != nil {
		return "unhealthy", err
	}
	if strings.TrimSpace(string(out)) != "active" {
		return "unhealthy", fmt.Errorf("systemd health check command returned: %s", string(out))
	}

	return "healthy", nil
}

func (s *WatchdogServer) checkGRPCHealth(endpoint string) (string, error) {
	return "healthy", nil
}

func (s *WatchdogServer) checkDatabaseHealth(endpoint string) (string, error) {
	return "healthy", nil
}

func (s *WatchdogServer) CheckServiceHealth(ctx context.Context, req *api.CheckServiceHealthRequest) (*api.HealthResponse, error) {
	if req.ServiceId == "" {
		return nil, status.Errorf(codes.InvalidArgument, "service ID cannot be empty")
	}

	serviceID, err := strconv.ParseInt(req.ServiceId, 10, 64)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid service ID format")
	}

	service, err := s.db.GetService(serviceID)
	if err != nil {
		return nil, status.Errorf(codes.NotFound, "service not found")
	}

	healthStatus := "unhealthy"

	switch service.Type {
	case "SERVICE_TYPE_HTTP":
		healthStatus, err = s.checkHTTPHealth(service.Endpoint)
		if err != nil {
			log.Printf("HTTP health check failed for service %d (%s): %v", serviceID, service.Endpoint, err)
			return &api.HealthResponse{
				Status:  "unhealthy",
				Message: fmt.Sprintf("Service is unreachable: %v", err),
			}, nil
		}
	case "SERVICE_TYPE_GRPC":
		healthStatus, err = s.checkGRPCHealth(service.Endpoint)
		if err != nil {
			log.Printf("gRPC health check failed for service %d (%s): %v", serviceID, service.Endpoint, err)
			return &api.HealthResponse{
				Status:  "unhealthy",
				Message: fmt.Sprintf("Service is unreachable: %v", err),
			}, nil
		}
	case "SERVICE_TYPE_SYSTEMD":
		healthStatus, err = s.checkSystemdHealth(service.Endpoint)
		if err != nil {
			log.Printf("Systemd health check failed for service %d (%s): %v", serviceID, service.Endpoint, err)
			return &api.HealthResponse{
				Status:  "unhealthy",
				Message: fmt.Sprintf("Service is unreachable: %v", err),
			}, nil
		}
	case "SERVICE_TYPE_DATABASE", "SERVICE_TYPE_CACHE", "SERVICE_TYPE_QUEUE", "SERVICE_TYPE_STORAGE", "SERVICE_TYPE_EXTERNAL_API", "SERVICE_TYPE_MICROSERVICE", "SERVICE_TYPE_OTHER":
		healthStatus, err = s.checkDatabaseHealth(service.Endpoint)
		if err != nil {
			log.Printf("Database health check failed for service %d (%s): %v", serviceID, service.Endpoint, err)
			return &api.HealthResponse{
				Status:  "unhealthy",
				Message: fmt.Sprintf("Service is unreachable: %v", err),
			}, nil
		}
	default:
		return nil, status.Errorf(codes.InvalidArgument, "Unsupported service type")
	}

	return &api.HealthResponse{
		Status:  healthStatus,
		Message: "Service health checked successfully",
	}, nil
}
