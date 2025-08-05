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

import { GithubComPd120424DMountainServiceApiContractsEmployeeV1ErrorResponse } from "./data-contracts";
import { ContentType, HttpClient, RequestParams } from "./http-client";

export class Oauth<
  SecurityDataType = unknown,
> extends HttpClient<SecurityDataType> {
  /**
   * @description OAuth2 password flow token endpoint for Swagger UI authentication
   *
   * @tags authentication
   * @name TokenCreate
   * @summary OAuth2 token endpoint
   * @request POST:/oauth/token
   * @secure
   */
  tokenCreate = (
    data: {
      /** Username */
      username: string;
      /** Password */
      password: string;
    },
    params: RequestParams = {},
  ) =>
    this.request<
      Record<string, any>,
      GithubComPd120424DMountainServiceApiContractsEmployeeV1ErrorResponse
    >({
      path: `/oauth/token`,
      method: "POST",
      body: data,
      secure: true,
      type: ContentType.UrlEncoded,
      format: "json",
      ...params,
    });
}
