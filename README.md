# Watchdog gRPC Service

A gRPC-based service monitoring and management system built with Go. The Watchdog service allows you to register, monitor, and manage the health status of various services in your infrastructure.

## Features

- **Service Registration**: Register services with the watchdog for monitoring
- **Health Monitoring**: Track service health and status
- **Service Discovery**: List all registered services
- **Status Management**: Update service status and heartbeat information
- **gRPC API**: High-performance gRPC interface
- **MySQL Persistence**: All data persisted to MySQL database
- **Service History**: Track status changes over time
- **Health Check Logging**: Monitor server health metrics
- **Reflection Support**: Built-in gRPC reflection for easy testing
- **Native Deployment**: Simple binary deployment with .env configuration
- **Ent Framework**: Schema-as-code with automatic migrations using Facebook's Ent
- **JavaScript/TypeScript SDK**: Full-featured SDK with dynamic protobuf support
- **Automated SDK Publishing**: Streamlined release process with version management

## Project Structure

```
watchdog/
├── api/                    # Generated protobuf code
│   ├── watchdog.pb.go      # Protocol buffer message definitions
│   └── watchdog_grpc.pb.go # gRPC service definitions
├── proto/                  # Protocol buffer schema
│   └── watchdog.proto      # API definition
├── ent/                    # Ent framework (schema-as-code)
│   ├── schema/            # Entity schema definitions
│   │   └── service.go     # Service entity schema
│   ├── generate.go        # Code generation script
│   ├── client.go          # Generated Ent client
│   ├── service*.go        # Generated service entity operations
│   └── migrate/           # Auto-migration support
├── server/                 # Server implementation
│   └── server.go          # Service methods implementation
├── database/               # Database layer
│   ├── ent_client.go      # Ent-based database client
│   └── interface.go       # Database interface and types
├── config/                 # Configuration management
│   └── config.go          # Configuration loading
├── sdk/                    # Official SDKs
│   ├── javascript/        # JavaScript/TypeScript SDK
│   └── proto/             # SDK proto files
├── examples/               # Client examples
│   └── nodejs/            # Node.js client implementation
├── cmd/                    # Application entry point
│   └── main.go            # Server startup code
├── bin/                    # Built binaries (created after build)
├── scripts/                # Utility scripts
│   ├── migrate-ent.go     # Ent-based migration script
│   └── deploy.sh          # Deployment script
├── docs/                   # Documentation
│   └── configuration.md   # Configuration guide
├── go.mod                  # Go module definition
├── Makefile               # Build automation
└── README.md              # This file
```

## API Reference

### Service Methods

#### GetHealth
Returns the health status of the watchdog service.

**Request**: `HealthRequest` (empty)
**Response**: `HealthResponse`
- `status` (string): Health status
- `message` (string): Status message

#### RegisterService
Registers a new service for monitoring.

**Request**: `RegisterServiceRequest`
- `name` (string): Service name
- `endpoint` (string): Service endpoint URL
- `type` (ServiceType): Service type classification

**Response**: `RegisterServiceResponse`
- `service_id` (string): Generated unique service ID
- `message` (string): Registration confirmation

#### UnregisterService
Removes a service from monitoring.

**Request**: `UnregisterServiceRequest`
- `service_id` (string): Service ID to unregister

**Response**: `UnregisterServiceResponse`
- `message` (string): Unregistration confirmation

#### ListServices
Lists all registered services.

**Request**: `ListServicesRequest` (empty)
**Response**: `ListServicesResponse`
- `services` (array): List of ServiceInfo objects

#### UpdateServiceStatus
Updates the status of a registered service.

**Request**: `UpdateServiceStatusRequest`
- `service_id` (string): Service ID
- `status` (string): New status

**Response**: `UpdateServiceStatusResponse`
- `message` (string): Update confirmation

### Data Types

#### ServiceInfo
- `id` (string): Unique service identifier
- `name` (string): Service name
- `endpoint` (string): Service endpoint
- `type` (ServiceType): Service type classification
- `status` (string): Current service status
- `last_heartbeat` (int64): Unix timestamp of last update

#### ServiceType Enum
- `SERVICE_TYPE_UNSPECIFIED` - Default/unspecified type
- `SERVICE_TYPE_HTTP` - HTTP/REST API services
- `SERVICE_TYPE_GRPC` - gRPC services
- `SERVICE_TYPE_DATABASE` - Database services
- `SERVICE_TYPE_CACHE` - Cache services (Redis, etc.)
- `SERVICE_TYPE_QUEUE` - Message queue services
- `SERVICE_TYPE_STORAGE` - Storage services
- `SERVICE_TYPE_EXTERNAL_API` - External API integrations
- `SERVICE_TYPE_MICROSERVICE` - Internal microservices
- `SERVICE_TYPE_OTHER` - Other service types
- `SERVICE_TYPE_SYSTEMD` - SystemD managed services

## Prerequisites

- Go 1.19 or later
- MySQL 8.0 or later
- Make (for build automation)
- Protocol Buffers compiler (protoc) - optional, for regenerating proto files

## Installation

1. **Clone the repository:**
```bash
git clone <repository-url>
cd watchdog
```

2. **Install dependencies:**
```bash
make deps
```

3. **Create and configure environment file:**
```bash
make env-setup
# Edit .env with your database credentials
```

4. **Set up MySQL database:**

**Option 1: Ent-based auto-migration (recommended)**
```bash
# First, create database and user manually:
mysql -u root -p
CREATE DATABASE IF NOT EXISTS watchdog_db CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;
CREATE USER IF NOT EXISTS 'watchdog'@'%' IDENTIFIED BY 'watchdog123';
GRANT ALL PRIVILEGES ON watchdog_db.* TO 'watchdog'@'%';
FLUSH PRIVILEGES;
exit

# Then run Ent-based migration (automatically creates/updates schema):
make db-migrate-ent

# Or preview what would be created without executing:
make db-migrate-ent-dry
```

**Option 2: Manual schema setup**
```bash
# Create database and user manually, then run Ent migration
mysql -u root -p
CREATE DATABASE IF NOT EXISTS watchdog_db CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;
CREATE USER IF NOT EXISTS 'watchdog'@'%' IDENTIFIED BY 'watchdog123';
GRANT ALL PRIVILEGES ON watchdog_db.* TO 'watchdog'@'%';
FLUSH PRIVILEGES;
exit

make db-migrate-ent
```

5. **Test the setup:**
```bash
# Test by running the migration in dry-run mode
make db-migrate-ent-dry
```

6. **Build the server:**
```bash
make build
```

## Usage

### Starting the Server

**Development mode:**
```bash
make run
```

**Production mode (with custom .env):**
```bash
./bin/watchdog-server
```

**Background service:**
```bash
nohup ./bin/watchdog-server > watchdog.log 2>&1 &
```

**With systemd (recommended for production):**
```bash
# Copy binary to system location
sudo cp bin/watchdog-server /usr/local/bin/
sudo cp .env /etc/watchdog/

# Create systemd service (see deployment section)
sudo systemctl start watchdog
```

The server will start on port `50051` by default (configurable via `PORT` environment variable).

### Available Make Commands

**Development:**
- `make build` - Build the server binary
- `make run` - Build and run the server
- `make clean` - Clean build artifacts
- `make test` - Run tests
- `make deps` - Install dependencies
- `make fmt` - Format code
- `make lint` - Run linter (requires golangci-lint)
- `make proto` - Generate protobuf code (requires protoc)
- `make build-all` - Build for multiple platforms

**Database:**
- `make db-migrate-ent` - Run Ent-based database migration (recommended)
- `make db-migrate-ent-dry` - Show what migration would do without executing
- `make ent-generate` - Generate Ent code from schema definitions

**JavaScript SDK:**
- `make sdk-install` - Install SDK dependencies
- `make sdk-build` - Build the JavaScript SDK
- `make sdk-test` - Run SDK tests
- `make sdk-lint` - Lint SDK code
- `make sdk-release-patch` - Release patch version (1.0.x)
- `make sdk-release-minor` - Release minor version (1.x.0)
- `make sdk-release-major` - Release major version (x.0.0)
- `make sdk-clean` - Clean SDK build artifacts

**Other:**
- `make help` - Show available commands

### Testing with grpcurl

If you have [grpcurl](https://github.com/fullstorydev/grpcurl) installed, you can test the API:

```bash
# Check health
grpcurl -plaintext localhost:50051 watchdog.WatchdogService/GetHealth

# Register a service
grpcurl -plaintext -d '{"name": "my-service", "endpoint": "http://localhost:8080", "type": "SERVICE_TYPE_HTTP"}' \
  localhost:50051 watchdog.WatchdogService/RegisterService

# List services
grpcurl -plaintext localhost:50051 watchdog.WatchdogService/ListServices

# Update service status
grpcurl -plaintext -d '{"service_id": "my-service-123456", "status": "healthy"}' \
  localhost:50051 watchdog.WatchdogService/UpdateServiceStatus

# Unregister service
grpcurl -plaintext -d '{"service_id": "my-service-123456"}' \
  localhost:50051 watchdog.WatchdogService/UnregisterService
```

### JavaScript/TypeScript SDK

The official JavaScript SDK provides a modern, type-safe interface with dynamic protobuf support:

#### Installation

```bash
npm install watchdog-grpc-sdk
```

#### Basic Usage

```javascript
const { WatchdogClient, ServiceType } = require('watchdog-grpc-sdk');

async function main() {
  const client = new WatchdogClient({
    host: 'localhost',
    port: 50051,
    timeout: 5000,
  });

  try {
    // Check health
    const health = await client.getHealth();
    console.log('Health:', health);

    // Register a service
    const registration = await client.registerService({
      name: 'my-web-service',
      endpoint: 'http://localhost:8080',
      type: ServiceType.SERVICE_TYPE_HTTP,
    });
    console.log('Registered:', registration);

    // List all services
    const services = await client.listServices();
    services.forEach(service => {
      console.log(`${service.getName()}: ${service.getStatus()}`);
    });

    // Update service status
    await client.updateServiceStatus({
      serviceId: registration.serviceId,
      status: 'healthy',
    });

    // Unregister service
    await client.unregisterService(registration.serviceId);
  } catch (error) {
    console.error('Error:', error.message);
  } finally {
    client.close();
  }
}

main();
```

#### TypeScript Usage

```typescript
import { WatchdogClient, ServiceType, WatchdogClientOptions } from 'watchdog-grpc-sdk';

const options: WatchdogClientOptions = {
  host: 'localhost',
  port: 50051,
  timeout: 5000,
};

const client = new WatchdogClient(options);

// Full type safety for all operations
const registration = await client.registerService({
  name: 'user-service',
  endpoint: 'grpc://user-service:9090',
  type: ServiceType.SERVICE_TYPE_GRPC,
});
```

#### SDK Features

- **Dynamic Protobuf Support**: No hardcoded imports - automatically adapts to proto changes
- **Full TypeScript Support**: Complete type definitions and IntelliSense support
- **Promise-based API**: Modern async/await patterns
- **Automatic Type Generation**: Generated from the same proto files as the server
- **Error Handling**: Descriptive error messages with troubleshooting hints
- **Connection Management**: Automatic connection lifecycle management

#### SDK Development

Build the SDK from source:

```bash
# Build the SDK
make sdk-build

# Run tests
make sdk-test

# Lint code
make sdk-lint

# Run examples
cd sdk/javascript
node examples/basic-usage.js
```

### Node.js Client Examples

Legacy Node.js client examples are available in the `examples/nodejs/` directory:

```bash
# Install dependencies
cd examples/nodejs
npm install

# Run the basic client example
node client.js

# Run advanced async client examples  
node async-client.js
```

See [examples/nodejs/README.md](examples/nodejs/README.md) for detailed usage instructions.

## Development

### Regenerating Protocol Buffers

If you modify the proto files, regenerate the Go code and SDK:

```bash
# Regenerate Go protobuf code
make proto

# Rebuild JavaScript SDK with new proto definitions
make sdk-build
```

Note: This requires the Protocol Buffers compiler (`protoc`) and the Go plugins to be installed.

### SDK Development Workflow

The JavaScript SDK uses dynamic protobuf imports to eliminate hardcoded dependencies:

1. **Modify proto files**: Update `proto/watchdog.proto`
2. **Regenerate code**: Run `make proto` to update Go code
3. **Rebuild SDK**: Run `make sdk-build` to regenerate JavaScript types
4. **Test changes**: Run `make sdk-test` to verify functionality
5. **Version and publish**: Use `make sdk-release-patch/minor/major` for releases

The SDK automatically adapts to proto changes without manual code updates.

### Code Formatting

Format your code before committing:

```bash
make fmt
```

### Running Tests

```bash
make test
```

## Deployment

### Systemd Service (Linux)

Create a systemd service for production deployment:

1. **Create service directory:**
```bash
sudo mkdir -p /etc/watchdog
sudo mkdir -p /var/log/watchdog
```

2. **Copy files:**
```bash
sudo cp bin/watchdog-server /usr/local/bin/
sudo cp .env /etc/watchdog/
sudo chown -R watchdog:watchdog /etc/watchdog
```

3. **Create service file:**
```bash
sudo tee /etc/systemd/system/watchdog.service > /dev/null <<EOF
[Unit]
Description=Watchdog gRPC Service
After=network.target mysql.service
Wants=mysql.service

[Service]
Type=simple
User=watchdog
Group=watchdog
WorkingDirectory=/etc/watchdog
ExecStart=/usr/local/bin/watchdog-server
Restart=always
RestartSec=5
StandardOutput=append:/var/log/watchdog/watchdog.log
StandardError=append:/var/log/watchdog/watchdog.log

# Security settings
NoNewPrivileges=yes
PrivateTmp=yes
ProtectSystem=strict
ProtectHome=yes
ReadWritePaths=/var/log/watchdog

[Install]
WantedBy=multi-user.target
EOF
```

4. **Start and enable service:**
```bash
sudo systemctl daemon-reload
sudo systemctl enable watchdog
sudo systemctl start watchdog
sudo systemctl status watchdog
```

### Process Manager (PM2)

For Node.js environments, you can use PM2:

```bash
# Install PM2
npm install -g pm2

# Create PM2 configuration
cat > ecosystem.config.js <<EOF
module.exports = {
  apps: [{
    name: 'watchdog',
    script: './bin/watchdog-server',
    instances: 1,
    exec_mode: 'fork',
    watch: false,
    env: {
      NODE_ENV: 'production'
    }
  }]
}
EOF

# Start with PM2
pm2 start ecosystem.config.js
pm2 save
pm2 startup
```

### Manual Deployment

**Production binary:**
```bash
make build-all
# Copy appropriate binary to target server
# Configure .env file
# Run: ./watchdog-server
```

**Background process:**
```bash
nohup ./bin/watchdog-server > watchdog.log 2>&1 &
echo $! > watchdog.pid
```

**Stop background process:**
```bash
kill $(cat watchdog.pid)
rm watchdog.pid
```

## Configuration

### Environment Variables

The server automatically loads configuration from `.env` files in this order of precedence:
1. `.env.local` (highest priority, git-ignored)
2. `.env` (git-ignored)
3. System environment variables
4. Default values (lowest priority)

The server can be configured using the following environment variables:

| Variable | Default | Description |
|----------|---------|-------------|
| `PORT` | `50051` | gRPC server port |
| `DB_HOST` | `localhost` | MySQL host address |
| `DB_PORT` | `3306` | MySQL port |
| `DB_USERNAME` | `watchdog` | MySQL username |
| `DB_PASSWORD` | `watchdog123` | MySQL password |
| `DB_DATABASE` | `watchdog_db` | MySQL database name |

### Setup .env File

Create your environment file:
```bash
make env-setup
```

This creates a `.env` file from the template. Edit it with your specific configuration:
```bash
# .env file
PORT=50051
DB_HOST=your-mysql-host.com
DB_USERNAME=your-username
DB_PASSWORD=your-password
DB_DATABASE=watchdog_db
```

### Configuration Examples

**Local Development (.env file):**
```bash
PORT=50051
DB_HOST=localhost
DB_PORT=3306
DB_USERNAME=watchdog
DB_PASSWORD=watchdog123
DB_DATABASE=watchdog_db
```

**AWS RDS (.env file):**
```bash
PORT=50051
DB_HOST=myinstance.123456789012.us-east-1.rds.amazonaws.com
DB_PORT=3306
DB_USERNAME=admin
DB_PASSWORD=mypassword
DB_DATABASE=watchdog_db
```

**Google Cloud SQL (.env file):**
```bash
PORT=50051
DB_HOST=10.1.2.3
DB_PORT=3306
DB_USERNAME=watchdog-user
DB_PASSWORD=secure-password
DB_DATABASE=watchdog_db
```

**Remote MySQL Server (.env file):**
```bash
PORT=50051
DB_HOST=mysql.example.com
DB_PORT=3306
DB_USERNAME=watchdog
DB_PASSWORD=watchdog123
DB_DATABASE=watchdog_db
```

**Environment Variables (alternative to .env file):**
```bash
export DB_HOST=your-mysql-host.com
export DB_USERNAME=your-username
export DB_PASSWORD=your-password
# ... etc
```

## Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Run tests and formatting
5. Submit a pull request

## License

This project is licensed under the MIT License - see the LICENSE file for details.

## Support

For questions or issues, please open an issue on the project repository.