import { Injectable } from '@angular/core';
import { HttpClient, HttpErrorResponse } from '@angular/common/http';
import { Observable, throwError } from 'rxjs';
import { catchError } from 'rxjs/operators';
import { environment } from '../../environments/environment';
import { ActivityResponse, ActivityCreateRequest, ActivityListResponse } from '../shared/models';

@Injectable({
  providedIn: 'root',
})
export class ActivityService {
  private baseApiUrl = environment.useMockApi
    ? '/api/v1'
    : `${environment.apiUrl}`;
  private activityApiUrl = this.baseApiUrl + "/activities"

  constructor(private http: HttpClient) { }

  getActivities(): Observable<ActivityResponse[]> {
    return this.http.get<ActivityResponse[]>(this.activityApiUrl).pipe(
      catchError(this.handleError)
    );
  }

  getActivitiesByUrgency(urgencyId: number): Observable<ActivityResponse[]> {
    return this.http.get<ActivityResponse[]>(`${this.activityApiUrl}?urgency_id=${urgencyId}`).pipe(
      catchError(this.handleError)
    );
  }

  getActivitiesWithPagination(params: {
    urgency_id?: number;
    employee_id?: number;
    page?: number;
    page_size?: number;
  }): Observable<ActivityListResponse> {
    const queryParams = new URLSearchParams();

    if (params.urgency_id) queryParams.append('urgency_id', params.urgency_id.toString());
    if (params.employee_id) queryParams.append('employee_id', params.employee_id.toString());
    if (params.page) queryParams.append('page', params.page.toString());
    if (params.page_size) queryParams.append('page_size', params.page_size.toString());

    const url = queryParams.toString() ? `${this.activityApiUrl}?${queryParams.toString()}` : this.activityApiUrl;

    return this.http.get<ActivityListResponse>(url).pipe(
      catchError(this.handleError)
    );
  }

  createActivity(activityRequest: ActivityCreateRequest): Observable<ActivityResponse> {
    return this.http.post<ActivityResponse>(this.activityApiUrl, activityRequest).pipe(
      catchError(this.handleError)
    );
  }

  private handleError(error: HttpErrorResponse): Observable<never> {
    if (error && error.error && typeof error.error === 'object' && 'error' in error.error) {
      const payload = error.error as { error?: string } & Record<string, any>;
      if (typeof payload.error === 'string' && payload.error) {
        return throwError(() => new Error(payload.error));
      }
    }

    // Fallback to the old behavior
    let errorMessage = 'Something went wrong. Please try again later.';
    if (error.error instanceof ErrorEvent) {
      errorMessage = `Client error: ${error.error.message}`;
    } else {
      if (error.status === 409) {
        errorMessage = (error.error as any)?.error || 'Conflict: Resource already exists';
      } else if (error.status === 400) {
        errorMessage = (error.error as any)?.error || 'Invalid data provided';
      }
    }

    return throwError(() => new Error(errorMessage));
  }
}
