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

export interface GithubComPd120424DMountainServiceApiContractsEmployeeV1ErrorResponse {
  error?: string;
}

export interface GithubComPd120424DMountainServiceApiContractsEmployeeV1MessageResponse {
  message?: string;
}

export interface InternalHandlerUploadProfilePictureResponse {
  blobName?: string;
  blobUrl?: string;
  message?: string;
  size?: number;
}

export interface V1ActiveEmergenciesResponse {
  hasActiveEmergencies?: boolean;
}

export interface V1AssignShiftRequest {
  shiftDate: string;
  /**
   * @min 1
   * @max 3
   */
  shiftType: number;
}

export interface V1AssignShiftResponse {
  id: number;
  shiftDate: string;
  shiftType: number;
}

export interface V1EmployeeCreateRequest {
  email: string;
  firstName: string;
  gender: string;
  lastName: string;
  password: string;
  phone: string;
  profilePicture?: string;
  profileType?: string;
  username: string;
}

export interface V1EmployeeLogin {
  password?: string;
  username?: string;
}

export interface V1EmployeeResponse {
  email?: string;
  firstName?: string;
  gender?: string;
  id?: number;
  lastName?: string;
  phone?: string;
  /** this may be represented as a byte array if we read the picture from somewhere for an example */
  profilePicture?: string;
  profileType?: string;
  username?: string;
}

export interface V1EmployeeUpdateRequest {
  email?: string;
  firstName?: string;
  gender?: string;
  lastName?: string;
  phone?: string;
  profilePicture?: string;
  profileType?: string;
  username?: string;
}

export interface V1OnCallEmployeesResponse {
  employees?: V1EmployeeResponse[];
}

export interface V1RemoveShiftRequest {
  shiftDate: string;
  /**
   * @min 1
   * @max 3
   */
  shiftType: number;
}

export interface V1ShiftAvailability {
  /** Whether the requesting employee is assigned to this shift */
  isAssignedToEmployee?: boolean;
  /** Whether the shift is at full capacity (2 medics + 4 technicians) */
  isFullyBooked?: boolean;
  /** Available slots for medics (0-2) */
  medicSlotsAvailable?: number;
  /** Available slots for technical staff (0-4) */
  technicalSlotsAvailable?: number;
}

export interface V1ShiftAvailabilityPerDay {
  firstShift?: V1ShiftAvailability;
  secondShift?: V1ShiftAvailability;
  thirdShift?: V1ShiftAvailability;
}

export interface V1ShiftAvailabilityResponse {
  days?: Record<string, V1ShiftAvailabilityPerDay>;
}

export interface V1ShiftResponse {
  createdAt?: string;
  id?: number;
  shiftDate?: string;
  /** 1: 6am-2pm, 2: 2pm-10pm, 3: 10pm-6am, < 1 or > 3: invalid */
  shiftType?: number;
}

export interface V1TokenResponse {
  /** @example "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..." */
  token?: string;
}
