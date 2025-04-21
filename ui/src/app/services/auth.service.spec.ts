import { TestBed } from '@angular/core/testing';
import { Router } from '@angular/router';
import { AuthService } from './auth.service';
import { sharedTestingProviders } from '../test-utils/shared-test-imports';

describe('AuthService', () => {
  let service: AuthService;
  let routerSpy: jasmine.SpyObj<Router>;

  const fakeToken = (expOffset: number) => {
    const payload = {
      exp: Math.floor(Date.now() / 1000) + expOffset,
    };
    const token = 'header.' + btoa(JSON.stringify(payload)) + '.signature';
    return token;
  };

  beforeEach(() => {
    routerSpy = jasmine.createSpyObj('Router', ['navigate']);

    TestBed.configureTestingModule({
      providers: [
        ...sharedTestingProviders,
        AuthService,
        { provide: Router, useValue: routerSpy }
      ]
    });

    service = TestBed.inject(AuthService);
    localStorage.clear();
  });

  afterEach(() => {
    service.stopPeriodicCheck();
    localStorage.clear();
  });

  it('should return false if no token is present', () => {
    expect(service.isAuthenticated()).toBeFalse();
  });

  it('should return false if token is expired', () => {
    localStorage.setItem('token', fakeToken(-60)); // expired 1 min ago
    expect(service.isAuthenticated()).toBeFalse();
  });

  it('should return true if token is valid', () => {
    localStorage.setItem('token', fakeToken(300)); // expires in 1 min
    expect(service.isAuthenticated()).toBeTrue();
  });

  it('should remove token and navigate to login on logout', () => {
    localStorage.setItem('token', fakeToken(60));
    service.logout();
    expect(localStorage.getItem('token')).toBeNull();
    expect(routerSpy.navigate).toHaveBeenCalledWith(['/login']);
  });
});
