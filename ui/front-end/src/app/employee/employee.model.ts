// src/app/employee/employee.model.ts
export interface Employee {
    id: number;
    firstName: string;
    lastName: string;
    gender: 'M' | 'F';
    username: string;
    password: string;
    confirmPassword: string;
    phoneNumber: string;
    email: string;
    profilePicture: string; // URL or path to the profile picture
    profileType: 'Medic' | 'Technical Staff';
  }
  