// Export the main client class
export { WatchdogClient } from './client';

// Export all client types dynamically
export type * from './client';

// Re-export all generated protobuf types dynamically
export * from './generated/watchdog_pb';