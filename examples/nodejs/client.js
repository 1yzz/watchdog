const grpc = require('@grpc/grpc-js');
const protoLoader = require('@grpc/proto-loader');
const path = require('path');

// Load the protobuf definition
const PROTO_PATH = path.join(__dirname, '../../proto/watchdog.proto');

const packageDefinition = protoLoader.loadSync(PROTO_PATH, {
  keepCase: true,
  longs: String,
  enums: String,
  defaults: true,
  oneofs: true
});

const watchdog = grpc.loadPackageDefinition(packageDefinition).watchdog;

// Create gRPC client
const client = new watchdog.WatchdogService('localhost:50051', grpc.credentials.createInsecure());

class WatchdogClient {
  constructor() {
    this.client = client;
  }

  // Health check
  async getHealth() {
    return new Promise((resolve, reject) => {
      this.client.GetHealth({}, (error, response) => {
        if (error) {
          reject(error);
        } else {
          resolve(response);
        }
      });
    });
  }

  // Register a service
  async registerService(name, endpoint, type = 'SERVICE_TYPE_UNSPECIFIED') {
    return new Promise((resolve, reject) => {
      this.client.RegisterService({ name, endpoint, type }, (error, response) => {
        if (error) {
          reject(error);
        } else {
          resolve(response);
        }
      });
    });
  }

  // Unregister a service
  async unregisterService(serviceId) {
    return new Promise((resolve, reject) => {
      this.client.UnregisterService({ service_id: serviceId }, (error, response) => {
        if (error) {
          reject(error);
        } else {
          resolve(response);
        }
      });
    });
  }

  // List all services
  async listServices() {
    return new Promise((resolve, reject) => {
      this.client.ListServices({}, (error, response) => {
        if (error) {
          reject(error);
        } else {
          resolve(response);
        }
      });
    });
  }

  // Update service status
  async updateServiceStatus(serviceId, status) {
    return new Promise((resolve, reject) => {
      this.client.UpdateServiceStatus({ service_id: serviceId, status }, (error, response) => {
        if (error) {
          reject(error);
        } else {
          resolve(response);
        }
      });
    });
  }
}

// Example usage
async function main() {
  const watchdogClient = new WatchdogClient();

  try {
    console.log('=== Watchdog gRPC Client Demo ===\n');

    // 1. Health check
    console.log('1. Checking health...');
    const health = await watchdogClient.getHealth();
    console.log('Health:', health);
    console.log();

    // 2. Register a service
    console.log('2. Registering a service...');
    const registration = await watchdogClient.registerService('my-web-service', 'http://localhost:3000', 'SERVICE_TYPE_HTTP');
    console.log('Registration:', registration);
    const serviceId = registration.service_id;
    console.log();

    // 3. List services
    console.log('3. Listing services...');
    const services = await watchdogClient.listServices();
    console.log('Services:', JSON.stringify(services, null, 2));
    console.log();

    // 4. Update service status
    console.log('4. Updating service status...');
    const statusUpdate = await watchdogClient.updateServiceStatus(serviceId, 'healthy');
    console.log('Status update:', statusUpdate);
    console.log();

    // 5. List services again to see the update
    console.log('5. Listing services after update...');
    const updatedServices = await watchdogClient.listServices();
    console.log('Updated services:', JSON.stringify(updatedServices, null, 2));
    console.log();

    // 6. Unregister the service
    console.log('6. Unregistering service...');
    const unregistration = await watchdogClient.unregisterService(serviceId);
    console.log('Unregistration:', unregistration);
    console.log();

    // 7. Final list to confirm removal
    console.log('7. Final service list...');
    const finalServices = await watchdogClient.listServices();
    console.log('Final services:', JSON.stringify(finalServices, null, 2));

  } catch (error) {
    console.error('Error:', error.message);
    if (error.code) {
      console.error('Error code:', error.code);
    }
  }
}

// Run the demo if this file is executed directly
if (require.main === module) {
  main().catch(console.error);
}

module.exports = { WatchdogClient };