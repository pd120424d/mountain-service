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
  GithubComPd120424DMountainServiceApiContractsEmployeeV1ErrorResponse,
  InternalHandlerUploadProfilePictureResponse,
  V1ActiveEmergenciesResponse,
  V1AssignShiftRequest,
  V1AssignShiftResponse,
  V1EmployeeResponse,
  V1EmployeeUpdateRequest,
  V1OnCallEmployeesResponse,
  V1RemoveShiftRequest,
  V1ShiftResponse,
} from "./data-contracts";
import { ContentType, HttpClient, RequestParams } from "./http-client";

export class Employees<
  SecurityDataType = unknown,
> extends HttpClient<SecurityDataType> {
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
    this.request<V1EmployeeResponse[][], any>({
      path: `/employees`,
      method: "GET",
      secure: true,
      format: "json",
      ...params,
    });
  /**
   * @description Креирање новог запосленог у систему (supports both JSON and multipart form data)
   *
   * @tags запослени
   * @name EmployeesCreate
   * @summary Креирање новог запосленог
   * @request POST:/employees
   * @secure
   */
  employeesCreate = (
    employee: {
      /** First Name (form data) */
      firstName?: string;
      /** Last Name (form data) */
      lastName?: string;
      /** Username (form data) */
      username?: string;
      /** Password (form data) */
      password?: string;
      /** Email (form data) */
      email?: string;
      /** Gender (form data) */
      gender?: string;
      /** Phone (form data) */
      phone?: string;
      /** Profile Type (form data) */
      profileType?: string;
      /**
       * Profile Picture (form data)
       * @format binary
       */
      profilePicture?: File;
    },
    params: RequestParams = {},
  ) =>
    this.request<
      V1EmployeeResponse,
      GithubComPd120424DMountainServiceApiContractsEmployeeV1ErrorResponse
    >({
      path: `/employees`,
      method: "POST",
      body: employee,
      secure: true,
      type: ContentType.FormData,
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
      V1OnCallEmployeesResponse,
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
   * @description Upload a profile picture for an employee
   *
   * @tags files
   * @name ProfilePictureCreate
   * @summary Upload profile picture
   * @request POST:/employees/{employeeId}/profile-picture
   * @secure
   */
  profilePictureCreate = (
    employeeId: number,
    data: {
      /** Profile picture file */
      file: File;
    },
    params: RequestParams = {},
  ) =>
    this.request<
      InternalHandlerUploadProfilePictureResponse,
      Record<string, any>
    >({
      path: `/employees/${employeeId}/profile-picture`,
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
   * @request DELETE:/employees/{employeeId}/profile-picture
   * @secure
   */
  profilePictureDelete = (
    employeeId: number,
    query: {
      /** Blob name to delete */
      blobName: string;
    },
    params: RequestParams = {},
  ) =>
    this.request<Record<string, any>, Record<string, any>>({
      path: `/employees/${employeeId}/profile-picture`,
      method: "DELETE",
      query: query,
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
  employeesUpdate = (
    id: number,
    employee: V1EmployeeUpdateRequest,
    params: RequestParams = {},
  ) =>
    this.request<
      V1EmployeeResponse,
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
    this.request<
      void,
      GithubComPd120424DMountainServiceApiContractsEmployeeV1ErrorResponse
    >({
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
      V1ActiveEmergenciesResponse,
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
   * @description Враћа листу упозорења о сменама за запосленог (нпр. недостају смене, није испуњена норма)
   *
   * @tags запослени
   * @name ShiftWarningsList
   * @summary Дохватање упозорења о сменама за запосленог
   * @request GET:/employees/{id}/shift-warnings
   * @secure
   */
  shiftWarningsList = (id: number, params: RequestParams = {}) =>
    this.request<
      Record<string, string[]>,
      GithubComPd120424DMountainServiceApiContractsEmployeeV1ErrorResponse
    >({
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
    this.request<V1ShiftResponse[], any>({
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
    shift: V1AssignShiftRequest,
    params: RequestParams = {},
  ) =>
    this.request<
      V1AssignShiftResponse,
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
    shift: V1RemoveShiftRequest,
    params: RequestParams = {},
  ) =>
    this.request<
      void,
      GithubComPd120424DMountainServiceApiContractsEmployeeV1ErrorResponse
    >({
      path: `/employees/${id}/shifts`,
      method: "DELETE",
      body: shift,
      secure: true,
      type: ContentType.Json,
      ...params,
    });
}
