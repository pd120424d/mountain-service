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
import { GithubComPd120424DMountainServiceApiEmployeeInternalHandlerActiveEmergenciesResponse, GithubComPd120424DMountainServiceApiEmployeeInternalHandlerAssignShiftRequest, GithubComPd120424DMountainServiceApiEmployeeInternalHandlerAssignShiftResponse, GithubComPd120424DMountainServiceApiEmployeeInternalHandlerEmployeeLogin, GithubComPd120424DMountainServiceApiEmployeeInternalHandlerEmployeeResponse, GithubComPd120424DMountainServiceApiEmployeeInternalHandlerEmployeeUpdateRequest, GithubComPd120424DMountainServiceApiEmployeeInternalHandlerErrorResponse, GithubComPd120424DMountainServiceApiEmployeeInternalHandlerMessageResponse, GithubComPd120424DMountainServiceApiEmployeeInternalHandlerOnCallEmployeesResponse, GithubComPd120424DMountainServiceApiEmployeeInternalHandlerRemoveShiftRequest, GithubComPd120424DMountainServiceApiEmployeeInternalHandlerShiftAvailabilityResponse, GithubComPd120424DMountainServiceApiEmployeeInternalHandlerShiftResponse, GithubComPd120424DMountainServiceApiEmployeeInternalHandlerTokenResponse, GithubComPd120424DMountainServiceApiEmployeeInternalHandlerUploadProfilePictureResponse, InternalHandlerActiveEmergenciesResponse, InternalHandlerAssignShiftRequest, InternalHandlerAssignShiftResponse, InternalHandlerEmployeeLogin, InternalHandlerEmployeeResponse, InternalHandlerEmployeeUpdateRequest, InternalHandlerErrorResponse, InternalHandlerMessageResponse, InternalHandlerOnCallEmployeesResponse, InternalHandlerRemoveShiftRequest, InternalHandlerShiftAvailabilityResponse, InternalHandlerShiftResponse, InternalHandlerTokenResponse, InternalHandlerUploadProfilePictureResponse, V1EmployeeResponse, V1ShiftAvailability, V1ShiftAvailabilityPerDay } from "./data-contracts"

export class Admin<SecurityDataType = unknown> extends HttpClient<SecurityDataType>  {

            /**
 * @description Брише све запослене, смене и повезане податке из система (само за админе)
 *
 * @tags админ
 * @name ResetDelete
 * @summary Ресетовање свих података
 * @request DELETE:/admin/reset
 * @secure
 */
resetDelete: (params: RequestParams = {}) =>
    this.request<InternalHandlerMessageResponse, InternalHandlerErrorResponse>({
        path: `/admin/reset`,
        method: 'DELETE',
                        secure: true,                format: "json",        ...params,
    }),            /**
 * @description Дохватање доступности смена за све запослене (само за админе)
 *
 * @tags админ
 * @name ShiftsAvailabilityList
 * @summary Дохватање доступности смена за админе
 * @request GET:/admin/shifts/availability
 * @secure
 */
shiftsAvailabilityList: (query?: {
  /** Број дана за које се проверава доступност (подразумевано 7) */
    days?: number,

}, params: RequestParams = {}) =>
    this.request<InternalHandlerShiftAvailabilityResponse, InternalHandlerErrorResponse>({
        path: `/admin/shifts/availability`,
        method: 'GET',
        query: query,                secure: true,                        ...params,
    }),    }
