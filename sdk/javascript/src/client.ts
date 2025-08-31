import * as grpc from '@grpc/grpc-js';
import { WatchdogServiceClient } from './generated/watchdog_grpc_pb';
import {
  HealthRequest,
  RegisterServiceRequest,
  UnregisterServiceRequest,
  ListServicesRequest,
  UpdateServiceStatusRequest,
  ServiceType,
  ServiceInfo,
} from './generated/watchdog_pb';

export interface WatchdogClientOptions {
  host: string;
  port: number;
  credentials?: grpc.ChannelCredentials;
  timeout?: number;
}

export interface ServiceRegistration {
  name: string;
  endpoint: string;
  type: ServiceType;
}

export interface ServiceUpdate {
  serviceId: string;
  status: string;
}

export class WatchdogClient {
  private client: WatchdogServiceClient;
  private timeout: number;

  constructor(options: WatchdogClientOptions) {
    const { host, port, credentials = grpc.credentials.createInsecure(), timeout = 5000 } = options;
    const target = `${host}:${port}`;
    
    this.client = new WatchdogServiceClient(target, credentials);
    this.timeout = timeout;
  }

  getHealth(): Promise<{ status: string; message: string }> {
    return new Promise((resolve, reject) => {
      const request = new HealthRequest();
      const metadata = new grpc.Metadata();
      
      this.client.getHealth(request, metadata, { deadline: this.createDeadline() }, (error, response) => {
        if (error) {
          reject(new Error(`Health check failed: ${error.message}`));
          return;
        }
        
        if (!response) {
          reject(new Error('No response received'));
          return;
        }

        resolve({
          status: response.getStatus(),
          message: response.getMessage(),
        });
      });
    });
  }

  registerService(service: ServiceRegistration): Promise<{ serviceId: string; message: string }> {
    return new Promise((resolve, reject) => {
      const request = new RegisterServiceRequest();
      request.setName(service.name);
      request.setEndpoint(service.endpoint);
      request.setType(service.type);
      const metadata = new grpc.Metadata();
      
      this.client.registerService(request, metadata, { deadline: this.createDeadline() }, (error, response) => {
        if (error) {
          reject(new Error(`Service registration failed: ${error.message}`));
          return;
        }
        
        if (!response) {
          reject(new Error('No response received'));
          return;
        }

        resolve({
          serviceId: response.getServiceId(),
          message: response.getMessage(),
        });
      });
    });
  }

  unregisterService(serviceId: string): Promise<{ message: string }> {
    return new Promise((resolve, reject) => {
      const request = new UnregisterServiceRequest();
      request.setServiceId(serviceId);
      const metadata = new grpc.Metadata();
      
      this.client.unregisterService(request, metadata, { deadline: this.createDeadline() }, (error, response) => {
        if (error) {
          reject(new Error(`Service unregistration failed: ${error.message}`));
          return;
        }
        
        if (!response) {
          reject(new Error('No response received'));
          return;
        }

        resolve({
          message: response.getMessage(),
        });
      });
    });
  }

  listServices(): Promise<ServiceInfo[]> {
    return new Promise((resolve, reject) => {
      const request = new ListServicesRequest();
      const metadata = new grpc.Metadata();
      
      this.client.listServices(request, metadata, { deadline: this.createDeadline() }, (error, response) => {
        if (error) {
          reject(new Error(`List services failed: ${error.message}`));
          return;
        }
        
        if (!response) {
          reject(new Error('No response received'));
          return;
        }

        resolve(response.getServicesList());
      });
    });
  }

  updateServiceStatus(update: ServiceUpdate): Promise<{ message: string }> {
    return new Promise((resolve, reject) => {
      const request = new UpdateServiceStatusRequest();
      request.setServiceId(update.serviceId);
      request.setStatus(update.status);
      const metadata = new grpc.Metadata();
      
      this.client.updateServiceStatus(request, metadata, { deadline: this.createDeadline() }, (error, response) => {
        if (error) {
          reject(new Error(`Status update failed: ${error.message}`));
          return;
        }
        
        if (!response) {
          reject(new Error('No response received'));
          return;
        }

        resolve({
          message: response.getMessage(),
        });
      });
    });
  }

  close(): void {
    this.client.close();
  }

  private createDeadline(): grpc.Deadline {
    return Date.now() + this.timeout;
  }
}

export * from './generated/watchdog_pb';
export { ServiceType } from './generated/watchdog_pb';