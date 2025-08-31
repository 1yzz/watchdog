import { WatchdogClient, ServiceType, WatchdogClientOptions } from '../dist/index';

async function demonstrateWatchdogSDK(): Promise<void> {
  const options: WatchdogClientOptions = {
    host: 'localhost',
    port: 50051,
    timeout: 5000,
  };

  const client = new WatchdogClient(options);

  try {
    console.log('🚀 Starting Watchdog SDK TypeScript Demo\n');

    // Health check
    const health = await client.getHealth();
    console.log('Health status:', health.status);
    console.log('Health message:', health.message);

    // Register multiple services
    const services = [
      {
        name: 'api-gateway',
        endpoint: 'http://api.example.com:8080',
        type: ServiceType.SERVICE_TYPE_HTTP,
      },
      {
        name: 'user-service',
        endpoint: 'grpc://user-service:9090',
        type: ServiceType.SERVICE_TYPE_GRPC,
      },
      {
        name: 'postgres-db',
        endpoint: 'postgres://db:5432/app',
        type: ServiceType.SERVICE_TYPE_DATABASE,
      },
      {
        name: 'redis-cache',
        endpoint: 'redis://cache:6379',
        type: ServiceType.SERVICE_TYPE_CACHE,
      },
    ];

    const registeredServices: string[] = [];

    // Register all services
    for (const service of services) {
      console.log(`\n📝 Registering ${service.name}...`);
      const result = await client.registerService(service);
      console.log(`✅ Registered with ID: ${result.serviceId}`);
      registeredServices.push(result.serviceId);
    }

    // List all services
    console.log('\n📋 Current services:');
    const allServices = await client.listServices();
    allServices.forEach((service) => {
      console.log(`- ${service.getName()} (${service.getEndpoint()}) - ${service.getStatus()}`);
    });

    // Update service statuses
    const statuses = ['healthy', 'warning', 'unhealthy', 'healthy'];
    for (let i = 0; i < registeredServices.length; i++) {
      console.log(`\n🔄 Updating ${services[i].name} status to: ${statuses[i]}`);
      await client.updateServiceStatus({
        serviceId: registeredServices[i],
        status: statuses[i],
      });
    }

    // Show updated services
    console.log('\n📊 Services after status updates:');
    const updatedServices = await client.listServices();
    updatedServices.forEach((service) => {
      const statusEmoji = service.getStatus() === 'healthy' ? '✅' : 
                         service.getStatus() === 'warning' ? '⚠️' : '❌';
      console.log(`${statusEmoji} ${service.getName()}: ${service.getStatus()}`);
    });

    // Clean up - unregister all services
    console.log('\n🧹 Cleaning up...');
    for (const serviceId of registeredServices) {
      await client.unregisterService(serviceId);
      console.log(`🗑️  Unregistered service ${serviceId}`);
    }

    console.log('\n🎉 Demo completed successfully!');

  } catch (error) {
    console.error('💥 Error during demo:', error);
  } finally {
    client.close();
    console.log('🔌 Client connection closed');
  }
}

// Run the demo
if (require.main === module) {
  demonstrateWatchdogSDK().catch(console.error);
}

export { demonstrateWatchdogSDK };