// src/app/employee/employee.service.ts
import { Injectable } from '@angular/core';
import { HttpClient, HttpErrorResponse } from '@angular/common/http';
import { catchError, Observable, throwError } from 'rxjs';
import { Employee, EmployeeCreateRequest } from './employee.model';
import { LoggingService } from '../services/logging.service';
import { environment } from '../../environments/environment'; // Import environment variables


@Injectable({
  providedIn: 'root'
})
export class EmployeeService {
  private baseApiUrl = environment.useMockApi
    ? '/api/v1' // Mock server URL
    : `${environment.apiUrl}`; // Real API
  private employeeApiUrl = this.baseApiUrl + "/employees"

  constructor(
    private http: HttpClient,
    private logger: LoggingService
  ) {
    this.logger.info(`Starting employee service with url: ${this.employeeApiUrl}`);
    this.logger.info(`Starting employee service with base apiUrl: ${this.baseApiUrl}`);
   }

  // Example methods
  getEmployees(): Observable<Employee[]> {
    return this.http.get<Employee[]>(this.employeeApiUrl).pipe(
      catchError(this.handleError)
    );;
  }

  getEmployeeById(id: number): Observable<Employee> {
    this.logger.info(`Fetching employee with ID: ${id}`);
    return this.http.get<Employee>(`${this.employeeApiUrl}/${id}`).pipe(
      catchError(this.handleError)
    );
  }

  addEmployee(employeeCreateRequest: EmployeeCreateRequest): Observable<Employee> {
    return this.http.post<Employee>(this.employeeApiUrl, employeeCreateRequest).pipe(
      catchError(this.handleError)
    );
  }

  updateEmployee(id: number, employee: Employee): Observable<Employee> {
    return this.http.put<Employee>(`${this.employeeApiUrl}/${id}`, employee).pipe(
      catchError(this.handleError)
    );
  }

  deleteEmployee(id: number): Observable<void> {
    return this.http.delete<void>(`${this.employeeApiUrl}/${id}`).pipe(
      catchError(this.handleError)
    );
  }

  private handleError(error: HttpErrorResponse): Observable<never> {
    console.error('An error occurred:', error);

    let errorMessage = 'Something went wrong; please try again later.';

    if (error.error instanceof ErrorEvent) {
      // Client-side error
      errorMessage = `Client error: ${error.error.message}`;
    } else {
      // Server-side error
      if (error.status === 409) {
        errorMessage = error.error?.error || 'Conflict: Resource already exists';
      } else if (error.status === 400) {
        errorMessage = error.error?.error || 'Invalid data provided';
      } else {
        errorMessage = `Server error: ${error.status} - ${error.message}`;
      }
    }

    this.logger.error(`Employee service error: ${errorMessage}`);
    return throwError(() => new Error(errorMessage));
  }
}
