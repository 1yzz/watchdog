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

    // 2. Register a service with unique name
    console.log('\n2. Registering a service...');
    const timestamp = Date.now();
    const serviceName = `test-service-${timestamp}`;
    
    const registration = await client.registerService({
      name: serviceName,
      endpoint: 'http://localhost:8080/health',
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
    
    // Provide more detailed error information
    if (error.message.includes('failed to register service')) {
      console.error('\n💡 Troubleshooting tips:');
      console.error('1. Check if a service with the same name and endpoint already exists');
      console.error('2. Ensure the endpoint is a valid URL format');
      console.error('3. Verify the service type is supported');
      console.error('4. Check server logs for more detailed error information');
    }
  } finally {
    // Always close the client
    client.close();
  }
}

main().catch(console.error);