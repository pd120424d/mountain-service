// Barrel exports for shared models
// This file will be updated as we add more generated models

// Generated models - Employee Service
export * from './generated/employee';

// Generated models - Urgency Service
export * from './generated/urgency';

// Generated models - Activity Service
export * from './generated/activity';

// Extensions (manual customizations)
export * from './extensions';

// Re-export commonly used types with cleaner names

// Employee Service
export type {
  V1EmployeeResponse as Employee,
  V1EmployeeResponse as EmployeeResponse,
  V1EmployeeCreateRequest as EmployeeCreateRequest,
  V1EmployeeUpdateRequest as EmployeeUpdateRequest,
  V1EmployeeLogin as EmployeeLogin,
  V1TokenResponse as TokenResponse,
  GithubComPd120424DMountainServiceApiContractsEmployeeV1ErrorResponse as ErrorResponse,
  GithubComPd120424DMountainServiceApiContractsEmployeeV1MessageResponse as MessageResponse,
  V1ShiftAvailability as ShiftAvailability,
  V1ShiftAvailabilityPerDay as ShiftAvailabilityPerDay,
  V1ShiftAvailabilityResponse as ShiftAvailabilityResponse,
  V1AssignShiftRequest as AssignShiftRequest,
  V1AssignShiftResponse as AssignShiftResponse,
  V1RemoveShiftRequest as RemoveShiftRequest,
  V1ShiftResponse as ShiftResponse,
  InternalHandlerUploadProfilePictureResponse as UploadProfilePictureResponse
} from './generated/employee';

// Urgency Service
export type {
  UrgencyResponse,
  UrgencyCreateRequest,
  UrgencyUpdateRequest
} from './generated/urgency';

export {
  UrgencyLevel,
  UrgencyStatus,
  UrgencyLevel as GeneratedUrgencyLevel,
  UrgencyStatus as GeneratedUrgencyStatus
} from './generated/urgency';

// Activity Service
export type {
  ActivityResponse,
  ActivityCreateRequest,
  ActivityListResponse,
  ActivityStatsResponse
} from './generated/activity';

export {
  ActivityType,
  ActivityLevel
} from './generated/activity';

// Re-export role constants and utilities for backward compatibility
export {
  MedicRole,
  TechnicalRole,
  AdministratorRole,
  type EmployeeRole,
  type EmployeeWithDisplayName,
  createDisplayName,
  isAdmin,
  isMedic,
  isTechnical
} from './extensions/employee-extensions';

// Legacy enum exports for backward compatibility
export enum EmployeeResponseProfileTypeEnum {
  Medic = "Medic",
  Technical = "Technical",
  Administrator = "Administrator",
}

export enum EmployeeCreateRequestProfileTypeEnum {
  Medic = "Medic",
  Technical = "Technical",
  Administrator = "Administrator",
}

export enum EmployeeUpdateRequestProfileTypeEnum {
  Medic = "Medic",
  Technical = "Technical",
  Administrator = "Administrator",
}

// Re-export urgency utilities
export {
  type Urgency,
  type UrgencyWithDisplayName,
  type LocationCoordinates,
  type EnhancedLocation,
  type UrgencyCreateRequestWithCoordinates,
  type UrgencyResponseWithCoordinates,
  createUrgencyDisplayName,
  withDisplayName,
  getUrgencyLevelColor,
  getUrgencyStatusColor,
  parseLocationString,
  formatLocationForApi,
  formatCoordinatesDisplay,
  calculateDistance,
  isValidCoordinates,
  isInMountainRegion
} from './extensions/urgency-extensions';

// Re-export activity utilities
export {
  type Activity,
  getActivityLevelColor,
  getActivityTypeIcon,
  getActivityTypeDisplayName,
  isSystemActivity,
  isEmployeeActivity,
  isUrgencyActivity,
  isShiftActivity,
  isNotificationActivity
} from './extensions/activity-extensions';
