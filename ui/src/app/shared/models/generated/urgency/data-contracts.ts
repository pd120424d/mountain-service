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

export enum V1UrgencyStatus {
  Open = "open",
  InProgress = "in_progress",
  Resolved = "resolved",
  Closed = "closed",
}

export enum V1UrgencyLevel {
  Low = "low",
  Medium = "medium",
  High = "high",
  Critical = "critical",
}

export interface V1AssignmentCreateRequest {
  employeeId: number;
}

export interface V1UrgencyCreateRequest {
  contactPhone: string;
  description: string;
  email?: string;
  firstName: string;
  lastName: string;
  level?: V1UrgencyLevel;
  location: string;
}

export interface V1UrgencyResponse {
  assignedAt?: string;
  assignedEmployeeId?: number;
  contactPhone?: string;
  createdAt?: string;
  description?: string;
  email?: string;
  firstName?: string;
  id?: number;
  lastName?: string;
  level?: V1UrgencyLevel;
  location?: string;
  status?: V1UrgencyStatus;
  updatedAt?: string;
}

export interface V1UrgencyUpdateRequest {
  contactPhone?: string;
  description?: string;
  email?: string;
  firstName?: string;
  lastName?: string;
  level?: V1UrgencyLevel;
  location?: string;
  status?: V1UrgencyStatus;
}
