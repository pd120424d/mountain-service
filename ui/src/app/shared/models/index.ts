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
  InternalHandlerEmployeeResponse as Employee,
  InternalHandlerEmployeeResponse as EmployeeResponse,

  InternalHandlerEmployeeUpdateRequest as EmployeeUpdateRequest,
  InternalHandlerEmployeeLogin as EmployeeLogin,
  InternalHandlerTokenResponse as TokenResponse,
  InternalHandlerErrorResponse as ErrorResponse,
  InternalHandlerMessageResponse as MessageResponse,
  V1ShiftAvailability as ShiftAvailability,
  V1ShiftAvailabilityPerDay as ShiftAvailabilityPerDay,
  InternalHandlerShiftAvailabilityResponse as ShiftAvailabilityResponse,
  InternalHandlerAssignShiftRequest as AssignShiftRequest,
  InternalHandlerAssignShiftResponse as AssignShiftResponse,
  InternalHandlerRemoveShiftRequest as RemoveShiftRequest,
  InternalHandlerShiftResponse as ShiftResponse,
  InternalHandlerUploadProfilePictureResponse as UploadProfilePictureResponse
} from './generated/employee';

// Urgency Service
export type {
  V1UrgencyResponse as UrgencyResponse,
  V1UrgencyCreateRequest as UrgencyCreateRequest,
  V1UrgencyUpdateRequest as UrgencyUpdateRequest
} from './generated/urgency';

export {
  V1UrgencyLevel as UrgencyLevel,
  V1UrgencyStatus as UrgencyStatus,
  V1UrgencyLevel as GeneratedUrgencyLevel,
  V1UrgencyStatus as GeneratedUrgencyStatus
} from './generated/urgency';


// Activity Service
export type {
  V1ActivityResponse as ActivityResponse,
  V1ActivityListResponse as ActivityListResponse
} from './generated/activity';

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
  isInMountainRegion,
  hasAcceptedAssignment
} from './extensions/urgency-extensions';

export {
  type Activity,
  getActivityIcon,
  getActivityDisplayTime,
  formatActivityDescription
} from './extensions/activity-extensions';
