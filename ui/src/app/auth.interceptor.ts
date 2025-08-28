import { inject } from '@angular/core';
import { HttpInterceptorFn, HttpRequest, HttpHandlerFn, HttpErrorResponse } from '@angular/common/http';
import { catchError, throwError } from 'rxjs';
import { AuthService } from './services/auth.service';

export const authInterceptor: HttpInterceptorFn = (req: HttpRequest<unknown>, next: HttpHandlerFn) => {
  const token = localStorage.getItem('token');
  const authService = inject(AuthService);

  if (token) {
    req = req.clone({
      setHeaders: { Authorization: `Bearer ${token}` },
    });
  }

  return next(req).pipe(
    catchError((error: HttpErrorResponse) => {
      if (error.status === 401) {
        const isExpiredOrMissing = !authService.isAuthenticated();
        const url = (req && (req as any).url) ? (req as any).url.toString() : '';
        const isAuthEndpoint = url.includes('/login') || url.includes('/oauth/token');
        // Only force logout if token is missing/expired or the 401 came from auth endpoints
        if (isExpiredOrMissing || isAuthEndpoint) {
          authService.logout();
        }
      }
      return throwError(() => error);
    })
  );
};
