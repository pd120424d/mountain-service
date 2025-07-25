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

export interface GithubComPd120424DMountainServiceApiContractsUrgencyV1UrgencyCreateRequest {
  contactPhone: string;
  description: string;
  email: string;
  level?: GithubComPd120424DMountainServiceApiContractsUrgencyV1UrgencyLevel;
  location: string;
  name: string;
}

export enum GithubComPd120424DMountainServiceApiContractsUrgencyV1UrgencyLevel {
  Low = "low",
  Medium = "medium",
  High = "high",
  Critical = "critical",
}

export interface GithubComPd120424DMountainServiceApiContractsUrgencyV1UrgencyResponse {
  contactPhone?: string;
  createdAt?: string;
  description?: string;
  email?: string;
  id?: number;
  level?: GithubComPd120424DMountainServiceApiContractsUrgencyV1UrgencyLevel;
  location?: string;
  name?: string;
  status?: GithubComPd120424DMountainServiceApiContractsUrgencyV1UrgencyStatus;
  updatedAt?: string;
}

export enum GithubComPd120424DMountainServiceApiContractsUrgencyV1UrgencyStatus {
  Open = "open",
  InProgress = "in_progress",
  Resolved = "resolved",
  Closed = "closed",
}

export interface GithubComPd120424DMountainServiceApiContractsUrgencyV1UrgencyUpdateRequest {
  contactPhone?: string;
  description?: string;
  email?: string;
  level?: GithubComPd120424DMountainServiceApiContractsUrgencyV1UrgencyLevel;
  location?: string;
  name?: string;
  status?: GithubComPd120424DMountainServiceApiContractsUrgencyV1UrgencyStatus;
}
