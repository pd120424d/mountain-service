
import { Injectable } from '@angular/core';
import { HttpClient, HttpErrorResponse } from '@angular/common/http';
import { Observable, throwError } from 'rxjs';
import { catchError, tap } from 'rxjs/operators';
import { environment } from '../../environments/environment';
import { Urgency, UrgencyCreateRequest, UrgencyUpdateRequest } from '../shared/models';

@Injectable({
  providedIn: 'root',
})
export class UrgencyService {
  private baseApiUrl = environment.useMockApi
    ? '/api/v1'
    : `${environment.apiUrl}`;
  private urgencyApiUrl = this.baseApiUrl + "/urgencies"

  constructor(private http: HttpClient) { }

  getUrgencies(): Observable<Urgency[]> {
    return this.http.get<Urgency[]>(this.urgencyApiUrl).pipe(
      catchError(this.handleError)
    );
  }

  getUrgencyById(id: number): Observable<Urgency> {
    return this.http.get<Urgency>(`${this.urgencyApiUrl}/${id}`).pipe(
      catchError(this.handleError)
    );
  }

  addUrgency(urgencyRequest: UrgencyCreateRequest): Observable<Urgency> {
    return this.http.post<Urgency>(this.urgencyApiUrl, urgencyRequest).pipe(
      catchError(this.handleError)
    );
  }

  updateUrgency(id: number, urgencyRequest: UrgencyUpdateRequest): Observable<Urgency> {
    return this.http.put<Urgency>(`${this.urgencyApiUrl}/${id}`, urgencyRequest).pipe(
      catchError(this.handleError)
    );
  }

  deleteUrgency(id: number): Observable<void> {
    return this.http.delete<void>(`${this.urgencyApiUrl}/${id}`).pipe(
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




