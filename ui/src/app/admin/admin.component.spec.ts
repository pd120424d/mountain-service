import { ComponentFixture, TestBed } from '@angular/core/testing';
import { Router } from '@angular/router';
import { AuthService } from '../services/auth.service';
import { sharedTestingProviders } from '../test-utils/shared-test-imports';
import { AdminComponent } from './admin.component';
import { of, throwError } from 'rxjs';
import { TranslateService } from '@ngx-translate/core';

describe('AdminComponent', () => {
  let component: AdminComponent;
  let fixture: ComponentFixture<AdminComponent>;

  beforeEach(async () => {
    const authServiceSpy = jasmine.createSpyObj('AuthService', ['isAdmin', 'resetAllData', 'restartService']);
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

  it('should not open restart modal when no deployment is selected', () => {
    component.selectedDeployment = '';
    component.openRestartModal();
    expect(component.showRestartModal).toBeFalse();
  });

  it('should open restart modal when deployment is selected and reset flags', () => {
    const authService = TestBed.inject(AuthService) as jasmine.SpyObj<AuthService>;
    authService.restartService.and.returnValue(of({ message: 'ok' }));

    component.restartError = true;
    component.restartSuccess = true;
    component.restartMessage = 'old';

    component.selectedDeployment = 'employee-service';
    component.openRestartModal();

    expect(component.showRestartModal).toBeTrue();
    expect(component.restartError).toBeFalse();
    expect(component.restartSuccess).toBeFalse();
    expect(component.restartMessage).toBe('');
  });

  it('should close restart modal on cancel', () => {
    component.showRestartModal = true;
    component.cancelRestart();
    expect(component.showRestartModal).toBeFalse();
  });

  it('should handle successful restart', () => {
    const authService = TestBed.inject(AuthService) as jasmine.SpyObj<AuthService>;
    authService.restartService.and.returnValue(of({ message: 'restart triggered' }));

    component.selectedDeployment = 'employee-service';
    component.openRestartModal();
    component.confirmRestart();

    expect(component.isRestarting).toBeFalse();
    expect(component.showRestartModal).toBeFalse();
    expect(component.restartSuccess).toBeTrue();
    expect(component.restartError).toBeFalse();
    expect(component.restartMessage).toContain('restart');
  });

  it('should handle restart error and show translated message', () => {
    const authService = TestBed.inject(AuthService) as jasmine.SpyObj<AuthService>;
    const translate = TestBed.inject(TranslateService);
    spyOn(translate, 'instant').and.returnValue('TRANSLATED_ERROR');

    authService.restartService.and.returnValue(throwError(() => ({ error: { error: 'boom' } })));

    component.selectedDeployment = 'employee-service';
    component.openRestartModal();
    component.confirmRestart();

    expect(component.isRestarting).toBeFalse();
    expect(component.restartSuccess).toBeFalse();
    expect(component.restartError).toBeTrue();
    expect(component.restartMessage === 'boom' || component.restartMessage === 'TRANSLATED_ERROR').toBeTrue();
  });


});
