export interface Employee {
  id: number;
  firstName: string;
  lastName: string;
  email: string;
  phone: string;
  profilePicture?: string; // Optional
  profileType: EmployeeRole;
  username: string;
  gender: string;
}

export interface EmployeeCreateRequest {
  firstName: string;
  lastName: string;
  email: string;
  phone: string;
  username: string;
  password: string;
  profileType: typeof MedicRole | typeof TechnicalRole;
  gender: string;
  profilePicture?: string; // Optional
}

export interface EmployeeUpdateRequest {
  firstName?: string;
  lastName?: string;
  email?: string;
  age?: number; // Optional, per schema
}

export type EmployeeRole = typeof MedicRole | typeof TechnicalRole | typeof AdministratorRole;

export const MedicRole = 'Medic';
export const TechnicalRole = 'Technical';
export const AdministratorRole = 'Administrator';
