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
  GithubComPd120424DMountainServiceApiContractsEmployeeV1ActiveEmergenciesResponse,
  GithubComPd120424DMountainServiceApiContractsEmployeeV1AssignShiftRequest,
  GithubComPd120424DMountainServiceApiContractsEmployeeV1AssignShiftResponse,
  GithubComPd120424DMountainServiceApiContractsEmployeeV1EmployeeCreateRequest,
  GithubComPd120424DMountainServiceApiContractsEmployeeV1EmployeeResponse,
  GithubComPd120424DMountainServiceApiContractsEmployeeV1EmployeeUpdateRequest,
  GithubComPd120424DMountainServiceApiContractsEmployeeV1ErrorResponse,
  GithubComPd120424DMountainServiceApiContractsEmployeeV1OnCallEmployeesResponse,
  GithubComPd120424DMountainServiceApiContractsEmployeeV1RemoveShiftRequest,
  GithubComPd120424DMountainServiceApiContractsEmployeeV1ShiftResponse,
} from "./data-contracts";
import { ContentType, HttpClient, RequestParams } from "./http-client";

export class Employees<SecurityDataType = unknown> extends HttpClient<SecurityDataType> {
  /**
   * @description Преузимање свих запослених
   *
   * @tags запослени
   * @name EmployeesList
   * @summary Преузимање листе запослених
   * @request GET:/employees
   * @secure
   */
  employeesList = (params: RequestParams = {}) =>
    this.request<GithubComPd120424DMountainServiceApiContractsEmployeeV1EmployeeResponse[][], any>({
      path: `/employees`,
      method: "GET",
      secure: true,
      format: "json",
      ...params,
    });
  /**
   * @description Креирање новог запосленог у систему
   *
   * @tags запослени
   * @name EmployeesCreate
   * @summary Креирање новог запосленог
   * @request POST:/employees
   * @secure
   */
  employeesCreate = (
    employee: GithubComPd120424DMountainServiceApiContractsEmployeeV1EmployeeCreateRequest,
    params: RequestParams = {},
  ) =>
    this.request<
      GithubComPd120424DMountainServiceApiContractsEmployeeV1EmployeeResponse,
      GithubComPd120424DMountainServiceApiContractsEmployeeV1ErrorResponse
    >({
      path: `/employees`,
      method: "POST",
      body: employee,
      secure: true,
      type: ContentType.Json,
      format: "json",
      ...params,
    });
  /**
   * @description Враћа листу запослених који су тренутно на дужности, са опционим бафером у случају да се близу крај тренутне смене
   *
   * @tags запослени
   * @name OnCallList
   * @summary Претрага запослених који су тренутно на дужности
   * @request GET:/employees/on-call
   * @secure
   */
  onCallList = (
    query?: {
      /** Бафер време пре краја смене (нпр. '1h') */
      shift_buffer?: string;
    },
    params: RequestParams = {},
  ) =>
    this.request<
      GithubComPd120424DMountainServiceApiContractsEmployeeV1OnCallEmployeesResponse,
      GithubComPd120424DMountainServiceApiContractsEmployeeV1ErrorResponse
    >({
      path: `/employees/on-call`,
      method: "GET",
      query: query,
      secure: true,
      type: ContentType.Json,
      format: "json",
      ...params,
    });
  /**
   * @description Ажурирање запосленог по ID-ју
   *
   * @tags запослени
   * @name EmployeesUpdate
   * @summary Ажурирање запосленог
   * @request PUT:/employees/{id}
   * @secure
   */
  employeesUpdate = (
    id: number,
    employee: GithubComPd120424DMountainServiceApiContractsEmployeeV1EmployeeUpdateRequest,
    params: RequestParams = {},
  ) =>
    this.request<
      GithubComPd120424DMountainServiceApiContractsEmployeeV1EmployeeResponse,
      GithubComPd120424DMountainServiceApiContractsEmployeeV1ErrorResponse
    >({
      path: `/employees/${id}`,
      method: "PUT",
      body: employee,
      secure: true,
      type: ContentType.Json,
      ...params,
    });
  /**
   * @description Брисање запосленог по ID-ју
   *
   * @tags запослени
   * @name EmployeesDelete
   * @summary Брисање запосленог
   * @request DELETE:/employees/{id}
   * @secure
   */
  employeesDelete = (id: number, params: RequestParams = {}) =>
    this.request<void, GithubComPd120424DMountainServiceApiContractsEmployeeV1ErrorResponse>({
      path: `/employees/${id}`,
      method: "DELETE",
      secure: true,
      ...params,
    });
  /**
   * @description Проверава да ли запослени има активне хитне случајеве
   *
   * @tags запослени
   * @name ActiveEmergenciesList
   * @summary Провера активних хитних случајева за запосленог
   * @request GET:/employees/{id}/active-emergencies
   * @secure
   */
  activeEmergenciesList = (id: number, params: RequestParams = {}) =>
    this.request<
      GithubComPd120424DMountainServiceApiContractsEmployeeV1ActiveEmergenciesResponse,
      GithubComPd120424DMountainServiceApiContractsEmployeeV1ErrorResponse
    >({
      path: `/employees/${id}/active-emergencies`,
      method: "GET",
      secure: true,
      type: ContentType.Json,
      format: "json",
      ...params,
    });
  /**
   * @description Дохватање смена за запосленог по ID-ју
   *
   * @tags запослени
   * @name ShiftsList
   * @summary Дохватање смена за запосленог
   * @request GET:/employees/{id}/shifts
   * @secure
   */
  shiftsList = (id: number, params: RequestParams = {}) =>
    this.request<GithubComPd120424DMountainServiceApiContractsEmployeeV1ShiftResponse[], any>({
      path: `/employees/${id}/shifts`,
      method: "GET",
      secure: true,
      ...params,
    });
  /**
   * @description Додељује смену запосленом по ID-ју
   *
   * @tags запослени
   * @name ShiftsCreate
   * @summary Додељује смену запосленом
   * @request POST:/employees/{id}/shifts
   * @secure
   */
  shiftsCreate = (
    id: number,
    shift: GithubComPd120424DMountainServiceApiContractsEmployeeV1AssignShiftRequest,
    params: RequestParams = {},
  ) =>
    this.request<
      GithubComPd120424DMountainServiceApiContractsEmployeeV1AssignShiftResponse,
      GithubComPd120424DMountainServiceApiContractsEmployeeV1ErrorResponse
    >({
      path: `/employees/${id}/shifts`,
      method: "POST",
      body: shift,
      secure: true,
      type: ContentType.Json,
      ...params,
    });
  /**
   * @description Уклањање смене за запосленог по ID-ју и подацима о смени
   *
   * @tags запослени
   * @name ShiftsDelete
   * @summary Уклањање смене за запосленог
   * @request DELETE:/employees/{id}/shifts
   * @secure
   */
  shiftsDelete = (
    id: number,
    shift: GithubComPd120424DMountainServiceApiContractsEmployeeV1RemoveShiftRequest,
    params: RequestParams = {},
  ) =>
    this.request<void, GithubComPd120424DMountainServiceApiContractsEmployeeV1ErrorResponse>({
      path: `/employees/${id}/shifts`,
      method: "DELETE",
      body: shift,
      secure: true,
      type: ContentType.Json,
      ...params,
    });
}
