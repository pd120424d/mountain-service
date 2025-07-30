import { TestBed } from '@angular/core/testing';
import { HttpInterceptorFn, HttpRequest, HttpErrorResponse } from '@angular/common/http';
import { HttpTestingController, provideHttpClientTesting } from '@angular/common/http/testing';
import { provideHttpClient, withInterceptors } from '@angular/common/http';
import { authInterceptor } from './auth.interceptor';
import { AuthService } from './services/auth.service';
import { sharedTestingProviders } from './test-utils/shared-test-imports';
import { throwError } from 'rxjs';

describe('authInterceptor', () => {
  let httpMock: HttpTestingController;
  let authService: jasmine.SpyObj<AuthService>;

  const interceptor: HttpInterceptorFn = (req, next) => 
    TestBed.runInInjectionContext(() => authInterceptor(req, next));

  beforeEach(() => {
    const authSpy = jasmine.createSpyObj('AuthService', ['logout']);

    TestBed.configureTestingModule({
      providers: [
        ...sharedTestingProviders,
        provideHttpClient(withInterceptors([authInterceptor])),
        provideHttpClientTesting(),
        { provide: AuthService, useValue: authSpy }
      ]
    });

    httpMock = TestBed.inject(HttpTestingController);
    authService = TestBed.inject(AuthService) as jasmine.SpyObj<AuthService>;
  });

  afterEach(() => {
    httpMock.verify();
    localStorage.clear();
  });

  it('should add Authorization header when token exists', () => {
    const token = 'test-token';
    localStorage.setItem('token', token);

    const mockRequest = new HttpRequest('GET', '/test');
    const mockNext = jasmine.createSpy('next').and.returnValue(throwError(() => new Error('test')));

    TestBed.runInInjectionContext(() => {
      authInterceptor(mockRequest, mockNext).subscribe({
        error: () => {} // Handle the error to complete the test
      });
    });

    expect(mockNext).toHaveBeenCalledWith(
      jasmine.objectContaining({
        headers: jasmine.objectContaining({
          lazyUpdate: jasmine.arrayContaining([
            jasmine.objectContaining({
              name: 'Authorization',
              value: `Bearer ${token}`
            })
          ])
        })
      })
    );
  });

  it('should not add Authorization header when no token exists', () => {
    localStorage.removeItem('token');

    const mockRequest = new HttpRequest('GET', '/test');
    const mockNext = jasmine.createSpy('next').and.returnValue(throwError(() => new Error('test')));

    TestBed.runInInjectionContext(() => {
      authInterceptor(mockRequest, mockNext).subscribe({
        error: () => {} // Handle the error to complete the test
      });
    });

    expect(mockNext).toHaveBeenCalledWith(mockRequest);
  });

  it('should call authService.logout() on 401 error', () => {
    const mockRequest = new HttpRequest('GET', '/test');
    const error401 = new HttpErrorResponse({ status: 401 });
    const mockNext = jasmine.createSpy('next').and.returnValue(throwError(() => error401));

    TestBed.runInInjectionContext(() => {
      authInterceptor(mockRequest, mockNext).subscribe({
        error: (error) => {
          expect(error).toBe(error401);
          expect(authService.logout).toHaveBeenCalled();
        }
      });
    });
  });

  it('should not call authService.logout() on non-401 errors', () => {
    const mockRequest = new HttpRequest('GET', '/test');
    const error500 = new HttpErrorResponse({ status: 500 });
    const mockNext = jasmine.createSpy('next').and.returnValue(throwError(() => error500));

    TestBed.runInInjectionContext(() => {
      authInterceptor(mockRequest, mockNext).subscribe({
        error: (error) => {
          expect(error).toBe(error500);
          expect(authService.logout).not.toHaveBeenCalled();
        }
      });
    });
  });
});
