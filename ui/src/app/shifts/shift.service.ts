import { Injectable } from '@angular/core';
import { HttpClient } from '@angular/common/http';
import { Observable } from 'rxjs';
import { AssignShiftRequest, AssignShiftResponse, RemoveShiftRequest, ShiftResponse } from './shift.model';

@Injectable({
  providedIn: 'root',
})
export class ShiftService {
  private readonly baseUrl = '/api/v1';

  constructor(private http: HttpClient) {}

  getShifts(employeeId: number): Observable<ShiftResponse[]> {
    return this.http.get<ShiftResponse[]>(`${this.baseUrl}/employees/${employeeId}/shifts`);
  }

  assignShift(employeeId: number, request: AssignShiftRequest): Observable<AssignShiftResponse> {
    return this.http.post<AssignShiftResponse>(`${this.baseUrl}/employees/${employeeId}/shifts`, request);
  }

  removeShift(employeeId: number, request: RemoveShiftRequest): Observable<void> {
    return this.http.request<void>('delete', `${this.baseUrl}/employees/${employeeId}/shifts`, {
      body: request,
    });
  }

  getShiftAvailability(date: string): Observable<Record<number, Record<string, number>>> {
    return this.http.get<Record<number, Record<string, number>>>(`${this.baseUrl}/shifts/availability`, {
      params: { date },
    });
  }
}
