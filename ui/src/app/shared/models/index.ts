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
  MessageResponse
} from './generated/employee';

// Urgency Service
export type {
  UrgencyResponse,
  UrgencyCreateRequest,
  UrgencyUpdateRequest,
  UrgencyLevel as GeneratedUrgencyLevel,
  UrgencyStatus as GeneratedUrgencyStatus
} from './generated/urgency';

// Activity Service
export type {
  ActivityResponse,
  ActivityCreateRequest,
  ActivityListRequest,
  ActivityListResponse,
  ActivityStatsResponse,
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

// Re-export urgency utilities and legacy enums
export {
  UrgencyLevel,
  Status,
  type Urgency,
  mapGeneratedLevelToLegacy,
  mapLegacyLevelToGenerated,
  mapGeneratedStatusToLegacy,
  mapLegacyStatusToGenerated,
  getUrgencyLevelColor,
  getStatusColor
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
