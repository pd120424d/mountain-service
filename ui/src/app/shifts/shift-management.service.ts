import { Injectable } from '@angular/core';
import { HttpClient } from '@angular/common/http';
import { Observable } from 'rxjs';
import { Employee } from '../employee/employee.model';

@Injectable({ providedIn: 'root' })
export class ShiftManagementService {
  constructor(private http: HttpClient) {}

  getShiftAvailability(): Observable<any> {
    return this.http.get('/api/v1/shifts/availability');
  }

  getAllEmployees(): Observable<Employee[]> {
    return this.http.get<Employee[]>('/api/v1/employees');
  }

  assignEmployeeToShift(shiftType: number, employeeId: string): Observable<any> {
    return this.http.post('/api/v1/shifts/assign', {
      shiftType,
      employeeId
    });
  }

  removeEmployeeFromShift(shiftType: number, employeeId: string): Observable<any> {
    return this.http.request('delete', '/api/v1/shifts/remove', {
      body: { shiftType, employeeId }
    });
  }
}

