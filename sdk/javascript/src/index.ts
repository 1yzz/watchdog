export { WatchdogClient } from './client';
export type {
  WatchdogClientOptions,
  ServiceRegistration,
  ServiceUpdate,
} from './client';
export {
  ServiceType,
  ServiceInfo,
  HealthRequest,
  HealthResponse,
  RegisterServiceRequest,
  RegisterServiceResponse,
  UnregisterServiceRequest,
  UnregisterServiceResponse,
  ListServicesRequest,
  ListServicesResponse,
  UpdateServiceStatusRequest,
  UpdateServiceStatusResponse,
} from './generated/watchdog_pb';