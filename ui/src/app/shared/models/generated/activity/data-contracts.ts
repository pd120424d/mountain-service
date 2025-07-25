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

export interface GithubComPd120424DMountainServiceApiContractsActivityV1ActivityCreateRequest {
  actorId?: number;
  actorName?: string;
  description: string;
  level: GithubComPd120424DMountainServiceApiContractsActivityV1ActivityLevel;
  metadata?: string;
  targetId?: number;
  targetType?: string;
  title: string;
  type: GithubComPd120424DMountainServiceApiContractsActivityV1ActivityType;
}

export enum GithubComPd120424DMountainServiceApiContractsActivityV1ActivityLevel {
  ActivityLevelInfo = "info",
  ActivityLevelWarning = "warning",
  ActivityLevelError = "error",
  ActivityLevelCritical = "critical",
}

export interface GithubComPd120424DMountainServiceApiContractsActivityV1ActivityListResponse {
  activities?: GithubComPd120424DMountainServiceApiContractsActivityV1ActivityResponse[];
  page?: number;
  pageSize?: number;
  total?: number;
  totalPages?: number;
}

export interface GithubComPd120424DMountainServiceApiContractsActivityV1ActivityResponse {
  /** ID of the user who performed the action */
  actorId?: number;
  /** Name of the user who performed the action */
  actorName?: string;
  createdAt?: string;
  description?: string;
  id?: number;
  level?: GithubComPd120424DMountainServiceApiContractsActivityV1ActivityLevel;
  /** JSON string with additional data */
  metadata?: string;
  /** ID of the target entity */
  targetId?: number;
  /** Type of the target entity (employee, urgency, etc.) */
  targetType?: string;
  title?: string;
  type?: GithubComPd120424DMountainServiceApiContractsActivityV1ActivityType;
  updatedAt?: string;
}

export interface GithubComPd120424DMountainServiceApiContractsActivityV1ActivityStatsResponse {
  activitiesByLevel?: Record<string, number>;
  activitiesByType?: Record<string, number>;
  activitiesLast24h?: number;
  activitiesLast30Days?: number;
  activitiesLast7Days?: number;
  recentActivities?: GithubComPd120424DMountainServiceApiContractsActivityV1ActivityResponse[];
  totalActivities?: number;
}

export enum GithubComPd120424DMountainServiceApiContractsActivityV1ActivityType {
  ActivityEmployeeCreated = "employee_created",
  ActivityEmployeeUpdated = "employee_updated",
  ActivityEmployeeDeleted = "employee_deleted",
  ActivityEmployeeLogin = "employee_login",
  ActivityShiftAssigned = "shift_assigned",
  ActivityShiftRemoved = "shift_removed",
  ActivityUrgencyCreated = "urgency_created",
  ActivityUrgencyUpdated = "urgency_updated",
  ActivityUrgencyDeleted = "urgency_deleted",
  ActivityEmergencyAssigned = "emergency_assigned",
  ActivityEmergencyAccepted = "emergency_accepted",
  ActivityEmergencyDeclined = "emergency_declined",
  ActivityNotificationSent = "notification_sent",
  ActivityNotificationFailed = "notification_failed",
  ActivitySystemReset = "system_reset",
}
