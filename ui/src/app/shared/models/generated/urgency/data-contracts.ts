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

export enum UrgencyLevel {
  Low = "low",
  Medium = "medium",
  High = "high",
  Critical = "critical",
}

export enum UrgencyStatus {
  Open = "open",
  InProgress = "in_progress",
  Resolved = "resolved",
  Closed = "closed",
}

export interface UrgencyCreateRequest {
  firstName: string;
  lastName: string;
  /** @format email */
  email?: string;
  contactPhone: string;
  location: string;
  description: string;
  level?: UrgencyLevel;
}

export interface UrgencyUpdateRequest {
  firstName?: string;
  lastName?: string;
  /** @format email */
  email?: string;
  contactPhone?: string;
  location?: string;
  description?: string;
  level?: UrgencyLevel;
  status?: UrgencyStatus;
}

export interface UrgencyResponse {
  /** @format int64 */
  id?: number;
  firstName?: string;
  lastName?: string;
  email?: string;
  contactPhone?: string;
  location?: string;
  description?: string;
  level?: UrgencyLevel;
  status?: UrgencyStatus;
  /** @format date-time */
  createdAt?: string;
  /** @format int64 */
  assignedEmployeeId?: number;
  /** @format date-time */
  assignedAt?: string;
  /** @format date-time */
  updatedAt?: string;
}

export interface UrgencyList {
  urgencies?: UrgencyResponse[];
}

export interface EmergencyAssignmentResponse {
  /** @format int64 */
  id?: number;
  /** @format int64 */
  urgencyId?: number;
  /** @format int64 */
  employeeId?: number;
  status?: string;
  /** @format date-time */
  assignedAt?: string;
  /** @format date-time */
  createdAt?: string;
  /** @format date-time */
  updatedAt?: string;
}

export interface AssignmentAcceptRequest {
  /** @format int64 */
  assignmentId: number;
}

export interface AssignmentDeclineRequest {
  /** @format int64 */
  assignmentId: number;
  reason?: string;
}

export interface EmployeeAssignmentsResponse {
  assignments?: EmergencyAssignmentResponse[];
}

export interface NotificationResponse {
  /** @format int64 */
  id?: number;
  /** @format int64 */
  urgencyId?: number;
  /** @format int64 */
  employeeId?: number;
  notificationType?: string;
  recipient?: string;
  message?: string;
  status?: string;
  attempts?: number;
  /** @format date-time */
  lastAttemptAt?: string;
  /** @format date-time */
  sentAt?: string;
  errorMessage?: string;
  /** @format date-time */
  createdAt?: string;
  /** @format date-time */
  updatedAt?: string;
}
