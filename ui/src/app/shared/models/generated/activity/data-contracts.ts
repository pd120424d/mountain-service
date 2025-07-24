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

export enum ActivityLevel {
  Info = "info",
  Warning = "warning",
  Error = "error",
  Critical = "critical",
}

export enum ActivityType {
  EmployeeCreated = "employee_created",
  EmployeeUpdated = "employee_updated",
  EmployeeDeleted = "employee_deleted",
  EmployeeLogin = "employee_login",
  ShiftAssigned = "shift_assigned",
  ShiftRemoved = "shift_removed",
  UrgencyCreated = "urgency_created",
  UrgencyUpdated = "urgency_updated",
  UrgencyDeleted = "urgency_deleted",
  EmergencyAssigned = "emergency_assigned",
  EmergencyAccepted = "emergency_accepted",
  EmergencyDeclined = "emergency_declined",
  NotificationSent = "notification_sent",
  NotificationFailed = "notification_failed",
  SystemReset = "system_reset",
}

export interface ErrorResponse {
  error?: string;
}

export interface MessageResponse {
  message?: string;
}

export interface ActivityResponse {
  /** @format int64 */
  id?: number;
  type?: ActivityType;
  level?: ActivityLevel;
  title?: string;
  description?: string;
  /** @format int64 */
  actorId?: number;
  actorName?: string;
  /** @format int64 */
  targetId?: number;
  targetType?: string;
  metadata?: string;
  /** @format date-time */
  createdAt?: string;
  /** @format date-time */
  updatedAt?: string;
}

export interface ActivityCreateRequest {
  type: ActivityType;
  level: ActivityLevel;
  title: string;
  description: string;
  /** @format int64 */
  actorId?: number;
  actorName?: string;
  /** @format int64 */
  targetId?: number;
  targetType?: string;
  metadata?: string;
}

export interface ActivityListRequest {
  type?: ActivityType;
  level?: ActivityLevel;
  /** @format int64 */
  actorId?: number;
  /** @format int64 */
  targetId?: number;
  targetType?: string;
  /** @format date-time */
  startDate?: string;
  /** @format date-time */
  endDate?: string;
  page?: number;
  pageSize?: number;
}

export interface ActivityListResponse {
  activities?: ActivityResponse[];
  /** @format int64 */
  total?: number;
  page?: number;
  pageSize?: number;
  totalPages?: number;
}

export interface ActivityStatsResponse {
  /** @format int64 */
  totalActivities?: number;
  activitiesByType?: Record<string, number>;
  activitiesByLevel?: Record<string, number>;
  recentActivities?: ActivityResponse[];
  /** @format int64 */
  activitiesLast24h?: number;
  /** @format int64 */
  activitiesLast7Days?: number;
  /** @format int64 */
  activitiesLast30Days?: number;
}
