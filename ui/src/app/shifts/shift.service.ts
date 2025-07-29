import { Injectable } from '@angular/core';
import { HttpClient, HttpErrorResponse } from '@angular/common/http';
import { Observable, throwError } from 'rxjs';
import { catchError, tap } from 'rxjs/operators';
import { Employee, ShiftAvailabilityResponse, AssignShiftRequest, AssignShiftResponse, RemoveShiftRequest, ShiftResponse } from '../shared/models';

@Injectable({
  providedIn: 'root',
})
export class ShiftManagementService {
  constructor(private http: HttpClient) { }

  getShiftAvailability(days: number = 7): Observable<ShiftAvailabilityResponse> {
    console.log(`Fetching shift availability for ${days} days`);
    return this.http.get<ShiftAvailabilityResponse>(
      `/api/v1/shifts/availability?days=${days}`
    ).pipe(
      tap(response => console.log('Shift availability response:', response)),
      catchError(this.handleError)
    );
  }

  getAllEmployees(): Observable<Employee[]> {
    console.log('Fetching all employees');
    return this.http.get<Employee[]>('/api/v1/employees').pipe(
      tap(response => console.log('Employees response:', response)),
      catchError(this.handleError)
    );
  }

  assignEmployeeToShift(shiftType: number, employeeId: string, date: Date): Observable<AssignShiftResponse> {
    const req = <AssignShiftRequest>{ shiftDate: date.toISOString().split('T')[0], shiftType: shiftType };
    console.log('Assigning employee to shift:', { employeeId, req });
    return this.http.post<AssignShiftResponse>(`/api/v1/employees/${employeeId}/shifts`, req).pipe(
      tap(response => console.log('Assignment response:', response)),
      catchError(this.handleError)
    );
  }
  removeEmployeeFromShiftByDetails(employeeId: string, shiftType: number, date: Date): Observable<any> {
    const req: RemoveShiftRequest = {
      shiftDate: date.toISOString().split('T')[0],
      shiftType: shiftType
    };
    console.log('Removing employee from shift by details:', { employeeId, req });
    return this.http.request('delete', `/api/v1/employees/${employeeId}/shifts`, {
      body: req
    }).pipe(
      tap(response => console.log('Removal by details response:', response)),
      catchError(this.handleError)
    );
  }

  getShiftWarnings(employeeId: string): Observable<{warnings: string[]}> {
    console.log(`Fetching shift warnings for employee ${employeeId}`);
    return this.http.get<{warnings: string[]}>(`/api/v1/employees/${employeeId}/shift-warnings`).pipe(
      tap(response => console.log('Shift warnings response:', response)),
      catchError(this.handleError)
    );
  }

  getEmployeeShifts(employeeId: string): Observable<ShiftResponse[]> {
    console.log(`Fetching shifts for employee ${employeeId}`);
    return this.http.get<ShiftResponse[]>(`/api/v1/employees/${employeeId}/shifts`).pipe(
      tap(response => console.log('Employee shifts response:', response)),
      catchError(this.handleError)
    );
  }

  private handleError(error: HttpErrorResponse): Observable<never> {
    console.error('ShiftManagementService error:', error);
    let errorMessage = 'An unknown error occurred';

    if (error.error instanceof ErrorEvent) {
      // Client-side error
      errorMessage = `Client error: ${error.error.message}`;
    } else {
      // Server-side error
      errorMessage = `Server error: ${error.status} - ${error.message}`;
      if (error.error && typeof error.error === 'string') {
        errorMessage += ` - ${error.error}`;
      } else if (error.error && error.error.message) {
        errorMessage += ` - ${error.error.message}`;
      }
    }

    console.error('Processed error message:', errorMessage);
    return throwError(() => new Error(errorMessage));
  }
}

