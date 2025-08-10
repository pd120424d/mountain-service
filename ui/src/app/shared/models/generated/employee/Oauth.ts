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
import { GithubComPd120424DMountainServiceApiEmployeeInternalHandlerActiveEmergenciesResponse, GithubComPd120424DMountainServiceApiEmployeeInternalHandlerAssignShiftRequest, GithubComPd120424DMountainServiceApiEmployeeInternalHandlerAssignShiftResponse, GithubComPd120424DMountainServiceApiEmployeeInternalHandlerEmployeeCreateRequest, GithubComPd120424DMountainServiceApiEmployeeInternalHandlerEmployeeLogin, GithubComPd120424DMountainServiceApiEmployeeInternalHandlerEmployeeResponse, GithubComPd120424DMountainServiceApiEmployeeInternalHandlerEmployeeUpdateRequest, GithubComPd120424DMountainServiceApiEmployeeInternalHandlerErrorResponse, GithubComPd120424DMountainServiceApiEmployeeInternalHandlerMessageResponse, GithubComPd120424DMountainServiceApiEmployeeInternalHandlerOnCallEmployeesResponse, GithubComPd120424DMountainServiceApiEmployeeInternalHandlerRemoveShiftRequest, GithubComPd120424DMountainServiceApiEmployeeInternalHandlerShiftAvailabilityResponse, GithubComPd120424DMountainServiceApiEmployeeInternalHandlerShiftResponse, GithubComPd120424DMountainServiceApiEmployeeInternalHandlerTokenResponse, GithubComPd120424DMountainServiceApiEmployeeInternalHandlerUploadProfilePictureResponse, InternalHandlerActiveEmergenciesResponse, InternalHandlerAssignShiftRequest, InternalHandlerAssignShiftResponse, InternalHandlerEmployeeCreateRequest, InternalHandlerEmployeeLogin, InternalHandlerEmployeeResponse, InternalHandlerEmployeeUpdateRequest, InternalHandlerErrorResponse, InternalHandlerMessageResponse, InternalHandlerOnCallEmployeesResponse, InternalHandlerRemoveShiftRequest, InternalHandlerShiftAvailabilityResponse, InternalHandlerShiftResponse, InternalHandlerTokenResponse, InternalHandlerUploadProfilePictureResponse, V1EmployeeResponse, V1ShiftAvailability, V1ShiftAvailabilityPerDay } from "./data-contracts"

export class Oauth<SecurityDataType = unknown> extends HttpClient<SecurityDataType>  {

            /**
 * @description OAuth2 password flow token endpoint for Swagger UI authentication
 *
 * @tags authentication
 * @name TokenCreate
 * @summary OAuth2 token endpoint
 * @request POST:/oauth/token
 * @secure
 */
tokenCreate: (data: {
  /** Username */
    username: string,
  /** Password */
    password: string,

}, params: RequestParams = {}) =>
    this.request<Record<string,any>, InternalHandlerErrorResponse>({
        path: `/oauth/token`,
        method: 'POST',
                body: data,        secure: true,        type: ContentType.UrlEncoded,        format: "json",        ...params,
    }),    }
