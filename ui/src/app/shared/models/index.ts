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
  EmployeeResponse as Employee,
  EmployeeCreateRequest,
  EmployeeUpdateRequest,
  EmployeeLogin,
  TokenResponse,
  ErrorResponse,
  MessageResponse,
  ShiftAvailability,
  ShiftAvailabilityPerDay,
  ShiftAvailabilityResponse,
  AssignShiftRequest,
  AssignShiftResponse,
  RemoveShiftRequest,
  ShiftResponse
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
  getUrgencyLevelColor,
  getUrgencyStatusColor
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
