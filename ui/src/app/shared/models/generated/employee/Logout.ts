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

export class Logout<SecurityDataType = unknown> extends HttpClient<SecurityDataType>  {

            /**
 * @description Одјављивање запосленог и поништавање токена
 *
 * @tags запослени
 * @name LogoutCreate
 * @summary Одјављивање запосленог
 * @request POST:/logout
 * @secure
 */
logoutCreate: (params: RequestParams = {}) =>
    this.request<InternalHandlerMessageResponse, InternalHandlerErrorResponse>({
        path: `/logout`,
        method: 'POST',
                        secure: true,                format: "json",        ...params,
    }),    }
