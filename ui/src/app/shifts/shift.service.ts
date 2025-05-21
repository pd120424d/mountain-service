import { Injectable } from '@angular/core';
import { HttpClient } from '@angular/common/http';
import { Observable } from 'rxjs';
import { Employee } from '../employee/employee.model';
import { DatePipe } from '@angular/common';
import { ShiftAvailabilityResponse, AssignShiftRequest, AssignShiftResponse, RemoveShiftRequest } from './shift.model';

@Injectable({
  providedIn: 'root',
})
export class ShiftManagementService {
  constructor(private http: HttpClient) { }

  getShiftAvailability(days: number = 7): Observable<ShiftAvailabilityResponse> {
    return this.http.get<ShiftAvailabilityResponse>(
      `/api/v1/shifts/availability?days=${days}`
    );
  }

  getAllEmployees(): Observable<Employee[]> {
    return this.http.get<Employee[]>('/api/v1/employees');
  }

  assignEmployeeToShift(shiftType: number, employeeId: string, date: Date): Observable<AssignShiftResponse> {
    var req = <AssignShiftRequest>{ shiftDate: date.toISOString().split('T')[0], shiftType: shiftType };
    return this.http.post<AssignShiftResponse>(`/api/v1/employees/${employeeId}/shifts`, req);
  }

  removeEmployeeFromShift(employeeId: string, shiftId: number): Observable<any> {
    return this.http.request('delete', `/api/v1/employees/${employeeId}/shifts`, {
      body: <RemoveShiftRequest>{ id: shiftId }
    });
  }
}

