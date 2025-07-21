import { ComponentFixture, TestBed } from '@angular/core/testing';

import { LoginComponent } from './login.component';
import { sharedTestingProviders } from '../test-utils/shared-test-imports';
import { of, throwError } from 'rxjs';
import { ActivatedRoute } from '@angular/router';
import { AuthService } from '../services/auth.service';

describe('LoginComponent', () => {
  let component: LoginComponent;
  let fixture: ComponentFixture<LoginComponent>;
  let authService: jasmine.SpyObj<AuthService>;

  beforeEach(async () => {
    const authServiceSpy = jasmine.createSpyObj('AuthService', ['login']);

    TestBed.configureTestingModule({
      imports: [LoginComponent],
      providers: [
        ...sharedTestingProviders,
        { provide: ActivatedRoute, useValue: { queryParams: of({}) } },
        { provide: AuthService, useValue: authServiceSpy }
      ]
    }).compileComponents();

    fixture = TestBed.createComponent(LoginComponent);
    component = fixture.componentInstance;
    authService = TestBed.inject(AuthService) as jasmine.SpyObj<AuthService>;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });

  it('should show error message on 401 login failure', () => {
    const error = { status: 401 };
    authService.login.and.returnValue(throwError(() => error));
    component.credentials = { username: 'test', password: 'wrong' };

    component.onLogin();

    expect(component.loginError).toBe(true);
  });

  it('should reset error message on new login attempt', () => {
    component.loginError = true;
    authService.login.and.returnValue(of({ token: 'fake-token' }));
    component.credentials = { username: 'test', password: 'correct' };

    component.onLogin();

    expect(component.loginError).toBe(false);
  });
});
