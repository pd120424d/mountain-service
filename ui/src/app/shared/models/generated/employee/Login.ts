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
  V1EmployeeLogin,
  V1TokenResponse,
} from "./data-contracts";
import { ContentType, HttpClient, RequestParams } from "./http-client";

export class Login<
  SecurityDataType = unknown,
> extends HttpClient<SecurityDataType> {
  /**
   * @description Пријавање запосленог са корисничким именом и лозинком
   *
   * @tags запослени
   * @name LoginCreate
   * @summary Пријавање запосленог
   * @request POST:/login
   * @secure
   */
  loginCreate = (employee: V1EmployeeLogin, params: RequestParams = {}) =>
    this.request<
      V1TokenResponse,
      GithubComPd120424DMountainServiceApiContractsEmployeeV1ErrorResponse
    >({
      path: `/login`,
      method: "POST",
      body: employee,
      secure: true,
      type: ContentType.Json,
      format: "json",
      ...params,
    });
}
