/* eslint-disable */
/* tslint:disable */
// @ts-nocheck
/*
 * ---------------------------------------------------------------
 * ## THIS FILE WAS GENERATED VIA SWAGGER-TYPESCRIPT-API        ##
 * ##                                                           ##
 * ## AUTHOR: acacode                                           ##
 * ## SOURCE: https://github.com/acacode/swagger-typescript-api ##
 * ---------------------------------------------------------------
 */

export interface ErrorResponse {
  error?: string;
}

export interface MessageResponse {
  message?: string;
}

export interface ShiftAvailability {
  /** Available slots for medics (0-2) */
  medicSlotsAvailable?: number;
  /** Available slots for technical staff (0-4) */
  technicalSlotsAvailable?: number;
}

export interface ShiftAvailabilityPerDay {
  firstShift?: ShiftAvailability;
  secondShift?: ShiftAvailability;
  thirdShift?: ShiftAvailability;
}

export interface ShiftAvailabilityResponse {
  days?: Record<string, ShiftAvailabilityPerDay>;
}

export interface AssignShiftRequest {
  /** @format date */
  shiftDate: string;
  /**
   * @min 1
   * @max 3
   */
  shiftType: number;
}

export interface AssignShiftResponse {
  /** @format int64 */
  id?: number;
  /** @format date */
  shiftDate?: string;
  shiftType?: number;
}

export interface RemoveShiftRequest {
  /** @format date */
  shiftDate: string;
  /**
   * @min 1
   * @max 3
   */
  shiftType: number;
}

export interface ShiftResponse {
  /** @format int64 */
  id?: number;
  /** @format date-time */
  shiftDate?: string;
  shiftType?: number;
  /** @format date-time */
  createdAt?: string;
}

export interface EmployeeLogin {
  username?: string;
  password?: string;
}

export interface EmployeeResponse {
  /** @format int64 */
  id?: number;
  username?: string;
  firstName?: string;
  lastName?: string;
  gender?: string;
  phone?: string;
  email?: string;
  profilePicture?: string;
  profileType?: "Medic" | "Technical" | "Administrator";
}

export interface EmployeeCreateRequest {
  firstName: string;
  lastName: string;
  username: string;
  password: string;
  /** @format email */
  email: string;
  gender: string;
  phone: string;
  profilePicture?: string;
  profileType?: "Medic" | "Technical" | "Administrator";
}

export interface EmployeeUpdateRequest {
  firstName?: string;
  lastName?: string;
  username?: string;
  /** @format email */
  email?: string;
  gender?: string;
  phone?: string;
  profilePicture?: string;
  profileType?: "Medic" | "Technical" | "Administrator";
}

export interface TokenResponse {
  /** @example "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..." */
  token?: string;
}
