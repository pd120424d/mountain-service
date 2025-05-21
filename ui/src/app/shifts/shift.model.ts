import { MedicRole, TechnicalRole } from "../employee/employee.model";

export interface ShiftResponse {
  id: number;
  shiftDate: string; // ISO format date
  shiftType: number; // 1: 6am-2pm, 2: 2pm-10pm, 3: 10pm-6am
  createdAt?: string; // Optional
}

export interface AssignShiftRequest {
  shiftDate: string;
  shiftType: number;
}

export interface AssignShiftResponse {
  id: number;
  shiftDate: string;
  shiftType: number;
}

export interface RemoveShiftRequest {
  id: number;
}

export interface ShiftAvailabilityDto {
  medic: number;
  technical: number;
}

export interface ShiftAvailabilityPerDay {
  firstShift: ShiftAvailabilityDto;
  secondShift: ShiftAvailabilityDto;
  thirdShift: ShiftAvailabilityDto;
}

export interface ShiftAvailabilityResponse {
  days: {
    [date: string]: ShiftAvailabilityPerDay;
  };
}

