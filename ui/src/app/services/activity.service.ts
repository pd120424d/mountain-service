import { Injectable } from '@angular/core';
import { HttpClient, HttpErrorResponse } from '@angular/common/http';
import { Observable, throwError } from 'rxjs';
import { catchError, map } from 'rxjs/operators';
import { environment } from '../../environments/environment';
import { Activity, ActivityCreatePayload } from '../shared/models';

@Injectable({ providedIn: 'root' })
export class ActivityService {
  private baseApiUrl = environment.useMockApi ? '/api/v1' : `${environment.apiUrl}`;
  private activityApiUrl = this.baseApiUrl + '/activities';

  constructor(private http: HttpClient) {}

  private mapActivityFromServer(item: any): Activity {
    if (!item || typeof item !== 'object') return item as Activity;
    return {
      id: item.id,
      description: item.description,
      employeeId: item.employeeId,
      urgencyId: item.urgencyId,
      createdAt: item.createdAt,
      updatedAt: item.updatedAt,
    } as Activity;
  }

  getActivities(): Observable<Activity[]> {
    return this.http.get<any>(this.activityApiUrl).pipe(
      map((resp) => (resp?.activities ?? []).map((a: any) => this.mapActivityFromServer(a))),
      catchError(this.handleError)
    );
  }

  getActivitiesByUrgency(urgencyId: number): Observable<Activity[]> {
    const url = `${this.activityApiUrl}?urgencyId=${urgencyId}`;
    return this.http.get<any>(url).pipe(
      map((resp) => (resp?.activities ?? []).map((a: any) => this.mapActivityFromServer(a))),
      catchError(this.handleError)
    );
  }

  getActivitiesWithPagination(params: {
    urgencyId?: number;
    employeeId?: number;
    page?: number;
    pageSize?: number;
  }): Observable<{ activities: Activity[]; total: number; page: number; pageSize: number; totalPages: number }> {
    const queryParams = new URLSearchParams();
    if (params.urgencyId) queryParams.append('urgencyId', params.urgencyId.toString());
    if (params.employeeId) queryParams.append('employeeId', params.employeeId.toString());
    if (params.page) queryParams.append('page', params.page.toString());
    if (params.pageSize) queryParams.append('pageSize', params.pageSize.toString());

    const url = queryParams.toString() ? `${this.activityApiUrl}?${queryParams.toString()}` : this.activityApiUrl;

    return this.http.get<any>(url).pipe(
      map((resp) => ({
        activities: (resp?.activities ?? []).map((a: any) => this.mapActivityFromServer(a)),
        total: resp?.total ?? 0,
        page: resp?.page ?? 1,
        pageSize: resp?.pageSize ?? 20,
        totalPages: resp?.totalPages ?? 0,
      })),
      catchError(this.handleError)
    );
  }

  createActivity(activityRequest: ActivityCreatePayload): Observable<Activity> {
    const payload = {
      description: activityRequest.description,
      employeeId: activityRequest.employeeId,
      urgencyId: activityRequest.urgencyId,
    };
    return this.http.post<any>(this.activityApiUrl, payload).pipe(
      map((resp) => this.mapActivityFromServer(resp)),
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
