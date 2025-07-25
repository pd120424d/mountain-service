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
