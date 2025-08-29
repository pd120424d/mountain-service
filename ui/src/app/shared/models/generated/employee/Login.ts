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
  GithubComPd120424DMountainServiceApiEmployeeInternalHandlerEmployeeLogin,
  GithubComPd120424DMountainServiceApiEmployeeInternalHandlerErrorResponse,
  GithubComPd120424DMountainServiceApiEmployeeInternalHandlerTokenResponse,
} from "./data-contracts";
import { ContentType, HttpClient, RequestParams } from "./http-client";

export class Login<SecurityDataType = unknown> extends HttpClient<SecurityDataType> {
  /**
   * @description Пријавање запосленог са корисничким именом и лозинком
   *
   * @tags запослени
   * @name LoginCreate
   * @summary Пријавање запосленог
   * @request POST:/login
   * @secure
   */
  loginCreate = (
    employee: GithubComPd120424DMountainServiceApiEmployeeInternalHandlerEmployeeLogin,
    params: RequestParams = {},
  ) =>
    this.request<
      GithubComPd120424DMountainServiceApiEmployeeInternalHandlerTokenResponse,
      GithubComPd120424DMountainServiceApiEmployeeInternalHandlerErrorResponse
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
