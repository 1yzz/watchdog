import * as grpc from '@grpc/grpc-js';
import { WatchdogServiceClient } from './generated/watchdog_grpc_pb';
import * as pb from './generated/watchdog_pb';

export interface WatchdogClientOptions {
  host: string;
  port: number;
  credentials?: grpc.ChannelCredentials;
  timeout?: number;
}

export interface ServiceRegistration {
  name: string;
  endpoint: string;
  type: pb.ServiceType;
}

export interface ServiceUpdate {
  serviceId: string;
  status: string;
  name?: string;
  endpoint?: string;
  type?: pb.ServiceType;
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

  /**
   * Generic method to call any gRPC service method dynamically
   */
  private callMethod<TRequest>(
    methodName: string,
    requestClass: new () => TRequest,
    requestData: Record<string, any>,
    responseFields: string[]
  ): Promise<Record<string, any>> {
    return new Promise((resolve, reject) => {
      const request = new requestClass();
      const metadata = new grpc.Metadata();

      // Dynamically set request fields
      Object.entries(requestData).forEach(([key, value]) => {
        if (value !== undefined && value !== null) {
          const setterName = `set${key.charAt(0).toUpperCase() + key.slice(1)}`;
          if (typeof (request as any)[setterName] === 'function') {
            (request as any)[setterName](value);
          }
        }
      });

      // Dynamically call the client method
      const clientMethod = (this.client as any)[methodName];
      if (typeof clientMethod !== 'function') {
        reject(new Error(`Method ${methodName} not found on client`));
        return;
      }

      clientMethod.call(
        this.client,
        request,
        metadata,
        { deadline: this.createDeadline() },
        (error: any, response: any) => {
          if (error) {
            reject(new Error(`${methodName} failed: ${error.message}`));
            return;
          }

          if (!response) {
            reject(new Error('No response received'));
            return;
          }

          // Dynamically extract response fields
          const result: Record<string, any> = {};
          responseFields.forEach(field => {
            const getterName = `get${field.charAt(0).toUpperCase() + field.slice(1)}`;
            if (typeof response[getterName] === 'function') {
              result[field] = response[getterName]();
            }
          });

          resolve(result);
        }
      );
    });
  }

  getHealth(): Promise<{ status: string; message: string }> {
    return this.callMethod(
      'getHealth',
      pb.HealthRequest,
      {},
      ['status', 'message']
    ) as Promise<{ status: string; message: string }>;
  }

  registerService(service: ServiceRegistration): Promise<{ serviceId: string; message: string }> {
    return this.callMethod(
      'registerService',
      pb.RegisterServiceRequest,
      {
        name: service.name,
        endpoint: service.endpoint,
        type: service.type,
      },
      ['serviceId', 'message']
    ) as Promise<{ serviceId: string; message: string }>;
  }

  unregisterService(serviceId: string): Promise<{ message: string }> {
    return this.callMethod(
      'unregisterService',
      pb.UnregisterServiceRequest,
      { serviceId },
      ['message']
    ) as Promise<{ message: string }>;
  }

  async listServices(): Promise<pb.ServiceInfo[]> {
    const result = await this.callMethod(
      'listServices',
      pb.ListServicesRequest,
      {},
      ['servicesList']
    );
    return result.servicesList as pb.ServiceInfo[];
  }

  updateService(update: ServiceUpdate): Promise<{ message: string }> {
    const requestData: Record<string, any> = {
      serviceId: update.serviceId,
      status: update.status,
    };

    // Add optional fields if provided
    if (update.name) requestData.name = update.name;
    if (update.endpoint) requestData.endpoint = update.endpoint;
    if (update.type !== undefined) requestData.type = update.type;

    return this.callMethod(
      'updateService',
      pb.UpdateServiceRequest,
      requestData,
      ['message']
    ) as Promise<{ message: string }>;
  }

  // Backward compatibility method for status-only updates
  updateServiceStatus(update: { serviceId: string; status: string }): Promise<{ message: string }> {
    return this.updateService(update);
  }

  close(): void {
    this.client.close();
  }

  private createDeadline(): grpc.Deadline {
    return Date.now() + this.timeout;
  }
}

// Re-export protobuf types for convenience
export { pb as ProtobufTypes };
export type ServiceType = pb.ServiceType;
export type ServiceInfo = pb.ServiceInfo;