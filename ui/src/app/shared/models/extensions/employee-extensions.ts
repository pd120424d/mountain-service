// Employee model extensions and utilities
// These extend the generated models with frontend-specific functionality

import type { EmployeeResponse } from '../index';

// Role constants for backward compatibility
export const MedicRole = 'Medic' as const;
export const TechnicalRole = 'Technical' as const;
export const AdministratorRole = 'Administrator' as const;

// Type for employee roles
export type EmployeeRole = typeof MedicRole | typeof TechnicalRole | typeof AdministratorRole;

// Extended employee interface with computed properties
export interface EmployeeWithDisplayName extends EmployeeResponse {
  displayName: string;
}

// Utility functions
export const createDisplayName = (employee: EmployeeResponse): string => {
  return `${employee.firstName} ${employee.lastName}`.trim();
};

export const isAdmin = (employee: EmployeeResponse): boolean => {
  return employee.profileType === AdministratorRole;
};

export const isMedic = (employee: EmployeeResponse): boolean => {
  return employee.profileType === MedicRole;
};

export const isTechnical = (employee: EmployeeResponse): boolean => {
  return employee.profileType === TechnicalRole;
};
