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
  GithubComPd120424DMountainServiceApiEmployeeInternalHandlerUploadProfilePictureResponse,
  InternalHandlerActiveEmergenciesResponse,
  InternalHandlerAssignShiftRequest,
  InternalHandlerAssignShiftResponse,
  InternalHandlerEmployeeResponse,
  InternalHandlerEmployeeUpdateRequest,
  InternalHandlerErrorResponse,
  InternalHandlerOnCallEmployeesResponse,
  InternalHandlerRemoveShiftRequest,
  InternalHandlerShiftResponse,
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
    this.request<InternalHandlerEmployeeResponse[][], any>({
      path: `/employees`,
      method: "GET",
      secure: true,
      format: "json",
      ...params,
    });
  /**
   * No description
   *
   * @name EmployeesCreate
   * @request POST:/employees
   * @secure
   */
  employeesCreate = (params: RequestParams = {}) =>
    this.request<InternalHandlerEmployeeResponse, InternalHandlerErrorResponse>({
      path: `/employees`,
      method: "POST",
      secure: true,
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
    this.request<InternalHandlerOnCallEmployeesResponse, InternalHandlerErrorResponse>({
      path: `/employees/on-call`,
      method: "GET",
      query: query,
      secure: true,
      type: ContentType.Json,
      format: "json",
      ...params,
    });
  /**
   * @description Преузимање запосленог по ID-ју
   *
   * @tags запослени
   * @name EmployeesDetail
   * @summary Преузимање запосленог по ID-ју
   * @request GET:/employees/{id}
   * @secure
   */
  employeesDetail = (id: number, params: RequestParams = {}) =>
    this.request<InternalHandlerEmployeeResponse, InternalHandlerErrorResponse>({
      path: `/employees/${id}`,
      method: "GET",
      secure: true,
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
  employeesUpdate = (id: number, employee: InternalHandlerEmployeeUpdateRequest, params: RequestParams = {}) =>
    this.request<InternalHandlerEmployeeResponse, InternalHandlerErrorResponse>({
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
    this.request<void, InternalHandlerErrorResponse>({
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
    this.request<InternalHandlerActiveEmergenciesResponse, InternalHandlerErrorResponse>({
      path: `/employees/${id}/active-emergencies`,
      method: "GET",
      secure: true,
      type: ContentType.Json,
      format: "json",
      ...params,
    });
  /**
   * @description Upload a profile picture for an employee
   *
   * @tags files
   * @name ProfilePictureCreate
   * @summary Upload profile picture
   * @request POST:/employees/{id}/profile-picture
   * @secure
   */
  profilePictureCreate = (
    id: number,
    data: {
      /** Profile picture file */
      file: File;
    },
    params: RequestParams = {},
  ) =>
    this.request<
      GithubComPd120424DMountainServiceApiEmployeeInternalHandlerUploadProfilePictureResponse,
      GithubComPd120424DMountainServiceApiEmployeeInternalHandlerErrorResponse
    >({
      path: `/employees/${id}/profile-picture`,
      method: "POST",
      body: data,
      secure: true,
      type: ContentType.FormData,
      format: "json",
      ...params,
    });
  /**
   * @description Delete a profile picture for an employee
   *
   * @tags files
   * @name ProfilePictureDelete
   * @summary Delete profile picture
   * @request DELETE:/employees/{id}/profile-picture
   * @secure
   */
  profilePictureDelete = (
    id: number,
    query: {
      /** Blob name to delete */
      blobName: string;
    },
    params: RequestParams = {},
  ) =>
    this.request<
      GithubComPd120424DMountainServiceApiEmployeeInternalHandlerMessageResponse,
      GithubComPd120424DMountainServiceApiEmployeeInternalHandlerErrorResponse
    >({
      path: `/employees/${id}/profile-picture`,
      method: "DELETE",
      query: query,
      secure: true,
      format: "json",
      ...params,
    });
  /**
   * @description Враћа листу упозорења о сменама за запосленог (нпр. недостају смене, није испуњена норма)
   *
   * @tags запослени
   * @name ShiftWarningsList
   * @summary Дохватање упозорења о сменама за запосленог
   * @request GET:/employees/{id}/shift-warnings
   * @secure
   */
  shiftWarningsList = (id: number, params: RequestParams = {}) =>
    this.request<Record<string, string[]>, InternalHandlerErrorResponse>({
      path: `/employees/${id}/shift-warnings`,
      method: "GET",
      secure: true,
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
    this.request<InternalHandlerShiftResponse[], any>({
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
  shiftsCreate = (id: number, shift: InternalHandlerAssignShiftRequest, params: RequestParams = {}) =>
    this.request<InternalHandlerAssignShiftResponse, InternalHandlerErrorResponse>({
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
  shiftsDelete = (id: number, shift: InternalHandlerRemoveShiftRequest, params: RequestParams = {}) =>
    this.request<void, InternalHandlerErrorResponse>({
      path: `/employees/${id}/shifts`,
      method: "DELETE",
      body: shift,
      secure: true,
      type: ContentType.Json,
      ...params,
    });
}
