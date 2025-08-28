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

export interface GithubComPd120424DMountainServiceApiContractsActivityV1ErrorResponse {
  error?: string;
}

export interface V1ActivityCreateRequest {
  description: string;
  employeeId: number;
  urgencyId: number;
}

export interface V1ActivityListResponse {
  activities?: V1ActivityResponse[];
  page?: number;
  pageSize?: number;
  total?: number;
  totalPages?: number;
}

export interface V1ActivityResponse {
  createdAt?: string;
  description?: string;
  /** ID of the employee who created the activity */
  employeeId?: number;
  id?: number;
  updatedAt?: string;
  /** ID of the urgency this activity relates to */
  urgencyId?: number;
}

export interface V1ActivityStatsResponse {
  activitiesLast24h?: number;
  activitiesLast30Days?: number;
  activitiesLast7Days?: number;
  recentActivities?: V1ActivityResponse[];
  totalActivities?: number;
}
