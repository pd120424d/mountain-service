
import { Injectable } from '@angular/core';
import { HttpClient } from '@angular/common/http';
import { interval, Observable, Subscription, tap } from 'rxjs';
import { Router } from '@angular/router';
import { environment } from '../../environments/environment'; // Import environment variables
import { EmployeeRole, MedicRole } from '../employee/employee.model';

@Injectable({
  providedIn: 'root',
})
export class AuthService {
  private apiUrl = `${environment.apiUrl}`;
  private checkInterval = 300 * 1000; // Check every 5 minutes
  private intervalSub: Subscription | null = null;

  constructor(private http: HttpClient, private router: Router) {
    this.startPeriodicTokenCheck();
  }

  login(credentials: { username: string; password: string }): Observable<{ token: string }> {
    return this.http.post<{ token: string }>(`${this.apiUrl}/login`, credentials).pipe(
      tap(response => {
        localStorage.setItem('token', response.token);
      })
    );
  }

  logout(): void {
    localStorage.removeItem('token');
    this.router.navigate(['/login']);
  }

  isAuthenticated(): boolean {
    const token = localStorage.getItem('token');
    if (!token) return false;

    const payload = this.decodeToken(token);
    return payload && payload.exp > Date.now() / 1000;
  }

  private decodeToken(token: string): any {
    try {
      return JSON.parse(atob(token.split('.')[1]));
    } catch (e) {
      return null;
    }
  }

  private startPeriodicTokenCheck(): void {
    this.intervalSub = interval(this.checkInterval).subscribe(() => {
      if (!this.isAuthenticated() && !this.isUnauthorizedRoute()) {
        this.logout();
      }
    });
  }

  private isUnauthorizedRoute(): boolean {
    return this.router.url !== '/login' && this.router.url !== '/employees/new'
  }

  stopPeriodicCheck(): void {
    this.intervalSub?.unsubscribe();
  }

  getRole(): EmployeeRole {
    const token = localStorage.getItem('token');
    if (!token) return MedicRole;

    const payload = JSON.parse(atob(token.split('.')[1])); // Decode JWT payload
    return payload.role || MedicRole;
  }

  getUserId(): string {
    const token = localStorage.getItem('token');
    if (!token) return '';

    const payload = JSON.parse(atob(token.split('.')[1])); // Decode JWT payload
    return payload.id || '';
  }
}
