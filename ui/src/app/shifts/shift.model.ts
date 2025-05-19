import { MedicRole, TechnicalRole } from "../employee/employee.model";

export interface ShiftResponse {
    id: number;
    shiftDate: string; // ISO format date
    shiftType: number; // 1: 6am-2pm, 2: 2pm-10pm, 3: 10pm-6am
    createdAt?: string; // Optional
  }
  
  export interface AssignShiftRequest {
    profileType: typeof MedicRole | typeof TechnicalRole;
    shiftDate: string;
    shiftType: number;
  }
  
  export interface AssignShiftResponse {
    id: number;
    shiftDate: string;
    shiftType: number;
    profileType: string;
  }
  
  export interface RemoveShiftRequest {
    id: number;
  }
  