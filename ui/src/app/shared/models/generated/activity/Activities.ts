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

import {
  GithubComPd120424DMountainServiceApiContractsActivityV1ActivityCreateRequest,
  GithubComPd120424DMountainServiceApiContractsActivityV1ActivityListResponse,
  GithubComPd120424DMountainServiceApiContractsActivityV1ActivityResponse,
  GithubComPd120424DMountainServiceApiContractsActivityV1ActivityStatsResponse,
} from "./data-contracts";
import { ContentType, HttpClient, RequestParams } from "./http-client";

export class Activities<SecurityDataType = unknown> extends HttpClient<SecurityDataType> {
  /**
   * @description Get a paginated list of activities with optional filtering
   *
   * @tags activities
   * @name ActivitiesList
   * @summary List activities
   * @request GET:/activities
   * @secure
   */
  activitiesList = (
    query?: {
      /**
       * Page number
       * @default 1
       */
      page?: number;
      /**
       * Page size
       * @default 10
       */
      pageSize?: number;
      /** Activity type filter */
      type?: string;
      /** Activity level filter */
      level?: string;
    },
    params: RequestParams = {},
  ) =>
    this.request<GithubComPd120424DMountainServiceApiContractsActivityV1ActivityListResponse, Record<string, any>>({
      path: `/activities`,
      method: "GET",
      query: query,
      secure: true,
      format: "json",
      ...params,
    });
  /**
   * @description Create a new activity in the system
   *
   * @tags activities
   * @name ActivitiesCreate
   * @summary Create a new activity
   * @request POST:/activities
   * @secure
   */
  activitiesCreate = (
    activity: GithubComPd120424DMountainServiceApiContractsActivityV1ActivityCreateRequest,
    params: RequestParams = {},
  ) =>
    this.request<GithubComPd120424DMountainServiceApiContractsActivityV1ActivityResponse, Record<string, any>>({
      path: `/activities`,
      method: "POST",
      body: activity,
      secure: true,
      type: ContentType.Json,
      format: "json",
      ...params,
    });
  /**
   * @description Delete all activities from the system
   *
   * @tags activities
   * @name ResetDelete
   * @summary Reset all activity data
   * @request DELETE:/activities/reset
   * @secure
   */
  resetDelete = (params: RequestParams = {}) =>
    this.request<Record<string, any>, Record<string, any>>({
      path: `/activities/reset`,
      method: "DELETE",
      secure: true,
      ...params,
    });
  /**
   * @description Get comprehensive activity statistics
   *
   * @tags activities
   * @name StatsList
   * @summary Get activity statistics
   * @request GET:/activities/stats
   * @secure
   */
  statsList = (params: RequestParams = {}) =>
    this.request<GithubComPd120424DMountainServiceApiContractsActivityV1ActivityStatsResponse, Record<string, any>>({
      path: `/activities/stats`,
      method: "GET",
      secure: true,
      format: "json",
      ...params,
    });
  /**
   * @description Get a specific activity by its ID
   *
   * @tags activities
   * @name ActivitiesDetail
   * @summary Get activity by ID
   * @request GET:/activities/{id}
   * @secure
   */
  activitiesDetail = (id: number, params: RequestParams = {}) =>
    this.request<GithubComPd120424DMountainServiceApiContractsActivityV1ActivityResponse, Record<string, any>>({
      path: `/activities/${id}`,
      method: "GET",
      secure: true,
      format: "json",
      ...params,
    });
  /**
   * @description Delete a specific activity by its ID
   *
   * @tags activities
   * @name ActivitiesDelete
   * @summary Delete activity
   * @request DELETE:/activities/{id}
   * @secure
   */
  activitiesDelete = (id: number, params: RequestParams = {}) =>
    this.request<Record<string, any>, Record<string, any>>({
      path: `/activities/${id}`,
      method: "DELETE",
      secure: true,
      ...params,
    });
}
