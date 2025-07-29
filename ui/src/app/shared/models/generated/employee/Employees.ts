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
  AssignShiftRequest,
  AssignShiftResponse,
  ErrorResponse,
  RemoveShiftRequest,
  ShiftResponse,
} from "./data-contracts";
import { ContentType, HttpClient, RequestParams } from "./http-client";

export class Employees<
  SecurityDataType = unknown,
> extends HttpClient<SecurityDataType> {
  /**
   * @description Assigns a shift to an employee
   *
   * @name ShiftsCreate
   * @summary Assign shift to employee
   * @request POST:/employees/{id}/shifts
   */
  shiftsCreate = (
    id: number,
    body: AssignShiftRequest,
    params: RequestParams = {},
  ) =>
    this.request<AssignShiftResponse, ErrorResponse>({
      path: `/employees/${id}/shifts`,
      method: "POST",
      body: body,
      type: ContentType.Json,
      format: "json",
      ...params,
    });
  /**
   * @description Gets all shifts for an employee
   *
   * @name ShiftsList
   * @summary Get employee shifts
   * @request GET:/employees/{id}/shifts
   */
  shiftsList = (id: number, params: RequestParams = {}) =>
    this.request<ShiftResponse[], any>({
      path: `/employees/${id}/shifts`,
      method: "GET",
      format: "json",
      ...params,
    });
  /**
   * @description Removes an employee from a shift
   *
   * @name ShiftsDelete
   * @summary Remove employee from shift
   * @request DELETE:/employees/{id}/shifts
   */
  shiftsDelete = (
    id: number,
    body: RemoveShiftRequest,
    params: RequestParams = {},
  ) =>
    this.request<void, ErrorResponse>({
      path: `/employees/${id}/shifts`,
      method: "DELETE",
      body: body,
      type: ContentType.Json,
      ...params,
    });
  /**
   * @description Returns warnings about shift coverage and quota for an employee
   *
   * @name ShiftWarningsList
   * @summary Get shift warnings for employee
   * @request GET:/employees/{id}/shift-warnings
   */
  shiftWarningsList = (id: number, params: RequestParams = {}) =>
    this.request<
      {
        warnings?: string[];
      },
      ErrorResponse
    >({
      path: `/employees/${id}/shift-warnings`,
      method: "GET",
      format: "json",
      ...params,
    });
}
