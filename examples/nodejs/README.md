# Node.js gRPC Client for Watchdog Service

This directory contains Node.js examples for calling the Watchdog gRPC service.

## Setup

1. Install dependencies:
```bash
cd examples/nodejs
npm install
```

2. Make sure the Watchdog gRPC server is running:
```bash
# From the project root
make run
```

## Examples

### Basic Client (`client.js`)

A simple client that demonstrates all Watchdog service methods:

```bash
node client.js
```

**Features:**
- Health checking
- Service registration
- Service listing
- Status updates
- Service unregistration
- Promise-based API

### Async Client (`async-client.js`)

Advanced client with multiple patterns:

```bash
node async-client.js
```

**Features:**
- Promisified gRPC methods using `util.promisify`
- Retry logic with exponential backoff
- Service monitoring loops
- Error handling patterns
- Clean resource management

## API Usage

### Creating a Client

```javascript
const grpc = require('@grpc/grpc-js');
const protoLoader = require('@grpc/proto-loader');

// Load proto file
const packageDefinition = protoLoader.loadSync('../../proto/watchdog.proto', {
  keepCase: true,
  longs: String,
  enums: String,
  defaults: true,
  oneofs: true
});

const watchdog = grpc.loadPackageDefinition(packageDefinition).watchdog;
const client = new watchdog.WatchdogService('localhost:50051', grpc.credentials.createInsecure());
```

### Method Examples

#### Health Check
```javascript
client.GetHealth({}, (error, response) => {
  if (error) {
    console.error('Error:', error);
  } else {
    console.log('Health:', response);
    // { status: 'healthy', message: 'Watchdog service is running' }
  }
});
```

#### Register Service
```javascript
const request = {
  name: 'my-service',
  endpoint: 'http://localhost:3000',
  type: 'SERVICE_TYPE_HTTP'  // Service type
};

client.RegisterService(request, (error, response) => {
  if (error) {
    console.error('Error:', error);
  } else {
    console.log('Registered:', response);
    // { service_id: '1', message: 'Service my-service registered successfully with ID 1' }
  }
});
```

#### List Services
```javascript
client.ListServices({}, (error, response) => {
  if (error) {
    console.error('Error:', error);
  } else {
    console.log('Services:', response.services);
    // Array of ServiceInfo objects
  }
});
```

#### Update Service Status
```javascript
const request = {
  service_id: 'my-service-1234567890',
  status: 'healthy'
};

client.UpdateServiceStatus(request, (error, response) => {
  if (error) {
    console.error('Error:', error);
  } else {
    console.log('Updated:', response);
    // { message: 'Service status updated successfully' }
  }
});
```

#### Unregister Service
```javascript
const request = {
  service_id: 'my-service-1234567890'
};

client.UnregisterService(request, (error, response) => {
  if (error) {
    console.error('Error:', error);
  } else {
    console.log('Unregistered:', response);
    // { message: 'Service unregistered successfully' }
  }
});
```

## Advanced Patterns

### Promise-based Client

```javascript
const { promisify } = require('util');

class WatchdogClient {
  constructor() {
    this.client = new watchdog.WatchdogService('localhost:50051', grpc.credentials.createInsecure());
    this.getHealth = promisify(this.client.GetHealth.bind(this.client));
    this.registerService = promisify(this.client.RegisterService.bind(this.client));
    // ... other methods
  }

  async checkHealth() {
    return await this.getHealth({});
  }
}
```

### Error Handling

```javascript
try {
  const response = await client.getHealth();
  console.log('Success:', response);
} catch (error) {
  switch (error.code) {
    case grpc.status.UNAVAILABLE:
      console.error('Service unavailable');
      break;
    case grpc.status.DEADLINE_EXCEEDED:
      console.error('Request timeout');
      break;
    case grpc.status.INVALID_ARGUMENT:
      console.error('Invalid request:', error.message);
      break;
    default:
      console.error('Unknown error:', error);
  }
}
```

### Connection with Metadata

```javascript
const metadata = new grpc.Metadata();
metadata.add('authorization', 'Bearer your-token');

client.GetHealth({}, metadata, (error, response) => {
  // Handle response
});
```

### Connection Options

```javascript
const client = new watchdog.WatchdogService('localhost:50051', grpc.credentials.createInsecure(), {
  'grpc.keepalive_time_ms': 30000,
  'grpc.keepalive_timeout_ms': 5000,
  'grpc.keepalive_permit_without_calls': true,
  'grpc.http2.max_pings_without_data': 0,
  'grpc.http2.min_time_between_pings_ms': 10000
});
```

## Service Info Object

```javascript
{
  id: "1",
  name: "my-web-service", 
  endpoint: "http://localhost:3000",
  type: "SERVICE_TYPE_HTTP",
  status: "healthy",
  last_heartbeat: "1672531200" // Unix timestamp
}
```

## Service Types

The following service types are supported:

- `SERVICE_TYPE_UNSPECIFIED` - Default/unspecified type
- `SERVICE_TYPE_HTTP` - HTTP/REST API services
- `SERVICE_TYPE_GRPC` - gRPC services  
- `SERVICE_TYPE_DATABASE` - Database services (MySQL, PostgreSQL, etc.)
- `SERVICE_TYPE_CACHE` - Cache services (Redis, Memcached, etc.)
- `SERVICE_TYPE_QUEUE` - Message queue services (RabbitMQ, Apache Kafka, etc.)
- `SERVICE_TYPE_STORAGE` - Storage services (S3, MinIO, etc.)
- `SERVICE_TYPE_EXTERNAL_API` - External API integrations
- `SERVICE_TYPE_MICROSERVICE` - Internal microservices
- `SERVICE_TYPE_OTHER` - Other service types

## Common Status Values

- `"active"` - Service is registered but status not updated
- `"healthy"` - Service is running normally
- `"degraded"` - Service has issues but still functional
- `"down"` - Service is not responding
- `"maintenance"` - Service is under maintenance

## Troubleshooting

1. **Connection Refused**: Make sure the gRPC server is running on `localhost:50051`
2. **Proto Loading Issues**: Verify the path to `watchdog.proto` is correct
3. **Method Not Found**: Ensure your proto file matches the server's proto definition

## Dependencies

- `@grpc/grpc-js`: Core gRPC library for Node.js
- `@grpc/proto-loader`: Dynamic proto loading
- `grpc-tools`: Code generation tools (optional)

## Further Reading

- [gRPC Node.js Documentation](https://grpc.github.io/grpc/node/)
- [Protocol Buffers](https://developers.google.com/protocol-buffers)
- [gRPC Core Concepts](https://grpc.io/docs/what-is-grpc/core-concepts/)