# Watchdog gRPC SDK for JavaScript/Node.js

A TypeScript/JavaScript SDK for interacting with the Watchdog gRPC service. This SDK provides a high-level, promise-based API for service registration, monitoring, and management.

## Features

- 🚀 **Promise-based API** - Modern async/await support
- 📝 **Full TypeScript Support** - Complete type definitions included
- 🔧 **Auto-generated from Protobuf** - Always in sync with server API
- ⚡ **High Performance** - Built on @grpc/grpc-js for optimal performance
- 🛡️ **Error Handling** - Comprehensive error handling with descriptive messages
- 🔌 **Easy Integration** - Simple to integrate into any Node.js project

## Installation

```bash
npm install watchdog-grpc-sdk
```

## Quick Start

### JavaScript (CommonJS)

```javascript
const { WatchdogClient, ServiceType } = require('watchdog-grpc-sdk');

async function main() {
  const client = new WatchdogClient({
    host: 'localhost',
    port: 50051,
  });

  try {
    // Check service health
    const health = await client.getHealth();
    console.log('Watchdog is', health.status);

    // Register a service
    const result = await client.registerService({
      name: 'my-api',
      endpoint: 'http://localhost:8080',
      type: ServiceType.SERVICE_TYPE_HTTP,
    });
    console.log('Service registered with ID:', result.serviceId);

    // List all services
    const services = await client.listServices();
    console.log('Total services:', services.length);

  } catch (error) {
    console.error('Error:', error.message);
  } finally {
    client.close();
  }
}

main();
```

### TypeScript

```typescript
import { WatchdogClient, ServiceType, WatchdogClientOptions } from 'watchdog-grpc-sdk';

const options: WatchdogClientOptions = {
  host: 'localhost',
  port: 50051,
  timeout: 5000, // 5 second timeout
};

const client = new WatchdogClient(options);

// Register a service with full type safety
const registration = await client.registerService({
  name: 'user-service',
  endpoint: 'grpc://user-service:9090',
  type: ServiceType.SERVICE_TYPE_GRPC,
});
```

## API Reference

### WatchdogClient

#### Constructor

```typescript
new WatchdogClient(options: WatchdogClientOptions)
```

**Options:**
- `host: string` - Watchdog server hostname
- `port: number` - Watchdog server port
- `credentials?: grpc.ChannelCredentials` - gRPC credentials (default: insecure)
- `timeout?: number` - Request timeout in milliseconds (default: 5000)

#### Methods

##### getHealth()

Check the health status of the Watchdog service.

```typescript
const health = await client.getHealth();
// Returns: { status: string, message: string }
```

##### registerService(service)

Register a new service for monitoring.

```typescript
const result = await client.registerService({
  name: 'api-gateway',
  endpoint: 'http://gateway.example.com:8080',
  type: ServiceType.SERVICE_TYPE_HTTP,
});
// Returns: { serviceId: string, message: string }
```

##### unregisterService(serviceId)

Remove a service from monitoring.

```typescript
const result = await client.unregisterService('service-123');
// Returns: { message: string }
```

##### listServices()

Get all registered services.

```typescript
const services = await client.listServices();
// Returns: ServiceInfo[]

// Access service properties:
services.forEach(service => {
  console.log(service.getId());        // string
  console.log(service.getName());      // string
  console.log(service.getEndpoint());  // string
  console.log(service.getType());      // ServiceType
  console.log(service.getStatus());    // string
});
```

##### updateServiceStatus(update)

Update the status of a registered service.

```typescript
const result = await client.updateServiceStatus({
  serviceId: 'service-123',
  status: 'healthy',
});
// Returns: { message: string }
```

##### close()

Close the gRPC client connection.

```typescript
client.close();
```

### Service Types

The SDK supports the following service types:

```typescript
enum ServiceType {
  SERVICE_TYPE_UNSPECIFIED = 0,
  SERVICE_TYPE_HTTP = 1,
  SERVICE_TYPE_GRPC = 2,
  SERVICE_TYPE_DATABASE = 3,
  SERVICE_TYPE_CACHE = 4,
  SERVICE_TYPE_QUEUE = 5,
  SERVICE_TYPE_STORAGE = 6,
  SERVICE_TYPE_EXTERNAL_API = 7,
  SERVICE_TYPE_MICROSERVICE = 8,
  SERVICE_TYPE_OTHER = 9,
}
```

## Examples

### Basic Service Monitoring

```javascript
const { WatchdogClient, ServiceType } = require('watchdog-grpc-sdk');

class ServiceMonitor {
  constructor() {
    this.client = new WatchdogClient({
      host: process.env.WATCHDOG_HOST || 'localhost',
      port: parseInt(process.env.WATCHDOG_PORT) || 50051,
    });
    this.serviceId = null;
  }

  async start() {
    try {
      // Register this service
      const result = await this.client.registerService({
        name: 'my-application',
        endpoint: 'http://localhost:3000',
        type: ServiceType.SERVICE_TYPE_HTTP,
      });
      
      this.serviceId = result.serviceId;
      console.log('Service registered:', this.serviceId);

      // Start heartbeat
      this.startHeartbeat();
    } catch (error) {
      console.error('Failed to register service:', error.message);
    }
  }

  startHeartbeat() {
    setInterval(async () => {
      if (this.serviceId) {
        try {
          await this.client.updateServiceStatus({
            serviceId: this.serviceId,
            status: 'healthy',
          });
          console.log('Heartbeat sent');
        } catch (error) {
          console.error('Heartbeat failed:', error.message);
        }
      }
    }, 30000); // Every 30 seconds
  }

  async stop() {
    if (this.serviceId) {
      await this.client.unregisterService(this.serviceId);
      console.log('Service unregistered');
    }
    this.client.close();
  }
}

// Usage
const monitor = new ServiceMonitor();
monitor.start();

// Graceful shutdown
process.on('SIGINT', async () => {
  await monitor.stop();
  process.exit(0);
});
```

### Service Discovery

```javascript
const { WatchdogClient } = require('watchdog-grpc-sdk');

async function discoverServices() {
  const client = new WatchdogClient({
    host: 'localhost',
    port: 50051,
  });

  try {
    const services = await client.listServices();
    
    // Group services by type
    const servicesByType = {};
    services.forEach(service => {
      const type = service.getType();
      if (!servicesByType[type]) {
        servicesByType[type] = [];
      }
      servicesByType[type].push({
        id: service.getId(),
        name: service.getName(),
        endpoint: service.getEndpoint(),
        status: service.getStatus(),
      });
    });

    console.log('Services by type:', servicesByType);
  } finally {
    client.close();
  }
}

discoverServices();
```

### Error Handling

```typescript
import { WatchdogClient, ServiceType } from 'watchdog-grpc-sdk';

async function robustServiceRegistration() {
  const client = new WatchdogClient({
    host: 'localhost',
    port: 50051,
    timeout: 3000,
  });

  const maxRetries = 3;
  let retries = 0;

  while (retries < maxRetries) {
    try {
      const result = await client.registerService({
        name: 'critical-service',
        endpoint: 'http://localhost:8080',
        type: ServiceType.SERVICE_TYPE_HTTP,
      });

      console.log('✅ Service registered successfully:', result.serviceId);
      break;

    } catch (error) {
      retries++;
      console.error(`❌ Registration attempt ${retries} failed:`, error.message);
      
      if (retries >= maxRetries) {
        console.error('🚨 Max retries reached, giving up');
        throw error;
      }

      // Wait before retrying
      await new Promise(resolve => setTimeout(resolve, 1000 * retries));
    }
  }

  client.close();
}
```

## Development

### Building from Source

```bash
git clone <repository-url>
cd watchdog/sdk/javascript
npm install
npm run build
```

### Running Examples

```bash
# Make sure Watchdog server is running on localhost:50051
npm run build

# Run JavaScript example
node examples/basic-usage.js

# Run TypeScript example (requires ts-node)
npx ts-node examples/typescript-usage.ts
```

### Available Scripts

- `npm run build` - Build the SDK (protobuf generation + TypeScript compilation)
- `npm run proto:generate` - Generate protobuf code only
- `npm run compile` - Compile TypeScript only
- `npm test` - Run tests
- `npm run lint` - Run ESLint
- `npm run docs` - Generate API documentation

## Requirements

- Node.js 14.0.0 or later
- Watchdog gRPC server running and accessible

## License

MIT License - see LICENSE file for details.

## Support

For questions and support:
- Check the [main Watchdog documentation](../../README.md)
- Open an issue on the project repository
- Review the examples in the `/examples` directory