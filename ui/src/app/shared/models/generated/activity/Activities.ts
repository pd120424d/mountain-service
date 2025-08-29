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



import { HttpClient, RequestParams, ContentType, HttpResponse } from "./http-client";
import { GithubComPd120424DMountainServiceApiContractsActivityV1ErrorResponse, V1ActivityCreateRequest, V1ActivityListResponse, V1ActivityResponse, V1ActivityStatsResponse } from "./data-contracts"

export class Activities<SecurityDataType = unknown> extends HttpClient<SecurityDataType>  {

            /**
 * @description Извлачење листе активности са опционим филтрирањем и страничењем
 *
 * @tags activities
 * @name ActivitiesList
 * @summary Листа активности
 * @request GET:/activities
 * @secure
 */
activitiesList: (query?: {
  /**
   * Page number
   * @default 1
   */
    page?: number,
  /**
   * Page size
   * @default 10
   */
    pageSize?: number,
  /** Activity type filter */
    type?: string,
  /** Activity level filter */
    level?: string,

}, params: RequestParams = {}) =>
    this.request<V1ActivityListResponse, Record<string,any>>({
        path: `/activities`,
        method: 'GET',
        query: query,                secure: true,                format: "json",        ...params,
    }),            /**
 * @description Креирање нове активности у систему
 *
 * @tags activities
 * @name ActivitiesCreate
 * @summary Креирање нове активности
 * @request POST:/activities
 * @secure
 */
activitiesCreate: (activity: V1ActivityCreateRequest, params: RequestParams = {}) =>
    this.request<V1ActivityResponse, GithubComPd120424DMountainServiceApiContractsActivityV1ErrorResponse>({
        path: `/activities`,
        method: 'POST',
                body: activity,        secure: true,        type: ContentType.Json,        format: "json",        ...params,
    }),            /**
 * @description Брисање свих активности из система
 *
 * @tags activities
 * @name ResetDelete
 * @summary Ресетовање свих података о активностима
 * @request DELETE:/activities/reset
 * @secure
 */
resetDelete: (params: RequestParams = {}) =>
    this.request<Record<string,any>, Record<string,any>>({
        path: `/activities/reset`,
        method: 'DELETE',
                        secure: true,                        ...params,
    }),            /**
 * @description Преузимање свеобухватних статистика активности
 *
 * @tags activities
 * @name StatsList
 * @summary Статистике активности
 * @request GET:/activities/stats
 * @secure
 */
statsList: (params: RequestParams = {}) =>
    this.request<V1ActivityStatsResponse, Record<string,any>>({
        path: `/activities/stats`,
        method: 'GET',
                        secure: true,                format: "json",        ...params,
    }),            /**
 * @description Преузимање одређене активности по њеном ID
 *
 * @tags activities
 * @name ActivitiesDetail
 * @summary Преузимање активности по ID
 * @request GET:/activities/{id}
 * @secure
 */
activitiesDetail: (id: number, params: RequestParams = {}) =>
    this.request<V1ActivityResponse, Record<string,any>>({
        path: `/activities/${id}`,
        method: 'GET',
                        secure: true,                format: "json",        ...params,
    }),            /**
 * @description Брисање одређене активности по њеном ID
 *
 * @tags activities
 * @name ActivitiesDelete
 * @summary Брисање активности
 * @request DELETE:/activities/{id}
 * @secure
 */
activitiesDelete: (id: number, params: RequestParams = {}) =>
    this.request<Record<string,any>, Record<string,any>>({
        path: `/activities/${id}`,
        method: 'DELETE',
                        secure: true,                        ...params,
    }),    }
