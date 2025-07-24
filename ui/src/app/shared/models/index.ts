// Barrel exports for shared models
// This file will be updated as we add more generated models

// Generated models - Employee Service
export * from './generated/employee';

// Extensions (manual customizations)
export * from './extensions';

// Re-export commonly used types with cleaner names
export type {
  EmployeeResponse as Employee,
  EmployeeCreateRequest,
  EmployeeUpdateRequest,
  EmployeeLogin,
  TokenResponse,
  ErrorResponse,
  MessageResponse
} from './generated/employee';

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
