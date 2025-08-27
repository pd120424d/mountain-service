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

export interface GithubComPd120424DMountainServiceApiContractsEmployeeV1EmployeeResponse {
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

export interface GithubComPd120424DMountainServiceApiContractsEmployeeV1ShiftAvailability {
  /** Whether the requesting employee is assigned to this shift */
  isAssignedToEmployee?: boolean;
  /** Whether the shift is at full capacity (2 medics + 4 technicians) */
  isFullyBooked?: boolean;
  /** Available slots for medics (0-2) */
  medicSlotsAvailable?: number;
  /** Available slots for technical staff (0-4) */
  technicalSlotsAvailable?: number;
}

export interface GithubComPd120424DMountainServiceApiContractsEmployeeV1ShiftAvailabilityPerDay {
  firstShift?: GithubComPd120424DMountainServiceApiContractsEmployeeV1ShiftAvailability;
  secondShift?: GithubComPd120424DMountainServiceApiContractsEmployeeV1ShiftAvailability;
  thirdShift?: GithubComPd120424DMountainServiceApiContractsEmployeeV1ShiftAvailability;
}

export interface GithubComPd120424DMountainServiceApiEmployeeInternalHandlerActiveEmergenciesResponse {
  hasActiveEmergencies?: boolean;
}

export interface GithubComPd120424DMountainServiceApiEmployeeInternalHandlerAssignShiftRequest {
  shiftDate: string;
  /**
   * @min 1
   * @max 3
   */
  shiftType: number;
}

export interface GithubComPd120424DMountainServiceApiEmployeeInternalHandlerAssignShiftResponse {
  id: number;
  shiftDate: string;
  shiftType: number;
}

export interface GithubComPd120424DMountainServiceApiEmployeeInternalHandlerEmployeeCreateRequest {
  email: string;
  firstName: string;
  gender: string;
  lastName: string;
  password: string;
  phone: string;
  profilePicture?: string;
  profileType: string;
  username: string;
}

export interface GithubComPd120424DMountainServiceApiEmployeeInternalHandlerEmployeeLogin {
  password: string;
  username: string;
}

export interface GithubComPd120424DMountainServiceApiEmployeeInternalHandlerEmployeeResponse {
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

export interface GithubComPd120424DMountainServiceApiEmployeeInternalHandlerEmployeeUpdateRequest {
  email?: string;
  firstName?: string;
  gender?: string;
  lastName?: string;
  phone?: string;
  profilePicture?: string;
  profileType?: string;
  username?: string;
}

export interface GithubComPd120424DMountainServiceApiEmployeeInternalHandlerErrorResponse {
  error?: string;
}

export interface GithubComPd120424DMountainServiceApiEmployeeInternalHandlerMessageResponse {
  message?: string;
}

export interface GithubComPd120424DMountainServiceApiEmployeeInternalHandlerOnCallEmployeesResponse {
  employees?: GithubComPd120424DMountainServiceApiContractsEmployeeV1EmployeeResponse[];
}

export interface GithubComPd120424DMountainServiceApiEmployeeInternalHandlerRemoveShiftRequest {
  shiftDate: string;
  /**
   * @min 1
   * @max 3
   */
  shiftType: number;
}

export interface GithubComPd120424DMountainServiceApiEmployeeInternalHandlerShiftAvailabilityResponse {
  days?: Record<string, GithubComPd120424DMountainServiceApiContractsEmployeeV1ShiftAvailabilityPerDay>;
}

export interface GithubComPd120424DMountainServiceApiEmployeeInternalHandlerShiftResponse {
  createdAt?: string;
  id?: number;
  shiftDate?: string;
  /** 1: 6am-2pm, 2: 2pm-10pm, 3: 10pm-6am, < 1 or > 3: invalid */
  shiftType?: number;
}

export interface GithubComPd120424DMountainServiceApiEmployeeInternalHandlerTokenResponse {
  /** @example "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..." */
  token?: string;
}

export interface GithubComPd120424DMountainServiceApiEmployeeInternalHandlerUploadProfilePictureResponse {
  blobName?: string;
  blobUrl?: string;
  message?: string;
  size?: number;
}

export interface InternalHandlerActiveEmergenciesResponse {
  hasActiveEmergencies?: boolean;
}

export interface InternalHandlerAssignShiftRequest {
  shiftDate: string;
  /**
   * @min 1
   * @max 3
   */
  shiftType: number;
}

export interface InternalHandlerAssignShiftResponse {
  id: number;
  shiftDate: string;
  shiftType: number;
}

export interface InternalHandlerEmployeeCreateRequest {
  email: string;
  firstName: string;
  gender: string;
  lastName: string;
  password: string;
  phone: string;
  profilePicture?: string;
  profileType: string;
  username: string;
}

export interface InternalHandlerEmployeeLogin {
  password: string;
  username: string;
}

export interface InternalHandlerEmployeeResponse {
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

export interface InternalHandlerEmployeeUpdateRequest {
  email?: string;
  firstName?: string;
  gender?: string;
  lastName?: string;
  phone?: string;
  profilePicture?: string;
  profileType?: string;
  username?: string;
}

export interface InternalHandlerErrorResponse {
  error?: string;
}

export interface InternalHandlerMessageResponse {
  message?: string;
}

export interface InternalHandlerOnCallEmployeesResponse {
  employees?: GithubComPd120424DMountainServiceApiContractsEmployeeV1EmployeeResponse[];
}

export interface InternalHandlerRemoveShiftRequest {
  shiftDate: string;
  /**
   * @min 1
   * @max 3
   */
  shiftType: number;
}

export interface InternalHandlerShiftAvailabilityResponse {
  days?: Record<string, GithubComPd120424DMountainServiceApiContractsEmployeeV1ShiftAvailabilityPerDay>;
}

export interface InternalHandlerShiftResponse {
  createdAt?: string;
  id?: number;
  shiftDate?: string;
  /** 1: 6am-2pm, 2: 2pm-10pm, 3: 10pm-6am, < 1 or > 3: invalid */
  shiftType?: number;
}

export interface InternalHandlerTokenResponse {
  /** @example "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..." */
  token?: string;
}

export interface InternalHandlerUploadProfilePictureResponse {
  blobName?: string;
  blobUrl?: string;
  message?: string;
  size?: number;
}
