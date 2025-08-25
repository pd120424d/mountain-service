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

export interface ActivityResponse {
  /** @format int64 */
  id?: number;
  description?: string;
  /** @format int64 */
  employee_id?: number;
  /** @format int64 */
  urgency_id?: number;
  /** @format date-time */
  created_at?: string;
  /** @format date-time */
  updated_at?: string;
}

export interface ActivityCreateRequest {
  description: string;
  /** @format int64 */
  employee_id: number;
  /** @format int64 */
  urgency_id: number;
}

export interface ActivityListRequest {
  /** @format int64 */
  employee_id?: number;
  /** @format int64 */
  urgency_id?: number;
  /** @format date-time */
  start_date?: string;
  /** @format date-time */
  end_date?: string;
  page?: number;
  page_size?: number;
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
  total_activities?: number;
  recent_activities?: ActivityResponse[];
  /** @format int64 */
  activities_last_24h?: number;
  /** @format int64 */
  activities_last_7_days?: number;
  /** @format int64 */
  activities_last_30_days?: number;
}
