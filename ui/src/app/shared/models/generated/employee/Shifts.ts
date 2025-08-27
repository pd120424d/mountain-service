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

import { InternalHandlerErrorResponse, InternalHandlerShiftAvailabilityResponse } from "./data-contracts";
import { HttpClient, RequestParams } from "./http-client";

export class Shifts<SecurityDataType = unknown> extends HttpClient<SecurityDataType> {
  /**
   * @description Дохватање доступности смена за одређени дан
   *
   * @tags запослени
   * @name AvailabilityList
   * @summary Дохватање доступности смена
   * @request GET:/shifts/availability
   * @secure
   */
  availabilityList = (
    query?: {
      /** Дан за који се проверава доступност смена */
      date?: string;
    },
    params: RequestParams = {},
  ) =>
    this.request<InternalHandlerShiftAvailabilityResponse, InternalHandlerErrorResponse>({
      path: `/shifts/availability`,
      method: "GET",
      query: query,
      secure: true,
      ...params,
    });
}
