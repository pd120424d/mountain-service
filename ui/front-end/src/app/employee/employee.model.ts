export interface Employee {
  id: number;
  firstName: string;
  lastName: string;
  email: string;
  phoneNumber: string;
  profilePicture?: string; // Optional
  profileType: 'Medic' | 'Technical';
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
  profileType: 'Medic' | 'Technical';
  gender: string;
  profilePicture?: string; // Optional
}

export interface EmployeeUpdateRequest {
  firstName?: string;
  lastName?: string;
  email?: string;
  age?: number; // Optional, per schema
}
