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

export class Files<
  SecurityDataType = unknown,
> extends HttpClient<SecurityDataType> {
  /**
   * @description Get information about a profile picture
   *
   * @tags files
   * @name ProfilePictureInfoList
   * @summary Get profile picture info
   * @request GET:/files/profile-picture/info
   * @secure
   */
  profilePictureInfoList = (
    query: {
      /** Blob name */
      blobName: string;
    },
    params: RequestParams = {},
  ) =>
    this.request<Record<string, any>, Record<string, any>>({
      path: `/files/profile-picture/info`,
      method: "GET",
      query: query,
      secure: true,
      format: "json",
      ...params,
    });
}
