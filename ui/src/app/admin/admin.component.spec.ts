import { ComponentFixture, TestBed } from '@angular/core/testing';
import { Router } from '@angular/router';
import { AuthService } from '../services/auth.service';
import { sharedTestingProviders } from '../test-utils/shared-test-imports';
import { AdminComponent } from './admin.component';
import { of } from 'rxjs';

describe('AdminComponent', () => {
  let component: AdminComponent;
  let fixture: ComponentFixture<AdminComponent>;

  beforeEach(async () => {
    const authServiceSpy = jasmine.createSpyObj('AuthService', ['isAdmin', 'resetAllData']);
    const routerSpy = jasmine.createSpyObj('Router', ['navigate']);

    await TestBed.configureTestingModule({
      imports: [AdminComponent],
      providers: [
        ...sharedTestingProviders,
        { provide: AuthService, useValue: authServiceSpy },
        { provide: Router, useValue: routerSpy }
      ]
    }).compileComponents();

    fixture = TestBed.createComponent(AdminComponent);
    component = fixture.componentInstance;

    const authService = TestBed.inject(AuthService) as jasmine.SpyObj<AuthService>;
    authService.isAdmin.and.returnValue(true);
    authService.resetAllData.and.returnValue(of({ message: 'success' }));

    fixture.detectChanges();
  });

  it('should create and not fail with subscribe error', () => {
    expect(component).toBeTruthy();
    expect(component.isResetting).toBeFalse();
    expect(component.resetSuccess).toBeFalse();
    expect(component.resetError).toBeFalse();
    expect(component.errorMessage).toBe('');
  });

  it('should call resetAllData without subscribe error when user confirms', () => {
    spyOn(window, 'confirm').and.returnValue(true);

    const authService = TestBed.inject(AuthService) as jasmine.SpyObj<AuthService>;

    authService.resetAllData.and.returnValue(of({ message: 'success' }));

    expect(() => {
      component.onResetAllData();
    }).not.toThrow();

    expect(authService.resetAllData).toHaveBeenCalled();

    expect(component.isResetting).toBeFalse();
    expect(component.resetSuccess).toBeTrue();
    expect(component.resetError).toBeFalse();
  });

  it('should not call resetAllData when user cancels confirmation', () => {
    spyOn(window, 'confirm').and.returnValue(false);

    const authService = TestBed.inject(AuthService) as jasmine.SpyObj<AuthService>;

    expect(() => {
      component.onResetAllData();
    }).not.toThrow();

    expect(authService.resetAllData).not.toHaveBeenCalled();

    expect(component.isResetting).toBeFalse();
    expect(component.resetSuccess).toBeFalse();
    expect(component.resetError).toBeFalse();
  });

  it('should navigate to home when goBack is called', () => {
    const router = TestBed.inject(Router) as jasmine.SpyObj<Router>;

    component.goBack();

    expect(router.navigate).toHaveBeenCalledWith(['/']);
  });

});
