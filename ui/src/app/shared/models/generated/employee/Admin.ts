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
  GithubComPd120424DMountainServiceApiEmployeeInternalHandlerErrorResponse,
  GithubComPd120424DMountainServiceApiEmployeeInternalHandlerMessageResponse,
  GithubComPd120424DMountainServiceApiEmployeeInternalHandlerShiftAvailabilityResponse,
} from "./data-contracts";
import { HttpClient, RequestParams } from "./http-client";

export class Admin<SecurityDataType = unknown> extends HttpClient<SecurityDataType> {
  /**
   * @description Брише све запослене, смене и повезане податке из система (само за админе)
   *
   * @tags админ
   * @name ResetDelete
   * @summary Ресетовање свих података
   * @request DELETE:/admin/reset
   * @secure
   */
  resetDelete = (params: RequestParams = {}) =>
    this.request<
      GithubComPd120424DMountainServiceApiEmployeeInternalHandlerMessageResponse,
      GithubComPd120424DMountainServiceApiEmployeeInternalHandlerErrorResponse
    >({
      path: `/admin/reset`,
      method: "DELETE",
      secure: true,
      format: "json",
      ...params,
    });
  /**
   * @description Дохватање доступности смена за све запослене (само за админе)
   *
   * @tags админ
   * @name ShiftsAvailabilityList
   * @summary Дохватање доступности смена за админе
   * @request GET:/admin/shifts/availability
   * @secure
   */
  shiftsAvailabilityList = (
    query?: {
      /** Број дана за које се проверава доступност (подразумевано 7) */
      days?: number;
    },
    params: RequestParams = {},
  ) =>
    this.request<
      GithubComPd120424DMountainServiceApiEmployeeInternalHandlerShiftAvailabilityResponse,
      GithubComPd120424DMountainServiceApiEmployeeInternalHandlerErrorResponse
    >({
      path: `/admin/shifts/availability`,
      method: "GET",
      query: query,
      secure: true,
      ...params,
    });
}
