
import { Injectable } from '@angular/core';
import { HttpClient, HttpBackend, HttpHeaders } from '@angular/common/http';
import { interval, Observable, Subscription, tap, finalize, Subject } from 'rxjs';
import { Router } from '@angular/router';
import { environment } from '../../environments/environment';
import { EmployeeRole, MedicRole, AdministratorRole } from '../shared/models';

@Injectable({
  providedIn: 'root',
})
export class AuthService {
  private apiUrl = `${environment.apiUrl}`;
  private checkInterval = 300 * 1000; // Check every 5 minutes
  private intervalSub: Subscription | null = null;
  private bareHttp!: HttpClient;
  private authChangedSubject = new Subject<'login' | 'logout'>();
  authChanged$ = this.authChangedSubject.asObservable();


  constructor(private http: HttpClient, private router: Router, private httpBackend: HttpBackend) {
    this.bareHttp = new HttpClient(this.httpBackend);
    this.startPeriodicTokenCheck();
  }

  private startPeriodicTokenCheck(): void {
    this.intervalSub = interval(this.checkInterval).subscribe(() => {
      if (!this.isAuthenticated() && this.isAauthorizedRoute()) {
        this.logout();
      }
    });
  }

  decodeToken(token: string): any {
    try {
      return JSON.parse(atob(token.split('.')[1]));
    } catch (e) {
      return null;
    }
  }

  login(credentials: { username: string; password: string }): Observable<{ token: string }> {
    return this.http.post<{ token: string }>(`${this.apiUrl}/login`, credentials).pipe(
      tap(response => {
        localStorage.setItem('token', response.token);
        // notify listeners (e.g., header) to refresh current user and UI
        this.authChangedSubject.next('login');
      })
    );
  }

  logout(): void {
    const token = localStorage.getItem('token');
    const finish = () => {
      localStorage.removeItem('token');
      this.stopPeriodicCheck();
      // notify listeners
      this.authChangedSubject.next('logout');
      this.router.navigate(['/login']);
    };
    if (token) {
      const headers = new HttpHeaders({ Authorization: `Bearer ${token}` });
      // Use bare HttpClient to bypass interceptors and avoid 401 -> logout recursion
      this.bareHttp.post<{ message: string }>(`${this.apiUrl}/logout`, {}, { headers }).pipe(
        finalize(finish)
      ).subscribe({ next: () => {}, error: () => {} });
    } else {
      finish();
    }
  }

  isAuthenticated(): boolean {
    const token = localStorage.getItem('token');
    if (!token) return false;

    const payload = this.decodeToken(token);
    return payload && payload.exp > Date.now() / 1000;
  }

  private isAauthorizedRoute(): boolean {
    return this.router.url !== '/home' &&
     this.router.url !== '/login' &&
     this.router.url !== '/employees/new';
  }

  stopPeriodicCheck(): void {
    this.intervalSub?.unsubscribe();
  }

  getRole(): EmployeeRole {
    const token = localStorage.getItem('token');
    if (!token) return MedicRole;

    const payload = JSON.parse(atob(token.split('.')[1]));
    return payload.role || MedicRole;
  }

  getUserId(): string {
    const token = localStorage.getItem('token');
    if (!token) return '';

    const payload = JSON.parse(atob(token.split('.')[1]));
    return (payload.id).toString() || '';
  }

  isAdmin(): boolean {
    return this.getRole() === AdministratorRole;
  }

  resetAllData(): Observable<{ message: string }> {
    return this.http.delete<{ message: string }>(`${this.apiUrl}/admin/reset`);
  }
}
