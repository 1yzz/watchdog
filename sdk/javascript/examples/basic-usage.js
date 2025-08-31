const { WatchdogClient, ServiceType } = require('../dist/index');

async function main() {
  const client = new WatchdogClient({
    host: 'localhost',
    port: 50051,
    timeout: 5000,
  });

  try {
    console.log('Testing Watchdog gRPC SDK...\n');

    // 1. Check health
    console.log('1. Checking health...');
    const health = await client.getHealth();
    console.log('Health status:', health);

    // 2. Register a service
    console.log('\n2. Registering a service...');
    const registration = await client.registerService({
      name: 'my-web-service',
      endpoint: 'http://localhost:8080',
      type: ServiceType.SERVICE_TYPE_HTTP,
    });
    console.log('Registration:', registration);

    const serviceId = registration.serviceId;

    // 3. List services
    console.log('\n3. Listing all services...');
    const services = await client.listServices();
    console.log('Services:', services.map(s => ({
      id: s.getId(),
      name: s.getName(),
      endpoint: s.getEndpoint(),
      type: s.getType(),
      status: s.getStatus(),
    })));

    // 4. Update service status
    console.log('\n4. Updating service status...');
    const statusUpdate = await client.updateServiceStatus({
      serviceId: serviceId,
      status: 'healthy',
    });
    console.log('Status update:', statusUpdate);

    // 5. List services again to see the update
    console.log('\n5. Listing services after status update...');
    const updatedServices = await client.listServices();
    console.log('Updated services:', updatedServices.map(s => ({
      id: s.getId(),
      name: s.getName(),
      status: s.getStatus(),
    })));

    // 6. Unregister the service
    console.log('\n6. Unregistering service...');
    const unregistration = await client.unregisterService(serviceId);
    console.log('Unregistration:', unregistration);

    console.log('\n✅ All operations completed successfully!');

  } catch (error) {
    console.error('❌ Error:', error.message);
  } finally {
    // Always close the client
    client.close();
  }
}

main().catch(console.error);