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
  V1AssignmentCreateRequest,
  V1UrgencyCreateRequest,
  V1UrgencyResponse,
  V1UrgencyUpdateRequest,
} from "./data-contracts";
import { ContentType, HttpClient, RequestParams } from "./http-client";

export class Urgencies<SecurityDataType = unknown> extends HttpClient<SecurityDataType> {
  /**
   * @description Извлачење свих ургентних ситуација
   *
   * @tags urgency
   * @name UrgenciesList
   * @summary Извлачење листе ургентних ситуација
   * @request GET:/urgencies
   * @secure
   */
  urgenciesList = (params: RequestParams = {}) =>
    this.request<V1UrgencyResponse[][], any>({
      path: `/urgencies`,
      method: "GET",
      secure: true,
      format: "json",
      ...params,
    });
  /**
   * @description Креирање нове ургентне ситуације са свим потребним подацима
   *
   * @tags urgency
   * @name UrgenciesCreate
   * @summary Креирање нове ургентне ситуације
   * @request POST:/urgencies
   * @secure
   */
  urgenciesCreate = (urgency: V1UrgencyCreateRequest, params: RequestParams = {}) =>
    this.request<V1UrgencyResponse, any>({
      path: `/urgencies`,
      method: "POST",
      body: urgency,
      secure: true,
      type: ContentType.Json,
      format: "json",
      ...params,
    });
  /**
   * @description Извлачење ургентне ситуације по њеном ID
   *
   * @tags urgency
   * @name UrgenciesDetail
   * @summary Извлачење ургентне ситуације по ID
   * @request GET:/urgencies/{id}
   * @secure
   */
  urgenciesDetail = (id: number, params: RequestParams = {}) =>
    this.request<V1UrgencyResponse, any>({
      path: `/urgencies/${id}`,
      method: "GET",
      secure: true,
      format: "json",
      ...params,
    });
  /**
   * @description Ажурирање постојеће ургентне ситуације
   *
   * @tags urgency
   * @name UrgenciesUpdate
   * @summary Ажурирање ургентне ситуације
   * @request PUT:/urgencies/{id}
   * @secure
   */
  urgenciesUpdate = (id: number, urgency: V1UrgencyUpdateRequest, params: RequestParams = {}) =>
    this.request<V1UrgencyResponse, any>({
      path: `/urgencies/${id}`,
      method: "PUT",
      body: urgency,
      secure: true,
      type: ContentType.Json,
      format: "json",
      ...params,
    });
  /**
   * @description Брисање ургентне ситуације по ID
   *
   * @tags urgency
   * @name UrgenciesDelete
   * @summary Брисање ургентне ситуације
   * @request DELETE:/urgencies/{id}
   * @secure
   */
  urgenciesDelete = (id: number, params: RequestParams = {}) =>
    this.request<void, any>({
      path: `/urgencies/${id}`,
      method: "DELETE",
      secure: true,
      ...params,
    });
  /**
   * No description
   *
   * @tags urgency
   * @name AssignCreate
   * @summary Assign urgency to an employee
   * @request POST:/urgencies/{id}/assign
   * @secure
   */
  assignCreate = (id: number, payload: V1AssignmentCreateRequest, params: RequestParams = {}) =>
    this.request<Record<string, string>, Record<string, any>>({
      path: `/urgencies/${id}/assign`,
      method: "POST",
      body: payload,
      secure: true,
      type: ContentType.Json,
      ...params,
    });
  /**
   * No description
   *
   * @tags urgency
   * @name AssignDelete
   * @summary Unassign urgency
   * @request DELETE:/urgencies/{id}/assign
   * @secure
   */
  assignDelete = (id: number, params: RequestParams = {}) =>
    this.request<void, Record<string, any>>({
      path: `/urgencies/${id}/assign`,
      method: "DELETE",
      secure: true,
      ...params,
    });
}
