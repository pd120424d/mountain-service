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

import { HttpClient, RequestParams } from "./http-client";

export class Admin<SecurityDataType = unknown> extends HttpClient<SecurityDataType> {
  /**
   * @description Брисање свих ургентних ситуација (само за администраторе)
   *
   * @tags urgency
   * @name UrgenciesResetDelete
   * @summary Ресетовање свих података
   * @request DELETE:/admin/urgencies/reset
   * @secure
   */
  urgenciesResetDelete = (params: RequestParams = {}) =>
    this.request<void, any>({
      path: `/admin/urgencies/reset`,
      method: "DELETE",
      secure: true,
      ...params,
    });
}
