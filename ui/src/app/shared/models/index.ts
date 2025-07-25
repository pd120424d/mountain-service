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
  GithubComPd120424DMountainServiceApiContractsUrgencyV1UrgencyResponse as UrgencyResponse,
  GithubComPd120424DMountainServiceApiContractsUrgencyV1UrgencyCreateRequest as UrgencyCreateRequest,
  GithubComPd120424DMountainServiceApiContractsUrgencyV1UrgencyUpdateRequest as UrgencyUpdateRequest
} from './generated/urgency';

export {
  GithubComPd120424DMountainServiceApiContractsUrgencyV1UrgencyLevel as GeneratedUrgencyLevel,
  GithubComPd120424DMountainServiceApiContractsUrgencyV1UrgencyStatus as GeneratedUrgencyStatus
} from './generated/urgency';

// Activity Service
export type {
  GithubComPd120424DMountainServiceApiContractsActivityV1ActivityResponse as ActivityResponse,
  GithubComPd120424DMountainServiceApiContractsActivityV1ActivityCreateRequest as ActivityCreateRequest,
  GithubComPd120424DMountainServiceApiContractsActivityV1ActivityListResponse as ActivityListResponse,
  GithubComPd120424DMountainServiceApiContractsActivityV1ActivityStatsResponse as ActivityStatsResponse
} from './generated/activity';

export {
  GithubComPd120424DMountainServiceApiContractsActivityV1ActivityType as ActivityType,
  GithubComPd120424DMountainServiceApiContractsActivityV1ActivityLevel as ActivityLevel
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
