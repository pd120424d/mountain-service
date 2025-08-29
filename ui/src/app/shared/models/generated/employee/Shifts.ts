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

export class Shifts<SecurityDataType = unknown> extends HttpClient<SecurityDataType>  {

            /**
 * @description Дохватање доступности смена за одређени дан
 *
 * @tags запослени
 * @name AvailabilityList
 * @summary Дохватање доступности смена
 * @request GET:/shifts/availability
 * @secure
 */
availabilityList: (query?: {
  /** Дан за који се проверава доступност смена */
    date?: string,

}, params: RequestParams = {}) =>
    this.request<InternalHandlerShiftAvailabilityResponse, InternalHandlerErrorResponse>({
        path: `/shifts/availability`,
        method: 'GET',
        query: query,                secure: true,                        ...params,
    }),    }
