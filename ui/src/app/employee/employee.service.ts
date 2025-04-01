// src/app/employee/employee.service.ts
import { Injectable } from '@angular/core';
import { HttpClient } from '@angular/common/http';
import { catchError, Observable, throwError } from 'rxjs';
import { Employee } from './employee.model';
import { LoggingService } from '../logging.service';
import { environment } from '../../environments/environment'; // Import environment variables


@Injectable({
  providedIn: 'root'
})
export class EmployeeService {
  // private apiUrl = `${environment.apiUrl}/employees`;
  private baseApiUrl = environment.useMockApi
    ? '/api/v1' // Mock server URL
    : `${environment.apiUrl}`; // Real API
  // private apiUrl = 'http://localhost:8082/api/v1/employees';
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

  addEmployee(employee: Employee): Observable<Employee> {
    return this.http.post<Employee>(this.employeeApiUrl, employee).pipe(
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

  private handleError(error: any): Observable<never> {
    console.error('An error occurred:', error);
    return throwError(() => new Error('Something went wrong; please try again later.'));
  }
}
