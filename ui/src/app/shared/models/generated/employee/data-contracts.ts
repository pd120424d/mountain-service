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

export interface GithubComPd120424DMountainServiceApiContractsEmployeeV1ActiveEmergenciesResponse {
  hasActiveEmergencies?: boolean;
}

export interface GithubComPd120424DMountainServiceApiContractsEmployeeV1AssignShiftRequest {
  shiftDate: string;
  /**
   * @min 1
   * @max 3
   */
  shiftType: number;
}

export interface GithubComPd120424DMountainServiceApiContractsEmployeeV1AssignShiftResponse {
  id: number;
  shiftDate: string;
  shiftType: number;
}

export interface GithubComPd120424DMountainServiceApiContractsEmployeeV1EmployeeCreateRequest {
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

export interface GithubComPd120424DMountainServiceApiContractsEmployeeV1EmployeeLogin {
  password?: string;
  username?: string;
}

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

export interface GithubComPd120424DMountainServiceApiContractsEmployeeV1EmployeeUpdateRequest {
  email?: string;
  firstName?: string;
  gender?: string;
  lastName?: string;
  phone?: string;
  profilePicture?: string;
  profileType?: string;
  username?: string;
}

export interface GithubComPd120424DMountainServiceApiContractsEmployeeV1ErrorResponse {
  error?: string;
}

export interface GithubComPd120424DMountainServiceApiContractsEmployeeV1MessageResponse {
  message?: string;
}

export interface GithubComPd120424DMountainServiceApiContractsEmployeeV1OnCallEmployeesResponse {
  employees?: GithubComPd120424DMountainServiceApiContractsEmployeeV1EmployeeResponse[];
}

export interface GithubComPd120424DMountainServiceApiContractsEmployeeV1RemoveShiftRequest {
  shiftDate: string;
  /**
   * @min 1
   * @max 3
   */
  shiftType: number;
}

export interface GithubComPd120424DMountainServiceApiContractsEmployeeV1ShiftAvailability {
  available?: boolean;
  employees?: string[];
}

export interface GithubComPd120424DMountainServiceApiContractsEmployeeV1ShiftAvailabilityPerDay {
  shift1?: GithubComPd120424DMountainServiceApiContractsEmployeeV1ShiftAvailability;
  shift2?: GithubComPd120424DMountainServiceApiContractsEmployeeV1ShiftAvailability;
  shift3?: GithubComPd120424DMountainServiceApiContractsEmployeeV1ShiftAvailability;
}

export interface GithubComPd120424DMountainServiceApiContractsEmployeeV1ShiftAvailabilityResponse {
  days?: Record<string, GithubComPd120424DMountainServiceApiContractsEmployeeV1ShiftAvailabilityPerDay>;
}

export interface GithubComPd120424DMountainServiceApiContractsEmployeeV1ShiftResponse {
  createdAt?: string;
  id?: number;
  shiftDate?: string;
  /** 1: 6am-2pm, 2: 2pm-10pm, 3: 10pm-6am, < 1 or > 3: invalid */
  shiftType?: number;
}

export interface GithubComPd120424DMountainServiceApiContractsEmployeeV1TokenResponse {
  /** @example "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..." */
  token?: string;
}
