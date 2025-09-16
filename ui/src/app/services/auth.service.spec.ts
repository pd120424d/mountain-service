import { TestBed } from '@angular/core/testing';
import { Router } from '@angular/router';
import { AuthService } from './auth.service';
import { sharedTestingProviders } from '../test-utils/shared-test-imports';
import { of } from 'rxjs';
import { MedicRole, TechnicalRole } from '../shared/models';
import { HttpTestingController } from '@angular/common/http/testing';
import { environment } from '../../environments/environment';

describe('AuthService', () => {
  let service: AuthService;
  let routerSpy: jasmine.SpyObj<Router>;
  let httpMock: HttpTestingController;

  const fakeToken = (expOffset: number) => {
    const payload = {
      exp: Math.floor(Date.now() / 1000) + expOffset,
    };
    const token = 'header.' + btoa(JSON.stringify(payload)) + '.signature';
    return token;
  };

  beforeEach(() => {
    routerSpy = jasmine.createSpyObj('Router', ['navigate'], { url: '/' });

    TestBed.configureTestingModule({
      providers: [
        ...sharedTestingProviders,
        AuthService,
        { provide: Router, useValue: routerSpy }
      ]
    });

    service = TestBed.inject(AuthService);
    httpMock = TestBed.inject(HttpTestingController);
    localStorage.clear();
  });

  afterEach(() => {
    service.stopPeriodicCheck();
    localStorage.clear();
  });

  it('should be created', () => {
    expect(service).toBeTruthy();
  });

  it('should decode the token correctly', () => {
    const token = fakeToken(300);
    const decoded = service['decodeToken'](token);
    expect(decoded.exp).toBeGreaterThan(Date.now() / 1000);
  });

  it('should set token on successful login', () => {
    const mockToken = 'mock-token';
    spyOn(service['http'], 'post').and.returnValue(of({ token: mockToken }));
    service.login({ username: 'test', password: 'test' }).subscribe();
    expect(localStorage.getItem('token')).toBe(mockToken);
  });

  it('it succeeds when logout is called without token', () => {
    // Ensure no token exists
    localStorage.removeItem('token');

    service.logout();

    // No HTTP call should be made
    httpMock.expectNone(`${environment.apiUrl}/logout`);

    // Navigates immediately
    expect(routerSpy.navigate).toHaveBeenCalledWith(['/login']);
    expect(localStorage.getItem('token')).toBeNull();
  });


  it('succeeds when logout posts to backend and clears token', () => {
    const tokenPayload = { id: 1, role: 'Medic', exp: Math.floor(Date.now() / 1000) + 300 };
    const token = 'h.' + btoa(JSON.stringify(tokenPayload)) + '.s';
    localStorage.setItem('token', token);

    service.logout();

    const req = httpMock.expectOne(`${environment.apiUrl}/logout`);
    expect(req.request.method).toBe('POST');
    expect(req.request.headers.get('Authorization')).toBe(`Bearer ${token}`);

    req.flush({ message: 'Successfully logged out' });

    expect(localStorage.getItem('token')).toBeNull();
    expect(routerSpy.navigate).toHaveBeenCalledWith(['/login']);
  });

  it('succeeds when logout returns error and still clears token', () => {
    const tokenPayload = { id: 1, role: 'Medic', exp: Math.floor(Date.now() / 1000) + 300 };
    const token = 'h.' + btoa(JSON.stringify(tokenPayload)) + '.s';
    localStorage.setItem('token', token);

    service.logout();

    const req = httpMock.expectOne(`${environment.apiUrl}/logout`);
    expect(req.request.method).toBe('POST');

    req.flush({ error: 'boom' }, { status: 500, statusText: 'Server Error' });

    expect(localStorage.getItem('token')).toBeNull();
    expect(routerSpy.navigate).toHaveBeenCalledWith(['/login']);
  });

  it('it succeeds when logout posts and then clears token and navigates', () => {
    const token = fakeToken(60);
    localStorage.setItem('token', token);

    service.logout();

    const req = httpMock.expectOne(`${environment.apiUrl}/logout`);
    expect(req.request.method).toBe('POST');
    req.flush({ message: 'ok' });

    expect(localStorage.getItem('token')).toBeNull();
    expect(routerSpy.navigate).toHaveBeenCalledWith(['/login']);
  });

  it('should return false if no token is present', () => {
    expect(service.isAuthenticated()).toBeFalse();
  });

  it('should return false if token is expired', () => {
    localStorage.setItem('token', fakeToken(-60)); // expired 1 min ago
    expect(service.isAuthenticated()).toBeFalse();
  });

  it('should return true if token is valid', () => {
    localStorage.setItem('token', fakeToken(300)); // expires in 5 minutes
    expect(service.isAuthenticated()).toBeTrue();
  });

  it('should return false if on unauthorized route', () => {
    Object.defineProperty(routerSpy, 'url', { value: '/login', writable: true });
    expect(service['isAauthorizedRoute']()).toBeFalse();
  });

  it('should return true if on authorized route', () => {
    Object.defineProperty(routerSpy, 'url', { value: '/shifts', writable: true });
    expect(service['isAauthorizedRoute']()).toBeTrue();
  });

  it('should return MedicRole if no token is present', () => {
    expect(service.getRole()).toBe(MedicRole);
  });

  it('should return the role from the token', () => {
    const token = 'header.eyJpZCI6MSwiZXhwIjoxNjg3MjIyOTI3LCJyb2xlIjoiVGVjaG5pY2FsIn0=.signature';
    localStorage.setItem('token', token);
    expect(service.getRole()).toBe(TechnicalRole);
  });

  it('should return the user id from the token', () => {
    const token = 'header.eyJpZCI6MSwiZXhwIjoxNjg3MjIyOTI3fQ==.signature';
    localStorage.setItem('token', token);
    expect(service.getUserId()).toBe('1');
  });

  it('should return false if user is not admin', () => {
    expect(service.isAdmin()).toBeFalse();
  });

  it('should return true if user is admin', () => {
    const token = 'header.eyJpZCI6MSwiZXhwIjoxNjg3MjIyOTI3LCJyb2xlIjoiQWRtaW5pc3RyYXRvciJ9.signature';
    localStorage.setItem('token', token);
    expect(service.isAdmin()).toBeTrue();
  });

  it('should call delete on resetAllData', () => {
    spyOn(service['http'], 'delete').and.returnValue(of({ message: 'success' }));
    service.resetAllData().subscribe();
    expect(service['http'].delete).toHaveBeenCalled();
  });

  it('should check if route is authorized', () => {
    // Test unauthorized routes
    Object.defineProperty(routerSpy, 'url', { value: '/home', writable: true });
    expect(service['isAauthorizedRoute']()).toBeFalse();

    Object.defineProperty(routerSpy, 'url', { value: '/login', writable: true });
    expect(service['isAauthorizedRoute']()).toBeFalse();

    Object.defineProperty(routerSpy, 'url', { value: '/employees/new', writable: true });
    expect(service['isAauthorizedRoute']()).toBeFalse();

    // Test authorized route
    Object.defineProperty(routerSpy, 'url', { value: '/employees', writable: true });
    expect(service['isAauthorizedRoute']()).toBeTrue();
  });

  it('should stop periodic check when interval exists', () => {
    // Create a mock subscription
    const mockSubscription = jasmine.createSpyObj('Subscription', ['unsubscribe']);
    service['intervalSub'] = mockSubscription;

    service.stopPeriodicCheck();

    expect(mockSubscription.unsubscribe).toHaveBeenCalled();
    // The method doesn't set intervalSub to undefined, it just unsubscribes
    expect(service['intervalSub']).toBe(mockSubscription);
  });

  it('should call post on restartService with deployment payload', () => {
    const dep = 'employee-service';
    service.restartService(dep).subscribe();

    const req = httpMock.expectOne(`${environment.apiUrl}/admin/k8s/restart`);
    expect(req.request.method).toBe('POST');
    expect(req.request.body).toEqual({ deployment: dep });
    req.flush({ message: 'ok' });
  });

});
