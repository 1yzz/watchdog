const grpc = require('@grpc/grpc-js');
const protoLoader = require('@grpc/proto-loader');
const path = require('path');

// Alternative implementation using async/await with promisify
const { promisify } = require('util');

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

class AsyncWatchdogClient {
  constructor(serverAddress = 'localhost:50051') {
    this.client = new watchdog.WatchdogService(serverAddress, grpc.credentials.createInsecure());
    
    // Promisify all client methods
    this.getHealth = promisify(this.client.GetHealth.bind(this.client));
    this.registerService = promisify(this.client.RegisterService.bind(this.client));
    this.unregisterService = promisify(this.client.UnregisterService.bind(this.client));
    this.listServices = promisify(this.client.ListServices.bind(this.client));
    this.updateServiceStatus = promisify(this.client.UpdateServiceStatus.bind(this.client));
  }

  // Wrapper methods for better API
  async checkHealth() {
    return await this.getHealth({});
  }

  async register(name, endpoint, type = 'SERVICE_TYPE_UNSPECIFIED') {
    return await this.registerService({ name, endpoint, type });
  }

  async unregister(serviceId) {
    return await this.unregisterService({ service_id: serviceId });
  }

  async list() {
    return await this.listServices({});
  }

  async updateStatus(serviceId, status) {
    return await this.updateServiceStatus({ service_id: serviceId, status });
  }

  // Close the client connection
  close() {
    this.client.close();
  }
}

// Example with error handling and retry logic
class RobustWatchdogClient extends AsyncWatchdogClient {
  constructor(serverAddress = 'localhost:50051', options = {}) {
    super(serverAddress);
    this.maxRetries = options.maxRetries || 3;
    this.retryDelay = options.retryDelay || 1000;
  }

  async withRetry(operation, ...args) {
    let lastError;
    
    for (let attempt = 1; attempt <= this.maxRetries; attempt++) {
      try {
        return await operation(...args);
      } catch (error) {
        lastError = error;
        console.warn(`Attempt ${attempt} failed:`, error.message);
        
        if (attempt < this.maxRetries) {
          console.log(`Retrying in ${this.retryDelay}ms...`);
          await new Promise(resolve => setTimeout(resolve, this.retryDelay));
        }
      }
    }
    
    throw lastError;
  }

  async checkHealthWithRetry() {
    return await this.withRetry(this.checkHealth.bind(this));
  }

  async registerWithRetry(name, endpoint) {
    return await this.withRetry(this.register.bind(this), name, endpoint);
  }

  async listWithRetry() {
    return await this.withRetry(this.list.bind(this));
  }

  async updateStatusWithRetry(serviceId, status) {
    return await this.withRetry(this.updateStatus.bind(this), serviceId, status);
  }

  async unregisterWithRetry(serviceId) {
    return await this.withRetry(this.unregister.bind(this), serviceId);
  }
}

// Example usage with different patterns
async function demonstratePatterns() {
  console.log('=== Node.js gRPC Client Patterns ===\n');

  // Pattern 1: Basic async client
  console.log('Pattern 1: Basic Async Client');
  const basicClient = new AsyncWatchdogClient();
  
  try {
    const health = await basicClient.checkHealth();
    console.log('Health check:', health);

    const registration = await basicClient.register('test-service', 'http://localhost:4000', 'SERVICE_TYPE_HTTP');
    console.log('Service registered:', registration);

    const services = await basicClient.list();
    console.log('Services count:', services.services.length);

    await basicClient.unregister(registration.service_id);
    console.log('Service unregistered successfully');

  } catch (error) {
    console.error('Basic client error:', error.message);
  } finally {
    basicClient.close();
  }

  console.log('\n' + '='.repeat(50) + '\n');

  // Pattern 2: Robust client with retries
  console.log('Pattern 2: Robust Client with Retries');
  const robustClient = new RobustWatchdogClient('localhost:50051', {
    maxRetries: 3,
    retryDelay: 1000
  });

  try {
    const health = await robustClient.checkHealthWithRetry();
    console.log('Health check with retry:', health);

    const registration = await robustClient.registerWithRetry('robust-service', 'http://localhost:5000', 'SERVICE_TYPE_GRPC');
    console.log('Service registered with retry:', registration);

    const services = await robustClient.listWithRetry();
    console.log('Services with retry:', services.services.length);

    await robustClient.updateStatusWithRetry(registration.service_id, 'degraded');
    console.log('Status updated with retry');

    await robustClient.unregisterWithRetry(registration.service_id);
    console.log('Service unregistered with retry');

  } catch (error) {
    console.error('Robust client error:', error.message);
  } finally {
    robustClient.close();
  }
}

// Pattern 3: Service monitoring loop
async function serviceMonitoring() {
  console.log('\nPattern 3: Service Monitoring Loop');
  const client = new AsyncWatchdogClient();
  
  try {
    // Register a service for monitoring
    const registration = await client.register('monitored-service', 'http://localhost:6000', 'SERVICE_TYPE_MICROSERVICE');
    const serviceId = registration.service_id;
    console.log('Started monitoring service:', serviceId);

    // Simulate a monitoring loop
    const statuses = ['healthy', 'degraded', 'healthy', 'down', 'healthy'];
    
    for (const status of statuses) {
      await client.updateStatus(serviceId, status);
      console.log(`Updated service status to: ${status}`);
      
      const services = await client.list();
      const myService = services.services.find(s => s.id === serviceId);
      console.log(`Current status: ${myService.status}, Last heartbeat: ${myService.last_heartbeat}`);
      
      // Wait 2 seconds between updates
      await new Promise(resolve => setTimeout(resolve, 2000));
    }

    // Clean up
    await client.unregister(serviceId);
    console.log('Monitoring completed and service unregistered');

  } catch (error) {
    console.error('Monitoring error:', error.message);
  } finally {
    client.close();
  }
}

// Run demonstrations
async function main() {
  await demonstratePatterns();
  await serviceMonitoring();
}

if (require.main === module) {
  main().catch(console.error);
}

module.exports = { AsyncWatchdogClient, RobustWatchdogClient };