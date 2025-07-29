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

import { ShiftAvailabilityResponse } from "./data-contracts";
import { HttpClient, RequestParams } from "./http-client";

export class Shifts<
  SecurityDataType = unknown,
> extends HttpClient<SecurityDataType> {
  /**
   * @description Returns shift availability for the specified number of days
   *
   * @name AvailabilityList
   * @summary Get shift availability
   * @request GET:/shifts/availability
   */
  availabilityList = (
    query?: {
      /**
       * Number of days to get availability for
       * @default 7
       */
      days?: number;
    },
    params: RequestParams = {},
  ) =>
    this.request<ShiftAvailabilityResponse, any>({
      path: `/shifts/availability`,
      method: "GET",
      query: query,
      format: "json",
      ...params,
    });
}
