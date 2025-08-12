import { Injectable } from '@angular/core';
import { HttpClient, HttpErrorResponse } from '@angular/common/http';
import { Observable, throwError } from 'rxjs';
import { catchError } from 'rxjs/operators';
import { environment } from '../../environments/environment';
import { Activity, ActivityCreateRequest, ActivityListResponse } from '../shared/models';

@Injectable({
  providedIn: 'root',
})
export class ActivityService {
  private baseApiUrl = environment.useMockApi
    ? '/api/v1'
    : `${environment.apiUrl}`;
  private activityApiUrl = this.baseApiUrl + "/activities"

  constructor(private http: HttpClient) { }

  getActivities(): Observable<Activity[]> {
    return this.http.get<Activity[]>(this.activityApiUrl).pipe(
      catchError(this.handleError)
    );
  }

  getActivitiesByUrgency(urgencyId: number): Observable<Activity[]> {
    return this.http.get<Activity[]>(`${this.activityApiUrl}?targetId=${urgencyId}&targetType=urgency`).pipe(
      catchError(this.handleError)
    );
  }

  getActivitiesWithPagination(params: {
    targetId?: number;
    targetType?: string;
    page?: number;
    pageSize?: number;
  }): Observable<ActivityListResponse> {
    const queryParams = new URLSearchParams();
    
    if (params.targetId) queryParams.append('targetId', params.targetId.toString());
    if (params.targetType) queryParams.append('targetType', params.targetType);
    if (params.page) queryParams.append('page', params.page.toString());
    if (params.pageSize) queryParams.append('pageSize', params.pageSize.toString());

    const url = queryParams.toString() ? `${this.activityApiUrl}?${queryParams.toString()}` : this.activityApiUrl;
    
    return this.http.get<ActivityListResponse>(url).pipe(
      catchError(this.handleError)
    );
  }

  createActivity(activityRequest: ActivityCreateRequest): Observable<Activity> {
    return this.http.post<Activity>(this.activityApiUrl, activityRequest).pipe(
      catchError(this.handleError)
    );
  }

  private handleError(error: HttpErrorResponse): Observable<never> {
    let errorMessage = 'Something went wrong. Please try again later.';

    if (error.error instanceof ErrorEvent) {
      // Client-side error
      errorMessage = `Client error: ${error.error.message}`;
    } else {
      // Server-side error - only show specific messages for certain status codes
      if (error.status === 409) {
        errorMessage = error.error?.error || 'Conflict: Resource already exists';
      } else if (error.status === 400) {
        errorMessage = error.error?.error || 'Invalid data provided';
      }
    }

    return throwError(() => new Error(errorMessage));
  }
}
