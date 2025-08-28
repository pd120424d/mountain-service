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

import { InternalHandlerErrorResponse, InternalHandlerMessageResponse } from "./data-contracts";
import { HttpClient, RequestParams } from "./http-client";

export class Logout<SecurityDataType = unknown> extends HttpClient<SecurityDataType> {
  /**
   * @description Одјављивање запосленог и поништавање токена
   *
   * @tags запослени
   * @name LogoutCreate
   * @summary Одјављивање запосленог
   * @request POST:/logout
   * @secure
   */
  logoutCreate = (params: RequestParams = {}) =>
    this.request<InternalHandlerMessageResponse, InternalHandlerErrorResponse>({
      path: `/logout`,
      method: "POST",
      secure: true,
      format: "json",
      ...params,
    });
}
