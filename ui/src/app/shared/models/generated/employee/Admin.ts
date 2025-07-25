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
  GithubComPd120424DMountainServiceApiContractsEmployeeV1MessageResponse,
} from "./data-contracts";
import { HttpClient, RequestParams } from "./http-client";

export class Admin<SecurityDataType = unknown> extends HttpClient<SecurityDataType> {
  /**
   * @description Брише све запослене, смене и повезане податке из система (само за админе)
   *
   * @tags админ
   * @name ResetDelete
   * @summary Ресетовање свих података
   * @request DELETE:/admin/reset
   * @secure
   */
  resetDelete = (params: RequestParams = {}) =>
    this.request<
      GithubComPd120424DMountainServiceApiContractsEmployeeV1MessageResponse,
      GithubComPd120424DMountainServiceApiContractsEmployeeV1ErrorResponse
    >({
      path: `/admin/reset`,
      method: "DELETE",
      secure: true,
      format: "json",
      ...params,
    });
}
