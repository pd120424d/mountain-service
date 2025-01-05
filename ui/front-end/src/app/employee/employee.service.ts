// src/app/employee/employee.service.ts
import { Injectable } from '@angular/core';
import { HttpClient } from '@angular/common/http';
import { catchError, Observable, throwError } from 'rxjs';
import { Employee } from './employee.model';
import { environment } from '../../environments/environment'; // Import environment variables


@Injectable({
  providedIn: 'root'
})
export class EmployeeService {
  // private apiUrl = `${environment.apiUrl}/employees`;
  private apiUrl = `${environment.useMockApi}`
    ? '/api/v1/employees' // Mock server URL
    : `${environment.apiUrl}`; // Real API
  // private apiUrl = 'http://localhost:8082/api/v1/employees'; // Replace with your actual API URL

  constructor(private http: HttpClient) { }  // Inject HttpClient

  // Example methods
  getEmployees(): Observable<Employee[]> {
    return this.http.get<Employee[]>(this.apiUrl).pipe(
      catchError(this.handleError)
    );;
  }

  getEmployeeById(id: number): Observable<Employee> {
    return this.http.get<Employee>(`${this.apiUrl}/${id}`).pipe(
      catchError(this.handleError)
    );
  }

  addEmployee(employee: Employee): Observable<Employee> {
    return this.http.post<Employee>(this.apiUrl, employee).pipe(
      catchError(this.handleError)
    );
  }

  updateEmployee(id: number, employee: Employee): Observable<Employee> {
    return this.http.put<Employee>(`${this.apiUrl}/${id}`, employee).pipe(
      catchError(this.handleError)
    );
  }

  deleteEmployee(id: number): Observable<void> {
    return this.http.delete<void>(`${this.apiUrl}/${id}`).pipe(
      catchError(this.handleError)
    );
  }

  private handleError(error: any): Observable<never> {
    console.error('An error occurred:', error);
    return throwError(() => new Error('Something went wrong; please try again later.'));
  }
}
