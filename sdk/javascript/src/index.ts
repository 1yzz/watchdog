export { WatchdogClient } from './client';
export type {
  WatchdogClientOptions,
  ServiceRegistration,
  ServiceUpdate,
} from './client';

// Re-export all generated protobuf types dynamically
export * from './generated/watchdog_pb';